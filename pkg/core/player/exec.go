package player

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

// ErrActionNotReady is returned if the requested action is not ready; this could be
// due to any of the following:
//   - Insufficient energy (burst only)
//   - Ability on cooldown
//   - Player currently in animation
var (
	// exec-specfic errors
	ErrActionNotReady        = errors.New("action is not ready yet; cannot be executed")
	ErrPlayerNotReady        = errors.New("player still in animation; cannot execute action")
	ErrInvalidAirborneAction = errors.New("player must use low_plunge or high_plunge while airborne")
	ErrActionNoOp            = errors.New("action is a noop")
	// shared character-specific errors
	ErrInvalidChargeAction = errors.New("need to use attack right before charge")
)

// ReadyCheck returns nil action is ready, else returns error representing why action is not ready
func (h *Handler) ReadyCheck(t info.Action, k keys.Char, param map[string]int) error {
	// check animation state
	if h.Animation.IsAnimationLocked(t) {
		return ErrPlayerNotReady
	}
	char := h.chars[h.active]
	// check for energy, cd, etc..
	// TODO: make sure there is a default check for charge attack/dash stams in char implementation
	// this should deal with Ayaka/Mona's drain vs straight up consumption
	if ok, reason := char.ActionReady(t, param); !ok {
		h.Events.ActionFailed.Emit(event.ActionFailedEvent{
			CharIndex: h.active,
			Action:    t,
			Params:    param,
			Reason:    reason,
		})
		return ErrActionNotReady
	}

	stamCheck := func(t info.Action, param map[string]int) (float64, bool) {
		req := char.AbilStamCost(t, param)
		return req, h.stam >= req
	}

	switch t {
	case info.ActionCharge: // require special calc for stam
		amt, ok := stamCheck(t, param)
		if !ok {
			h.Log.NewEvent("insufficient stam: charge attack", glog.LogWarnings, -1).
				Write("have", h.stam).
				Write("cost", amt)
			h.Events.ActionFailed.Emit(event.ActionFailedEvent{
				CharIndex: h.active,
				Action:    t,
				Params:    param,
				Reason:    info.InsufficientStamina,
			})
			return ErrActionNotReady
		}
	case info.ActionDash: // require special calc for stam
		// dash handles it in the action itself
		amt, ok := stamCheck(t, param)
		if !ok {
			h.Log.NewEvent("insufficient stam: dash", glog.LogWarnings, -1).
				Write("have", h.stam).
				Write("cost", amt)
			h.Events.ActionFailed.Emit(event.ActionFailedEvent{
				CharIndex: h.active,
				Action:    t,
				Params:    param,
				Reason:    info.InsufficientStamina,
			})
			return ErrActionNotReady
		}

		// dash is still on cooldown and is locked out, cannot dash again until CD expires
		if h.dashLockout && h.dashCDExpirationFrame > *h.F {
			h.Log.NewEvent("dash on cooldown", glog.LogWarnings, -1).
				Write("dash_cd_expiration", h.dashCDExpirationFrame-*h.F)
			h.Events.ActionFailed.Emit(event.ActionFailedEvent{
				CharIndex: h.active,
				Action:    t,
				Params:    param,
				Reason:    info.DashCD,
			})
			return ErrActionNotReady
		}
	case info.ActionSwap:
		if h.active == h.charPos[k] {
			// even though noop this action is still ready
			return nil
		}
		if h.swapCD > 0 {
			h.Events.ActionFailed.Emit(event.ActionFailedEvent{
				CharIndex: h.active,
				Action:    t,
				Params:    param,
				Reason:    info.SwapCD,
			})
			return ErrActionNotReady
		}
	}

	return nil
}

