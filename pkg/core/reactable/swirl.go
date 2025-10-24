package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func calcSwirlAtkDurability(consumed, src info.Durability) info.Durability {
	if consumed < src {
		return 1.25*(0.5*consumed-1) + 25
	}
	return 1.25*(src-1) + 25
}

func (r *Reactable) queueSwirl(rt info.ReactionType, ele attributes.ElementType, tag info.AttackTag, icd info.ICDTag, dur info.Durability, charIndex int) {
	// swirl triggers two attacks; one self with no gauge
	// and one aoe with gauge
	ai := info.Attack{
		ActorIndex:       charIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(rt),
		AttackTag:        tag,
		ICDTag:           icd,
		ICDGroup:         info.ICDGroupReactionA,
		StrikeType:       info.StrikeTypeDefault,
		Element:          ele,
		IgnoreDefPercent: 1,
	}
	char := r.core.GetCharacter(charIndex)
	flatdmg, snap := calcReactionDmg(&ai, char)
	ai.FlatDmg = 0.6 * flatdmg
	// first attack is self no hitbox
	r.core.QueueAttackWithSnap(snap, info.QueueAttack{
		Info:     ai,
		Pattern:  combat.NewSingleTargetHit(r.self.Key()),
		DmgDelay: 1,
	})
	// next is aoe - hydro swirls never do AoE damage, as they only spread the element
	if ele == attributes.ElementWater {
		ai.FlatDmg = 0
	}
	ai.Durability = dur
	ai.Abil = string(rt) + " (aoe)"
	ap := combat.NewCircleHitOnTarget(r.self, nil, 5)
	ap.IgnoredKeys = []keys.Target{r.self.Key()}
	r.core.QueueAttackWithSnap(snap, info.QueueAttack{
		Info:     ai,
		Pattern:  ap,
		DmgDelay: 1 + 4, // 4f after self swirl
	})
}

func (r *Reactable) TrySwirlElectro(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.Durability[attributes.ElementElectric] < ZeroDur {
		return false
	}
	rd := r.reduce(attributes.ElementElectric, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	// queue an attack first
	r.core.Events().SwirlElectro.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// 0.1s gcd on swirl electro attack
	if !(r.swirlElectroGCD != -1 && r.core.F() < r.swirlElectroGCD) {
		r.swirlElectroGCD = r.core.F() + 0.1*60
		r.queueSwirl(
			info.ReactionTypeSwirlElectro,
			attributes.ElementElectric,
			info.AttackTagSwirlElectro,
			info.ICDTagSwirlElectro,
			atkDur,
			a.Info.ActorIndex,
		)
	}

	// at this point if any durability left, we need to check for prescence of
	// hydro in case of EC
	if a.Info.Durability > ZeroDur && r.Durability[attributes.ElementWater] > ZeroDur {
		// trigger swirl hydro
		r.TrySwirlHydro(a)
		// check EC clean up
		r.checkEC()
	}
	return true
}

func (r *Reactable) TrySwirlHydro(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.Durability[attributes.ElementWater] < ZeroDur {
		return false
	}
	rd := r.reduce(attributes.ElementWater, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	// queue an attack first
	r.core.Events().SwirlHydro.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// 0.1s gcd on swirl hydro attack
	if !(r.swirlHydroGCD != -1 && r.core.F() < r.swirlHydroGCD) {
		r.swirlHydroGCD = r.core.F() + 0.1*60
		r.queueSwirl(
			info.ReactionTypeSwirlHydro,
			attributes.ElementWater,
			info.AttackTagSwirlHydro,
			info.ICDTagSwirlHydro,
			atkDur,
			a.Info.ActorIndex,
		)
	}

	return true
}

func (r *Reactable) TrySwirlCryo(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.Durability[attributes.ElementIce] < ZeroDur {
		return false
	}
	rd := r.reduce(attributes.ElementIce, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	// queue an attack first
	r.core.Events().SwirlCryo.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// 0.1s gcd on swirl cryo attack
	if !(r.swirlCryoGCD != -1 && r.core.F() < r.swirlCryoGCD) {
		r.swirlCryoGCD = r.core.F() + 0.1*60
		r.queueSwirl(
			info.ReactionTypeSwirlCryo,
			attributes.ElementIce,
			info.AttackTagSwirlCryo,
			info.ICDTagSwirlCryo,
			atkDur,
			a.Info.ActorIndex,
		)
	}

	return true
}

func (r *Reactable) TrySwirlPyro(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.Durability[attributes.ElementFire] < ZeroDur {
		return false
	}
	rd := r.reduce(attributes.ElementFire, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	r.burningCheck()
	// queue an attack first
	r.core.Events().SwirlPyro.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// 0.1s gcd on swirl pyro attack
	if !(r.swirlPyroGCD != -1 && r.core.F() < r.swirlPyroGCD) {
		r.swirlPyroGCD = r.core.F() + 0.1*60
		r.queueSwirl(
			info.ReactionTypeSwirlPyro,
			attributes.ElementFire,
			info.AttackTagSwirlPyro,
			info.ICDTagSwirlPyro,
			atkDur,
			a.Info.ActorIndex,
		)
	}

	return true
}

func (r *Reactable) TrySwirlFrozen(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.Durability[attributes.ElementFrozen] < ZeroDur {
		return false
	}
	rd := r.reduce(attributes.ElementFrozen, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	// queue an attack first
	r.core.Events().SwirlCryo.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})
	// 0.1s gcd on swirl cryo attack
	if !(r.swirlCryoGCD != -1 && r.core.F() < r.swirlCryoGCD) {
		r.swirlCryoGCD = r.core.F() + 0.1*60
		r.queueSwirl(
			info.ReactionTypeSwirlCryo,
			attributes.ElementIce,
			info.AttackTagSwirlCryo,
			info.ICDTagSwirlCryo,
			atkDur,
			a.Info.ActorIndex,
		)
	}

	r.checkFreeze()
	return true
}
