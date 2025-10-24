package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryFreeze(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	// so if already frozen there are 2 cases:
	// 1. src exists but no other coexisting -> attach
	// 2. src does not exist but opposite coexists -> add to freeze durability
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.ElementWater:
		// if cryo exists we'll trigger freeze regardless if frozen already coexists
		if r.Durability[attributes.ElementIce] < ZeroDur {
			return false
		}
		consumed = r.triggerFreeze(r.Durability[attributes.ElementIce], a.Info.Durability)
		r.Durability[attributes.ElementIce] -= consumed
		r.Durability[attributes.ElementIce] = max(r.Durability[attributes.ElementIce], 0)
	case attributes.ElementIce:
		if r.Durability[attributes.ElementWater] < ZeroDur {
			return false
		}
		consumed := r.triggerFreeze(r.Durability[attributes.ElementWater], a.Info.Durability)
		r.Durability[attributes.ElementWater] -= consumed
		r.Durability[attributes.ElementWater] = max(r.Durability[attributes.ElementWater], 0)
	default:
		// should be here
		return false
	}
	a.Reacted = true
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	r.core.Events().Frozen.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})
	return true
}

func (r *Reactable) PoiseDMGCheck(a *info.AttackEvent) bool {
	if r.Durability[attributes.ElementFrozen] < ZeroDur {
		return false
	}
	if a.Info.StrikeType != info.StrikeTypeBlunt {
		return false
	}
	// remove frozen durability according to poise dmg
	r.Durability[attributes.ElementFrozen] -= info.Durability(0.15 * a.Info.PoiseDMG)
	r.checkFreeze()
	return true
}

func (r *Reactable) ShatterCheck(a *info.AttackEvent) bool {
	if r.Durability[attributes.ElementFrozen] < ZeroDur {
		return false
	}
	if a.Info.StrikeType != info.StrikeTypeBlunt && a.Info.Element != attributes.ElementRock {
		return false
	}
	// remove 200 freeze gauge if available
	r.Durability[attributes.ElementFrozen] -= 200
	r.checkFreeze()

	r.core.Events().Shatter.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// 0.2s gcd on shatter attack
	if !(r.shatterGCD != -1 && r.core.F() < r.shatterGCD) {
		r.shatterGCD = r.core.F() + 0.2*60
		// trigger shatter attack
		ai := info.Attack{
			ActorIndex:       a.Info.ActorIndex,
			DamageSrc:        r.self.Key(),
			Abil:             string(info.ReactionTypeShatter),
			AttackTag:        info.AttackTagShatter,
			ICDTag:           info.ICDTagShatter,
			ICDGroup:         info.ICDGroupReactionA,
			StrikeType:       info.StrikeTypeDefault,
			IgnoreDefPercent: 1,
		}
		char := r.core.GetCharacter(a.Info.ActorIndex)
		flatdmg, snap := calcReactionDmg(&ai, char)
		ai.FlatDmg = 3.0 * flatdmg
		// shatter is a self attack
		r.core.QueueAttackWithSnap(snap, info.QueueAttack{
			Info:    ai,
			Pattern: combat.NewSingleTargetHit(r.self.Key()),
		})
	}

	return true
}

// add to freeze durability and return amount of durability consumed
func (r *Reactable) triggerFreeze(a, b info.Durability) info.Durability {
	d := min(a, b)
	if r.FreezeResist >= 1 {
		return d
	}
	// trigger freeze should only addDurability and should not touch decay rate
	r.attachOverlap(attributes.ElementFrozen, 2*d, ZeroDur)
	return d
}

func (r *Reactable) checkFreeze() {
	if r.Durability[attributes.ElementFrozen] <= ZeroDur {
		r.Durability[attributes.ElementFrozen] = 0

		r.core.Events().AuraDurabilityDepleted.Emit(event.AuraDurabilityDepletedEvent{
			Target:  r.self.Key(),
			Element: attributes.ElementFrozen,
		})
		// trigger another attack here, purely for the purpose of breaking bubbles >.>
		ai := info.Attack{
			ActorIndex:  0,
			DamageSrc:   r.self.Key(),
			Abil:        "Freeze Broken",
			AttackTag:   info.AttackTagNone,
			ICDTag:      info.ICDTagNone,
			ICDGroup:    info.ICDGroupDefault,
			StrikeType:  info.StrikeTypeDefault,
			SourceIsSim: true,
			DoNotLog:    true,
		}
		//TODO: delay attack by 1 frame ok?
		r.core.QueueAttack(info.QueueAttack{
			Info:          ai,
			Pattern:       combat.NewSingleTargetHit(r.self.Key()),
			SnapshotDelay: -1,
		})
	}
}

func (r *Reactable) SetFreezeResist(resist float64) {
	r.FreezeResist = resist
}