// Exec will forcefully execute an action t regardless if t is ready or not. The assumption is
// that whatever caller of Exec would have first checked ReadyCheck where ever relevant
// before calling Exec.
//
// The separation allows for forcefully execution of certain actions such as swap bypassing
// swapCD if any
func (h *Handler) Exec(t info.Action, k keys.Char, param map[string]int) error {
	char := h.chars[h.active]

	// special airborne handler; if airborne the next action MUST be attack otherwise error
	if h.airborne != info.AirborneGrounded && t != info.ActionLowPlunge && t != info.ActionHighPlunge {
		return ErrInvalidAirborneAction
	}

	var err error
	switch t {
	case info.ActionCharge: // require special calc for stam
		req := char.AbilStamCost(t, param)
		h.UseStam(req, t)
		err = h.useAbility(t, param, char.ChargeAttack) // TODO: make sure characters are consuming stam in charge attack function
	case info.ActionDash:
		err = h.useAbility(t, param, char.Dash) // TODO: make sure characters are consuming stam in dashes
	case info.ActionJump:
		err = h.useAbility(t, param, char.Jump)
	case info.ActionWalk:
		err = h.useAbility(t, param, char.Walk)
	case info.ActionAim:
		err = h.useAbility(t, param, char.Aimed)
	case info.ActionSkill:
		err = h.useAbility(t, param, char.Skill)
	case info.ActionBurst:
		err = h.useAbility(t, param, char.Burst)
	case info.ActionAttack:
		err = h.useAbility(t, param, char.Attack)
	case info.ActionHighPlunge:
		err = h.useAbility(t, param, char.HighPlungeAttack)
		h.airborne = info.AirborneGrounded
	case info.ActionLowPlunge:
		err = h.useAbility(t, param, char.LowPlungeAttack)
		h.airborne = info.AirborneGrounded
	case info.ActionSwap:
		if h.active == h.charPos[k] {
			return ErrActionNoOp
		}
		if h.swapCD > 0 {
			// since we allow force swap, this is ok but will emit an extra log anyways just in case
			h.Log.NewEventBuildMsg(glog.LogActionEvent, h.active, "swapping ", h.chars[h.active].GetBase().Key.String(), " to ", h.chars[h.charPos[k]].GetBase().Key.String(), " (bypassed cd)").
				Write("swap_cd", h.swapCD)
			h.swapCD = 0
		} else {
			h.Log.NewEventBuildMsg(glog.LogActionEvent, h.active, "swapping ", h.chars[h.active].GetBase().Key.String(), " to ", h.chars[h.charPos[k]].GetBase().Key.String())
		}

		x := action.Info{
			Frames: func(info.Action) int {
				return h.Delays.Swap
			},
			AnimationLength: h.Delays.Swap,
			CanQueueAfter:   h.Delays.Swap,
			State:           info.AnimationStateSwap,
		}
		x.QueueAction(h.swap(k), h.Delays.Swap)
		h.Animation.SetActionUsed(h.active, t, &x)
		h.lastAction.Type = t
		h.lastAction.Param = param
		h.lastAction.Char = h.active
	default:
		return fmt.Errorf("invalid action: %v", t)
	}
	if err != nil {
		return err
	}

	if t != info.ActionAttack {
		h.ResetAllNormalCounter()
	}

	h.Events.ActionExec.Emit(event.ActionExecEvent{
		CharIndex: h.active,
		Action:    t,
		Params:    param,
	})

	return nil
}

func (h *Handler) useAbility(
	t info.Action,
	param map[string]int,
	f func(p map[string]int) (action.Info, error),
) error {
	actionToEvent := map[info.Action]event.ActionEventHandler{
		info.ActionDash:       h.Events.Dash,
		info.ActionSkill:      h.Events.Skill,
		info.ActionBurst:      h.Events.Burst,
		info.ActionAttack:     h.Events.Attack,
		info.ActionCharge:     h.Events.ChargeAttack,
		info.ActionLowPlunge:  h.Events.Plunge,
		info.ActionHighPlunge: h.Events.Plunge,
		info.ActionAim:        h.Events.AimShoot,
	}

	state, ok := actionToEvent[t]
	if ok {
		state.Emit(event.ActionEvent{
			CharIndex: h.active,
			Params:    param,
		})
	}
	info, err := f(param)
	if err != nil {
		return err
	}
	h.Animation.SetActionUsed(h.active, t, &info)
	if info.FramePausedOnHitlag == nil {
		info.FramePausedOnHitlag = h.chars[h.active].FramePausedOnHitlag
	}

	h.lastAction.Type = t
	h.lastAction.Param = param
	h.lastAction.Char = h.active

	h.Log.NewEventBuildMsg(
		glog.LogActionEvent,
		h.active,
		"executed ", t.String(),
	).
		Write("action", t.String()).
		Write("stam_post", h.stam).
		Write("swap_cd_post", h.swapCD)
	return nil
}
