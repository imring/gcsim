package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryOverload(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.ElementElectric:
		// must have pyro; pyro cant coexist (for now) so ok to ignore count?
		if r.Durability[attributes.ElementFire] < ZeroDur && r.Durability[attributes.ElementBurning] < ZeroDur {
			return false
		}
		// reduce; either gone or left; don't care how much actually reacted
		consumed = r.reduce(attributes.ElementFire, a.Info.Durability, 1)
		r.burningCheck()
	case attributes.ElementFire:
		// must have electro; gotta be careful with ec?
		if r.Durability[attributes.ElementElectric] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.ElementElectric, a.Info.Durability, 1)
	default:
		// should be here
		return false
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	// trigger event before attack is queued. this gives time for other actions to modify it
	r.core.Events().Overload.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// 0.1s gcd on overload attack
	if !(r.overloadGCD != -1 && r.core.F() < r.overloadGCD) {
		r.overloadGCD = r.core.F() + 0.1*60
		// trigger an overload attack
		atk := info.Attack{
			ActorIndex:       a.Info.ActorIndex,
			DamageSrc:        r.self.Key(),
			Abil:             string(info.ReactionTypeOverload),
			AttackTag:        info.AttackTagOverloadDamage,
			ICDTag:           info.ICDTagOverloadDamage,
			ICDGroup:         info.ICDGroupReactionB,
			StrikeType:       info.StrikeTypeBlunt,
			PoiseDMG:         90,
			Element:          attributes.ElementFire,
			IgnoreDefPercent: 1,
		}
		char := r.core.GetCharacter(a.Info.ActorIndex)
		flatdmg, snap := calcReactionDmg(&atk, char)
		atk.FlatDmg = 2.75 * flatdmg
		r.core.QueueAttackWithSnap(snap, info.QueueAttack{
			Info:     atk,
			Pattern:  combat.NewCircleHitOnTarget(r.self, nil, 3),
			DmgDelay: 1,
		})
	}

	return true
}
