package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TrySuperconduct(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	// this is for non frozen one
	if r.Durability[attributes.ElementFrozen] >= ZeroDur {
		return false
	}
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.ElementElectric:
		if r.Durability[attributes.ElementIce] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.ElementIce, a.Info.Durability, 1)
	case attributes.ElementIce:
		// could be ec potentially
		if r.Durability[attributes.ElementElectric] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.ElementElectric, a.Info.Durability, 1)
	default:
		return false
	}

	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true
	r.queueSuperconduct(a)
	return true
}

func (r *Reactable) TryFrozenSuperconduct(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	// this is for frozen
	if r.Durability[attributes.ElementFrozen] < ZeroDur {
		return false
	}
	switch a.Info.Element {
	case attributes.ElementElectric:
		//TODO: the assumption here is we first reduce cryo, and if there's any
		// src durability left, we reduce frozen. note that it's still only one
		// superconduct reaction
		a.Info.Durability -= r.reduce(attributes.ElementIce, a.Info.Durability, 1)
		r.reduce(attributes.ElementFrozen, a.Info.Durability, 1)
		a.Info.Durability = 0
		a.Reacted = true
	default:
		return false
	}

	r.queueSuperconduct(a)

	return false
}

func (r *Reactable) queueSuperconduct(a *info.AttackEvent) {
	r.core.Events().Superconduct.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// 0.1s gcd on superconduct attack
	if r.superconductGCD != -1 && r.core.F() < r.superconductGCD {
		return
	}
	r.superconductGCD = r.core.F() + 0.1*60

	// superconduct attack
	atk := info.Attack{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(info.ReactionTypeSuperconduct),
		AttackTag:        info.AttackTagSuperconductDamage,
		ICDTag:           info.ICDTagSuperconductDamage,
		ICDGroup:         info.ICDGroupReactionA,
		StrikeType:       info.StrikeTypeDefault,
		Element:          attributes.ElementIce,
		IgnoreDefPercent: 1,
	}
	char := r.core.GetCharacter(a.Info.ActorIndex)
	flatdmg, snap := calcReactionDmg(&atk, char)
	atk.FlatDmg = 1.5 * flatdmg
	r.core.QueueAttackWithSnap(snap, info.QueueAttack{
		Info:     atk,
		Pattern:  combat.NewCircleHitOnTarget(r.self, nil, 3),
		DmgDelay: 1,
	})
}
