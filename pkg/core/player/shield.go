package player

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

// Handler keeps track of all the active shields
// However we do care which is active since:
// - global shields: only active is shielded
// - 1 character shields: active only shieled if is target of shield
type ShieldHandler struct { //nolint:revive // cannot just name this Handler because then there is a conflict with Handler in player package
	log    glog.Logger
	events *event.System
	f      *int
	player *Handler

	shields []core.Shield
}

func NewShieldHandler(f *int, log glog.Logger, events *event.System, player *Handler) *ShieldHandler {
	h := &ShieldHandler{
		log:     log,
		events:  events,
		f:       f,
		player:  player,
		shields: make([]core.Shield, 0, info.ShieldEndType),
	}
	return h
}

func (h *ShieldHandler) Count() int { return len(h.shields) }

func (h *ShieldHandler) CharacterIsShielded(char, active int) bool {
	for _, v := range h.shields {
		target := v.ShieldTarget()
		if (target == -1 && char == active) || target == char {
			return true
		}
	}
	return false
}

func (h *ShieldHandler) Get(t info.ShieldType) core.Shield {
	for _, v := range h.shields {
		if v.Type() == t {
			return v
		}
	}
	return nil
}

// TODO: do shields get affected by hitlag? if so.. which timer? active char?
func (h *ShieldHandler) Add(shd core.Shield) {
	// we always assume over write of the same type and target
	ind := -1
	for i, v := range h.shields {
		if v.Type() == shd.Type() && v.ShieldTarget() == shd.ShieldTarget() {
			ind = i
		}
	}
	if ind > -1 {
		h.log.NewEvent("shield overridden", glog.LogShieldEvent, -1).
			Write("overwrite", true).
			Write("name", shd.Desc()).
			Write("hp", shd.CurrentHP()).
			Write("ele", shd.Element()).
			Write("expiry", shd.Expiry())
		h.shields[ind].OnOverwrite()
		h.shields[ind] = shd
	} else {
		h.shields = append(h.shields, shd)
		h.log.NewEvent("shield added", glog.LogShieldEvent, -1).
			Write("overwrite", false).
			Write("name", shd.Desc()).
			Write("hp", shd.CurrentHP()).
			Write("ele", shd.Element()).
			Write("expiry", shd.Expiry())
	}
	h.events.Shielded.Emit(event.ShieldedEvent{
		Type: shd.Type(),
	})
}

func (h *ShieldHandler) List() []core.Shield {
	return h.shields
}

func (h *ShieldHandler) OnDamage(char, active int, dmg float64, ele attributes.ElementType) float64 {
	// find shield bonuses
	bonus := h.ShieldBonus(char)
	mintaken := dmg // min of damage taken
	n := 0
	for _, v := range h.shields {
		target := v.ShieldTarget()
		if target == -1 && char != active {
			continue
		}
		if target != -1 && char != target {
			continue
		}
		preHp := v.CurrentHP()
		taken, ok := v.OnDamage(dmg, ele, bonus)
		h.log.NewEvent(
			"shield taking damage",
			glog.LogShieldEvent,
			-1,
		).Write("name", v.Desc()).
			Write("ele", v.Element()).
			Write("dmg", dmg).
			Write("previous_hp", preHp).
			Write("dmg_after_shield", taken).
			Write("current_hp", v.CurrentHP()).
			Write("expiry", v.Expiry())
		if taken < mintaken {
			mintaken = taken
		}
		if ok {
			h.shields[n] = v
			n++
		} else {
			// shield broken
			h.log.NewEvent(
				"shield broken",
				glog.LogShieldEvent,
				-1,
			).Write("name", v.Desc()).
				Write("ele", v.Element()).
				Write("expiry", v.Expiry())
			h.events.ShieldBreak.Emit(event.ShieldBreakEvent{
				Type: v.Type(),
			})
		}
	}
	h.shields = h.shields[:n]
	return mintaken
}

func (h *ShieldHandler) Tick() {
	n := 0
	for _, v := range h.shields {
		if v.Expiry() == *h.f {
			v.OnExpire()
			h.log.NewEvent("shield expired", glog.LogShieldEvent, -1).
				Write("name", v.Desc()).
				Write("hp", v.CurrentHP())
			h.events.ShieldBreak.Emit(event.ShieldBreakEvent{
				Type: v.Type(),
			})
		} else {
			h.shields[n] = v
			n++
		}
	}
	h.shields = h.shields[:n]
}

func (h *ShieldHandler) ShieldBonus(charIndex int) float64 {
	char := h.player.ByIndex(charIndex)
	return char.Stat(attributes.ShieldStrength)
}
