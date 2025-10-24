package enemy

import (
	"math"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/damage"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (e *Enemy) HandleAttack(atk *info.AttackEvent) float64 {
	e.Core().Events().EnemyHit.Emit(event.TargetHitEvent{
		Target:      e.Key(),
		AttackEvent: atk,
	})

	var amp string
	var cata string
	var dmg float64
	var grpMult float64
	var crit bool

	evt := e.Core().Log().NewEvent(atk.Info.Abil, glog.LogDamageEvent, atk.Info.ActorIndex).
		Write("target", e.Key()).
		Write("attack-tag", atk.Info.AttackTag).
		Write("ele", atk.Info.Element.String()).
		Write("damage", &dmg).
		Write("crit", &crit).
		Write("damage_grp_mult", &grpMult).
		Write("amp", &amp).
		Write("cata", &cata).
		Write("abil", atk.Info.Abil).
		Write("source_frame", atk.SourceFrame)
	evt.WriteBuildMsg(atk.Attacker.Logs...)

	if !atk.Info.SourceIsSim {
		logDetails := make([]any, 0, 5)
		for _, v := range atk.Attacker.Stats.Changes {
			modStatus := make([]string, 0, 2)
			modStatus = append(
				modStatus,
				"status: added",
				"expiry_frame: "+strconv.Itoa(v.Expiry),
			)
			modStatus = append(
				modStatus,
				attributes.PrettyPrintStatsMap(v.PropMap)...,
			)
			logDetails = append(logDetails, v.Reason, modStatus)
		}
		evt.Write("pre_damage_mods", logDetails)
	}

	e.tryReact(atk)

	dmg, crit = damage.Calc(atk, e.Core().Rand(), e.Core().Log())

	// reduce damage by damage group
	grpMult = 1.0
	if !atk.Info.SourceIsSim {
		grpMult = e.GroupTagDamageMult(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex)
		dmg *= grpMult
	}

	e.checkHitlag(atk)

	// TODO: particle drops

	// delay damage event to end of the frame
	e.Core().Tasks().Add(func() {
		// apply the damage
		actualDmg := e.applyDamage(atk, dmg)
		// e.Core.Combat.TotalDamage += actualDmg
		e.Core().Events().EnemyDamage.Emit(event.EnemyDamageEvent{
			Target:      e.Key(),
			AttackEvent: atk,
			Damage:      actualDmg,
			IsCrit:      crit,
		})

		// callbacks
		cb := info.AttackCB{
			Target:      e.Key(),
			AttackEvent: atk,
			Damage:      actualDmg,
			IsCrit:      crit,
		}
		for _, f := range atk.Callbacks {
			f(cb)
		}
	}, 0)

	// this works because string in golang is a slice underneath, so the &amp points to the slice info
	// that's why when the underlying string in amp changes (has to be reallocated) the pointer doesn't
	// change since it's just pointing to the slice "header"
	if atk.Info.Amped {
		amp = string(atk.Info.AmpType)
	}
	if atk.Info.Catalyzed {
		cata = string(atk.Info.CatalyzedType)
	}

	return dmg
}

func (e *Enemy) ModifyHP(amount float64) {
	e.hpRatio += amount / e.MaxHP()
	e.hpRatio = max(0, min(1, e.hpRatio))
}

func (e *Enemy) tryReact(atk *info.AttackEvent) {
	// if target is frozen prior to attack landing, set impulse to 0
	// let the break freeze attack to trigger actual impulse
	if e.GetAuraDurability(attributes.ElementFrozen) > info.ZeroDurability {
		atk.Info.NoImpulse = true
	}

	// check poise dmg and then shatter first
	e.PoiseDMGCheck(atk)
	e.ShatterCheck(atk)

	checkBurningICD := func() {
		// special global ICD for Burning DMG
		if atk.Info.ICDTag != info.ICDTagBurningDamage {
			return
		}
		// checks for ICD on all the other characters as well
		// TODO: chars
		// for i := 0; i < len(e.Core.Player.Chars()); i++ {
		// 	if i == atk.Info.ActorIndex {
		// 		continue
		// 	}
		// 	// burning durability wiped out to 0 if any of the other char still on icd re burning dmg
		// 	atk.Info.Durability *= info.Durability(e.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, i))
		// }
	}

	// check tags
	if atk.Info.Durability > 0 {
		// check for ICD first
		atk.Info.Durability *= info.Durability(e.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex))
		checkBurningICD()
		if atk.Info.Durability > 0 && atk.Info.Element != attributes.ElementNone {
			existing := e.Reactable.ActiveAuraString()
			applied := atk.Info.Durability
			e.React(atk)
			if atk.Reacted {
				e.Core().Log().NewEvent(
					"application",
					glog.LogElementEvent,
					atk.Info.ActorIndex,
				).
					Write("attack_tag", atk.Info.AttackTag).
					Write("applied_ele", atk.Info.Element.String()).
					Write("dur", applied).
					Write("abil", atk.Info.Abil).
					Write("target", e.Key()).
					Write("existing", existing).
					Write("after", e.Reactable.ActiveAuraString())
			}
		}
	}
}

func (e *Enemy) checkHitlag(atk *info.AttackEvent) {
	if !e.Core().IsHitlagEnabled() {
		return
	}

	willapply := true
	if atk.Info.HitlagOnHeadshotOnly {
		willapply = atk.Info.HitWeakPoint
	}
	dur := atk.Info.HitlagHaltFrames
	if e.Core().CanBeDefenseHalt() && atk.Info.CanBeDefenseHalted {
		dur += 3.6
	}
	dur = math.Ceil(dur)
	if willapply && dur > 0 {
		// apply hit lag to enemy
		e.ApplyHitlag(atk.Info.HitlagFactor, dur)
		// also apply hitlag to reactable
		// e.Reactable.ApplyHitlag(atk.Info.HitlagFactor, dur)
	}
}

func (e *Enemy) applyDamage(atk *info.AttackEvent, damage float64) float64 {
	// record dmg
	// do not let hp become negative because this function can be called multiple times in same frame
	actualDmg := min(damage, e.MaxHP()) // do not let dmg be greater than remaining enemy hp
	e.ModifyHP(-actualDmg)

	// check if target is dead
	if e.Core().IsDamageMode() && e.HP() <= 0 {
		e.Kill(atk)
		return actualDmg
	}

	// apply auras
	if atk.Info.Durability > 0 && !atk.Reacted && atk.Info.Element != attributes.ElementNone {
		// check for ICD first
		existing := e.Reactable.ActiveAuraString()
		applied := atk.Info.Durability
		e.AttachOrRefill(atk)
		e.Core().Log().NewEvent(
			"application",
			glog.LogElementEvent,
			atk.Info.ActorIndex,
		).
			Write("attack_tag", atk.Info.AttackTag).
			Write("applied_ele", atk.Info.Element.String()).
			Write("dur", applied).
			Write("abil", atk.Info.Abil).
			Write("target", e.Key()).
			Write("existing", existing).
			Write("after", e.Reactable.ActiveAuraString())
	}
	// just return damage without considering enemy hp here for both:
	// - damage mode if target not dead (otherwise would have entered the death if statement)
	// - duration mode (no concept of killing blow)
	return damage
}
