package modifier

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Handler struct {
	core      core.Core
	owner     keys.Target
	modifiers []*Instance
}

func New(core core.Core, owner keys.Target) *Handler {
	h := &Handler{
		core:  core,
		owner: owner,
	}
	h.subscribe()
	return h
}

func (h *Handler) Core() core.Core {
	return h.core
}

func (h *Handler) Owner() keys.Target {
	return h.owner
}

func (h *Handler) SetOwner(owner keys.Target) {
	h.owner = owner
}

func (h *Handler) Count() int {
	return len(h.modifiers)
}

func (h *Handler) CountByKey(key string) int {
	count := 0
	for _, mod := range h.modifiers {
		if mod.key == key {
			count += mod.count
		}
	}
	return count
}

func (h *Handler) GetDuration(key string) int {
	// find the longest duration
	dur := 0
	for _, mod := range h.modifiers {
		if mod.key == key {
			dur = max(dur, mod.Duration())
		}
	}
	return dur
}

func (h *Handler) HasModifier(key string) bool {
	for _, mod := range h.modifiers {
		if mod.key == key {
			return true
		}
	}
	return false
}

func (h *Handler) Tick(hitlag bool) {
	n := 0
	var removedMods []*Instance
	for _, mod := range h.modifiers {
		remove := false

		if !hitlag || !mod.hitlagAffected {
			mod.elementDurability -= mod.decayRate
			mod.currentThinkInterval-- // TODO: onThinkIntervalIsFixedUpdate
		}
		if mod.thinkInterval > 0 && mod.currentThinkInterval <= 0 {
			h.emitThinkInterval(mod)
			mod.currentThinkInterval = mod.thinkInterval
		}
		if mod.elementDurability <= 0 {
			remove = true
			removedMods = append(removedMods, mod)
		}

		if !remove {
			h.modifiers[n] = mod
			n++
		}
	}
	h.modifiers = h.modifiers[:n]
	h.emitRemove(removedMods)
}

func (h *Handler) Stats() []attributes.ModifierChange {
	stats := make([]attributes.ModifierChange, 0, len(h.modifiers))
	for _, mod := range h.modifiers {
		h.emitStat(mod)
		props := mod.Props()
		stats = append(stats, attributes.ModifierChange{
			Props:  props,
			Reason: mod.key,
			Expiry: h.core.F() + mod.Duration(),
		})
	}
	return stats
}

func (h *Handler) ExtendByHitlag(dur int) {
	for _, mod := range h.modifiers {
		if mod.hitlagAffected && mod.elementDurability > 0 {
			mod.elementDurability += info.Durability(dur) * mod.decayRate
			h.emitExtendDuration(mod)
		}
	}
}
