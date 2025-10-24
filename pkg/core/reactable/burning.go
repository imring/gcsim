package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryBurning(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}

	dendroDur := r.Durability[attributes.ElementGrass]

	// adding pyro or dendro just adds to durability
	switch a.Info.Element {
	case attributes.ElementFire:
		// if there's no existing pyro/burning or dendro/quicken then do nothing
		if r.Durability[attributes.ElementGrass] < ZeroDur && r.Durability[attributes.ElementOverdose] < ZeroDur {
			return false
		}
		// add to pyro durability
		// r.attachOrRefillNormalEle(ModifierPyro, a.Info.Durability)
	case attributes.ElementGrass:
		// if there's no existing pyro/burning or dendro/quicken then do nothing
		if r.Durability[attributes.ElementFire] < ZeroDur && r.Durability[attributes.ElementBurning] < ZeroDur {
			return false
		}
		dendroDur = max(dendroDur, a.Info.Durability*0.8)
		// add to dendro durability
		// r.attachOrRefillNormalEle(ModifierDendro, a.Info.Durability)
	default:
		return false
	}
	// a.Reacted = true

	if r.Durability[attributes.ElementBurningFuel] < ZeroDur {
		r.attachBurningFuel(max(dendroDur, r.Durability[attributes.ElementOverdose]), 1)
		r.attachBurning()

		r.core.Events().Burning.Emit(event.ReactionEvent{
			Target:      r.self.Key(),
			AttackEvent: a,
		})
		r.calcBurningDmg(a)

		if r.burningTickSrc == -1 {
			r.burningTickSrc = r.core.F()
			if t, ok := r.self.(core.Enemy); ok {
				// queue up burning ticks
				t.QueueEnemyTask(r.nextBurningTick(r.core.F(), 1, t), 15)
			}
		}
		return true
	}
	// overwrite burning fuel and recalc burning dmg
	if a.Info.Element == attributes.ElementGrass {
		r.attachBurningFuel(a.Info.Durability, 0.8)
	}
	r.calcBurningDmg(a)

	return false
}

func (r *Reactable) attachBurningFuel(dur, mult info.Durability) {
	// burning fuel always overwrites
	r.Durability[attributes.ElementBurningFuel] = mult * dur
	decayRate := mult * dur / (6*dur + 420)
	if decayRate < 10.0/60.0 {
		decayRate = 10.0 / 60.0
	}
	r.DecayRate[attributes.ElementBurningFuel] = decayRate
}

func (r *Reactable) calcBurningDmg(a *info.AttackEvent) {
	atk := info.Attack{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(info.ReactionTypeBurning),
		AttackTag:        info.AttackTagBurningDamage,
		ICDTag:           info.ICDTagBurningDamage,
		ICDGroup:         info.ICDGroupBurning,
		StrikeType:       info.StrikeTypeDefault,
		Element:          attributes.ElementFire,
		Durability:       25,
		IgnoreDefPercent: 1,
	}
	char := r.core.GetCharacter(a.Info.ActorIndex)
	flatdmg, snap := calcReactionDmg(&atk, char)
	atk.FlatDmg = 0.25 * flatdmg
	r.burningAtk = atk
	r.burningSnapshot = snap
}

func (r *Reactable) nextBurningTick(src, counter int, t core.Enemy) func() {
	return func() {
		if r.burningTickSrc != src {
			// source changed, do nothing
			return
		}
		// burning SHOULD be active still, since if not we would have
		// called cleanup and set source to -1
		if r.Durability[attributes.ElementBurningFuel] < ZeroDur || r.Durability[attributes.ElementBurning] < ZeroDur {
			return
		}
		// so burning is active, which means both auras must still have value > 0, so we can do dmg
		if counter != 9 {
			// skip the 9th tick because hyv spaghetti
			ai := r.burningAtk
			ap := combat.NewCircleHitOnTarget(r.self, nil, 1)
			r.core.QueueAttackWithSnap(r.burningSnapshot, info.QueueAttack{
				Info:    ai,
				Pattern: ap,
			})
			// self damage
			ai.Abil += info.SelfDamageSuffix
			ap.SkipTargets[info.TargettablePlayer] = false
			ap.SkipTargets[info.TargettableEnemy] = true
			ap.SkipTargets[info.TargettableGadget] = true
			r.core.QueueAttackWithSnap(r.burningSnapshot, info.QueueAttack{
				Info:    ai,
				Pattern: ap,
			})
		}
		counter++
		// queue up next tick
		t.QueueEnemyTask(r.nextBurningTick(src, counter, t), 15)
	}
}

// burningCheck purges modifiers if burning no longer active
func (r *Reactable) burningCheck() {
	if r.Durability[attributes.ElementBurning] < ZeroDur && r.Durability[attributes.ElementBurningFuel] > ZeroDur {
		// no more burning ticks
		r.burningTickSrc = -1
		// remove burning fuel
		r.Durability[attributes.ElementBurningFuel] = 0
		r.DecayRate[attributes.ElementBurningFuel] = 0
	}
}
