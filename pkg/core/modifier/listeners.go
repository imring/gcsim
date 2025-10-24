package modifier

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Listeners struct {
	// mod
	OnAdd           func(mod *Instance)
	OnRemove        func(mod *Instance)
	OnThinkInterval func(mod *Instance)

	// stat
	OnStat func(mod *Instance)

	// combat
	OnEnemyHit func(mod *Instance, target keys.Target, attack *info.AttackEvent)

	// action
	OnCharacterAction func(mod *Instance, char int, action info.Action, p map[string]int)
}

func (h *Handler) subscribe() {
	events := h.core.Events()

	events.EnemyHit.Subscribe(h.onEnemyHit)
	events.CharacterAction.Subscribe(h.onCharacterAction)
}

func (h *Handler) emitAdd(mod *Instance) {
	h.logModAdd(mod, false)

	f := mod.listeners.OnAdd
	if f != nil {
		f(mod)
	}
}

func (h *Handler) emitExtendDuration(mod *Instance) {
	h.logModAdd(mod, true)
}

func (h *Handler) emitThinkInterval(mod *Instance) {
	if mod.listeners.OnThinkInterval != nil {
		mod.listeners.OnThinkInterval(mod)
	}
}

func (h *Handler) emitStat(mod *Instance) {
	if mod.listeners.OnStat != nil {
		mod.listeners.OnStat(mod)
	}
}

func (h *Handler) emitRemove(mods []*Instance) {
	for _, mod := range mods {
		if mod.listeners.OnRemove != nil {
			mod.listeners.OnRemove(mod)
			// TODO: log
		}
	}
}

func (h *Handler) onEnemyHit(event event.TargetHitEvent) {
	for _, mod := range h.modifiers {
		if mod.listeners.OnEnemyHit != nil {
			mod.listeners.OnEnemyHit(mod, event.Target, event.AttackEvent)
		}
	}
}

func (h *Handler) onCharacterAction(event event.CharacterActionEvent) {
	for _, mod := range h.modifiers {
		if mod.listeners.OnCharacterAction != nil {
			mod.listeners.OnCharacterAction(mod, event.CharIndex, event.Action, event.Params)
		}
	}
}

func (h *Handler) logModAdd(mod *Instance, overwrote bool) {
	owner := -1
	if char := h.core.GetCharacterByTarget(mod.Owner()); char != nil {
		owner = char.GetIndex()
	}

	msg := "mod added"
	if overwrote {
		msg = "mod refreshed"
	}

	expiry := h.core.F() + mod.Duration()
	h.core.Log().NewEventBuildMsg(glog.LogStatusEvent, owner, msg).
		Write("overwrite", overwrote).
		Write("key", mod.Key()).
		Write("expiry", expiry).
		SetEnded(expiry)
}
