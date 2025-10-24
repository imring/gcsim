package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryAddEC(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	// if there's still frozen left don't try to ec
	// game actively rejects ec reaction if frozen is present
	if r.Durability[attributes.ElementFrozen] > ZeroDur {
		return false
	}

	// adding ec or hydro just adds to durability
	switch a.Info.Element {
	case attributes.ElementWater:
		// if there's no existing hydro or electro then do nothing
		if r.Durability[attributes.ElementElectric] < ZeroDur {
			return false
		}
		// add to hydro durability (can't add if the atk already reacted)
		//TODO: this shouldn't happen here
		if !a.Reacted {
			r.attachOrRefillNormalEle(attributes.ElementWater, a.Info.Durability)
		}
	case attributes.ElementElectric:
		// if there's no existing hydro or electro then do nothing
		if r.Durability[attributes.ElementWater] < ZeroDur {
			return false
		}
		// add to electro durability (can't add if the atk already reacted)
		if !a.Reacted {
			r.attachOrRefillNormalEle(attributes.ElementElectric, a.Info.Durability)
		}
	default:
		return false
	}

	a.Reacted = true
	r.core.Events().ElectroCharged.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// at this point ec is refereshed so we need to trigger a reaction
	// and change ownership
	atk := info.Attack{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(info.ReactionTypeElectroCharged),
		AttackTag:        info.AttackTagECDamage,
		ICDTag:           info.ICDTagECDamage,
		ICDGroup:         info.ICDGroupReactionB,
		StrikeType:       info.StrikeTypeDefault,
		Element:          attributes.ElementElectric,
		IgnoreDefPercent: 1,
	}
	char := r.core.GetCharacter(a.Info.ActorIndex)
	flatdmg, snap := calcReactionDmg(&atk, char)
	atk.FlatDmg = 2.0 * flatdmg
	r.ecAtk = atk
	r.ecSnapshot = snap

	// if this is a new ec then trigger tick immediately and queue up ticks
	// otherwise do nothing
	//TODO: need to check if refresh ec triggers new tick immediately or not
	if r.ecTickSrc == -1 {
		r.ecTickSrc = r.core.F()
		r.core.QueueAttackWithSnap(r.ecSnapshot, info.QueueAttack{
			Info:     r.ecAtk,
			Pattern:  combat.NewSingleTargetHit(r.self.Key()),
			DmgDelay: 10,
		})

		r.core.Tasks().Add(r.nextTick(r.core.F()), 60+10)
		// subscribe to wane ticks
		r.ecEventKey = r.core.Events().EnemyDamage.Subscribe(func(event event.EnemyDamageEvent) {
			// target should be first, then snapshot
			// TODO: there's no target index
			if event.Target != r.self.Key() {
				return
			}
			if a.Info.AttackTag != info.AttackTagECDamage {
				return
			}
			// ignore if this dmg instance has been wiped out due to icd
			if event.Damage == 0 {
				return
			}
			// ignore if we no longer have both electro and hydro
			if r.Durability[attributes.ElementElectric] < ZeroDur || r.Durability[attributes.ElementWater] < ZeroDur {
				r.core.Events().EnemyDamage.Unsubscribe(r.ecEventKey)
				r.ecEventKey = -1
				return
			}

			// wane in 0.1 seconds
			r.core.Tasks().Add(func() {
				r.waneEC()
			}, 6)
		})
	}

	// ticks are 60 frames since last tick
	// taking tick dmg resets last tick
	return true
}

func (r *Reactable) waneEC() {
	r.Durability[attributes.ElementElectric] -= 10
	r.Durability[attributes.ElementElectric] = max(0, r.Durability[attributes.ElementElectric])
	r.Durability[attributes.ElementWater] -= 10
	r.Durability[attributes.ElementWater] = max(0, r.Durability[attributes.ElementWater])
	r.core.Log().NewEvent("ec wane",
		glog.LogElementEvent,
		-1,
	).
		Write("aura", "ec").
		Write("target", r.self.Key()).
		Write("hydro", r.Durability[attributes.ElementWater]).
		Write("electro", r.Durability[attributes.ElementElectric])

	// ec is gone
	r.checkEC()
}

func (r *Reactable) checkEC() {
	if r.Durability[attributes.ElementElectric] < ZeroDur || r.Durability[attributes.ElementWater] < ZeroDur {
		r.core.Events().EnemyDamage.Unsubscribe(r.ecEventKey)
		r.ecTickSrc = -1
		r.ecEventKey = -1
		r.core.Log().NewEvent("ec expired",
			glog.LogElementEvent,
			-1,
		).
			Write("aura", "ec").
			Write("target", r.self.Key()).
			Write("hydro", r.Durability[attributes.ElementWater]).
			Write("electro", r.Durability[attributes.ElementElectric])
	}
}

func (r *Reactable) nextTick(src int) func() {
	return func() {
		if r.ecTickSrc != src {
			// source changed, do nothing
			return
		}
		// ec SHOULD be active still, since if not we would have
		// called cleanup and set source to -1
		if r.Durability[attributes.ElementElectric] < ZeroDur || r.Durability[attributes.ElementWater] < ZeroDur {
			return
		}

		// so ec is active, which means both aura must still have value > 0; so we can do dmg
		r.core.QueueAttackWithSnap(r.ecSnapshot, info.QueueAttack{
			Info:    r.ecAtk,
			Pattern: combat.NewSingleTargetHit(r.self.Key()),
		})

		// queue up next tick
		r.core.Tasks().Add(r.nextTick(src), 60)
	}
}
