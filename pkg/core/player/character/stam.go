package character

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

// ActionStam provides default implementation for stam cost for charge and dash
// character should override this though
func (c *Character) ActionStam(a info.Action, p map[string]int) float64 {
	switch a {
	case info.ActionCharge:
		// 20 sword (most)
		// 25 polearm
		// 40 per second claymore
		// 50 catalyst
		switch c.Weapon.Class {
		case info.WeaponClassSword:
			return 20
		case info.WeaponClassSpear:
			return 25
		case info.WeaponClassCatalyst:
			return 50
		case info.WeaponClassClaymore:
			return 0
		case info.WeaponClassBow:
			return 0
		default:
			return 0
		}
	case info.ActionDash:
		// 18 per
		return 18
	default:
		return 0
	}
}

func (c *Character) AbilStamCost(a info.Action, p map[string]int) float64 {
	r := 1 + c.Stat(attributes.CostStamP)
	if r < 0 {
		r = 0
	}
	return r * c.ActionStam(a, p)
}

func (c *Character) Dash(p map[string]int) (action.Info, error) {
	// Execute dash CD logic
	c.ApplyDashCD()

	// consume stamina at end of the dash
	c.QueueDashStaminaConsumption(p)

	length := c.DashLength()
	dashJumpLength := c.DashToJumpLength()
	return action.Info{
		Frames: func(a info.Action) int {
			switch a {
			case info.ActionJump:
				return dashJumpLength
			default:
				return length
			}
		},
		AnimationLength: length,
		CanQueueAfter:   dashJumpLength,
		State:           info.AnimationStateDash,
	}, nil
}

// set the dash CD. If the dash was on CD when this dash executes, lockout dash
func (c *Character) ApplyDashCD() {
	var evt glog.Event

	if c.core.PlayerRemainingDashCD() > 0 {
		c.core.SetPlayerDashCD(true, 1.5*60)
		evt = c.core.Log().NewEvent("dash cooldown triggered", glog.LogCooldownEvent, c.Index)
	} else {
		c.core.SetPlayerDashCD(false, 0.8*60)
		evt = c.core.Log().NewEvent("dash lockout evaluation started", glog.LogCooldownEvent, c.Index)
	}

	evt.Write("lockout", c.core.PlayerDashLockout()).
		Write("expiry", c.core.PlayerRemainingDashCD()).
		Write("expiry_frame", c.core.PlayerRemainingDashCD()+c.core.F())
}

func (c *Character) QueueDashStaminaConsumption(p map[string]int) {
	// consume stam at the end
	c.core.Tasks().Add(func() {
		stam := c.AbilStamCost(info.ActionDash, p)
		c.core.PlayerUseStam(stam, info.ActionDash)
	}, c.DashLength()-1)
}

func (c *Character) DashLength() int {
	switch c.CharBody {
	case info.BodyBoy, info.BodyLoli:
		return 21
	case info.BodyMale:
		return 19
	case info.BodyLady:
		return 22
	default:
		return 20
	}
}

func (c *Character) DashToJumpLength() int {
	switch c.CharBody {
	case info.BodyGirl, info.BodyLoli:
		return 4
	case info.BodyBoy:
		return 2
	default:
		return 3
	}
}

func (c *Character) Jump(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(info.XianyunAirborneBuff) {
		c.core.SetPlayerAirborne(info.AirborneXianyun)
		// 4/8 for claymore/bow/catalyst and 5/9 for sword/polearm
		lowPlunge := 4
		highPlunge := 8
		switch c.Weapon.Class {
		case info.WeaponClassSword, info.WeaponClassSpear:
			lowPlunge = 5
			highPlunge = 9
		}

		animLength := 60 // Upperbound for jump for high/low plunge
		return action.Info{
			Frames: func(a info.Action) int {
				switch a {
				case info.ActionLowPlunge:
					return lowPlunge
				case info.ActionHighPlunge:
					return highPlunge
				default:
					return animLength // This is expected to later lead to action error because no other action besides plunges can be done while AirborneXianyun
				}
			},
			AnimationLength: animLength,
			CanQueueAfter:   lowPlunge, // earliest cancel
			State:           info.AnimationStateJump,
		}, nil
	}
	f := c.JumpLength()
	return action.Info{
		Frames:          func(info.Action) int { return f },
		AnimationLength: f,
		CanQueueAfter:   f,
		State:           info.AnimationStateJump,
	}, nil
}

func (c *Character) JumpLength() int {
	if c.core.GetPlayerLastAction().Type == info.ActionDash {
		switch c.CharBody {
		case info.BodyGirl, info.BodyBoy:
			return 34
		default:
			return 37
		}
	}
	switch c.CharBody {
	case info.BodyBoy, info.BodyGirl:
		return 31
	case info.BodyMale:
		return 28
	case info.BodyLady:
		return 32
	case info.BodyLoli:
		return 29
	default:
		return 30
	}
}

func (c *Character) Walk(p map[string]int) (action.Info, error) {
	f, ok := p["f"]
	if !ok {
		f = 1
	}
	return action.Info{
		Frames:          func(info.Action) int { return f },
		AnimationLength: f,
		CanQueueAfter:   f,
		State:           info.AnimationStateWalk,
	}, nil
}
