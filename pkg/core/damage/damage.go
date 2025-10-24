package damage

import (
	"math/rand"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func Calc(atk *info.AttackEvent, rand *rand.Rand, log glog.Logger) (float64, bool) {
	var isCrit bool

	attackerProps := atk.Attacker.Stats.Props
	targetProps := atk.Target.Stats.Props

	// calculate using either attack or def or hp
	var a float64
	switch {
	case atk.Info.UseHP:
		a = attackerProps.MaxHP()
	case atk.Info.UseDef:
		a = attackerProps.TotalDEF()
	default:
		a = attackerProps.TotalATK()
	}

	base := atk.Info.Mult*a + atk.Info.FlatDmg
	dmgBonus := attackerProps.DMG(atk.Info.Element)
	damage := base * (1 + dmgBonus)

	// def
	defadj := max(-0.9, targetProps[attributes.DEFPercent])
	defignore := atk.Info.IgnoreDefPercent + attackerProps[attributes.DEFIgnore]
	def := min(targetProps.TotalDEF(), 0) * max(1-defignore, 0)
	defmod := float64(5*atk.Attacker.Level+500) / (def + 5*float64(atk.Attacker.Level+500))
	damage *= defmod

	// res
	res := targetProps.Resist(atk.Info.Element)
	resmod := 1 - res/2
	if res >= 0 && res < 0.75 {
		resmod = 1 - res
	} else if res > 0.75 {
		resmod = 1 / (4*res + 1)
	}
	damage *= resmod

	// crit
	precritdmg := damage
	critRate := max(min(attackerProps.CritRate(), 1), 0)
	critDmg := attackerProps.CritDMG()
	if atk.Info.HitWeakPoint || rand.Float64() <= critRate {
		damage *= 1 + critDmg
		isCrit = true
	}

	// em bonus
	preampdmg := damage
	em := attackerProps.EM()
	emBonus := (2.78 * em) / (1400 + em)
	var reactBonus float64
	var ampBonus float64
	if atk.Info.Amped {
		reactBonus = ReactionBonus(&atk.Info, attackerProps)
		ampBonus = atk.Info.AmpMult * (1 + emBonus + reactBonus)
		damage *= ampBonus
	}

	if log != nil {
		avgDmg := precritdmg * (1 + critRate*critDmg)

		log.NewEvent(
			atk.Info.Abil,
			glog.LogCalc,
			atk.Info.ActorIndex,
		).
			Write("src_frame", atk.SourceFrame).
			Write("damage", damage).
			Write("abil", atk.Info.Abil).
			Write("talent", atk.Info.Mult).
			Write("base_atk", attackerProps[attributes.ATKBase]).
			Write("flat_atk", attackerProps[attributes.ATK]).
			Write("atk_per", attackerProps[attributes.ATKPercent]).
			Write("use_def", atk.Info.UseDef).
			Write("base_def", attackerProps[attributes.DEFBase]).
			Write("flat_def", attackerProps[attributes.DEF]).
			Write("def_per", attackerProps[attributes.DEFPercent]).
			Write("base_hp", attackerProps[attributes.HPBase]).
			Write("flat_hp", attackerProps[attributes.HP]).
			Write("hp_per", attackerProps[attributes.HPPercent]).
			Write("catalyzed", atk.Info.Catalyzed).
			Write("flat_dmg", atk.Info.FlatDmg).
			Write("total_atk_def", a).
			Write("base_dmg", base).
			Write("ele", atk.Info.Element).
			Write("ele_per", dmgBonus-attackerProps[attributes.DmgP]).
			Write("bonus_dmg", dmgBonus).
			Write("ignore_def", atk.Info.IgnoreDefPercent).
			Write("def_adj", defadj).
			Write("target_lvl", atk.Target.Level).
			Write("char_lvl", atk.Attacker.Level).
			Write("def_mod", defmod).
			Write("res", res).
			Write("res_mod", resmod).
			Write("cr", critRate).
			Write("cd", critDmg).
			Write("pre_crit_dmg", precritdmg).
			Write("dmg_if_crit", precritdmg*(1+critDmg)).
			Write("avg_crit_dmg", avgDmg).
			Write("is_crit", isCrit).
			Write("pre_amp_dmg", preampdmg).
			Write("reaction_type", atk.Info.AmpType).
			Write("melt_vape", atk.Info.Amped).
			Write("react_mult", atk.Info.AmpMult).
			Write("em", em).
			Write("em_bonus", emBonus).
			Write("react_bonus", reactBonus).
			Write("amp_mult_total", ampBonus).
			Write("pre_crit_dmg_react", precritdmg*ampBonus).
			Write("dmg_if_crit_react", precritdmg*(1+critDmg)*ampBonus).
			Write("avg_crit_dmg_react", avgDmg*ampBonus)
	}

	return damage, isCrit
}

// leave as a function because some props have extra version
func ReactionBonus(atk *info.Attack, attackerProps attributes.Props) float64 {
	// amplifying reactions
	if atk.Amped {
		switch atk.AmpType {
		case info.ReactionTypeMelt:
			return attackerProps[attributes.MeltBonus]
		case info.ReactionTypeVaporize:
			return attackerProps[attributes.VaporizeBonus]
		}
	}

	// additive reactions
	if atk.Catalyzed {
		switch atk.CatalyzedType {
		case info.ReactionTypeAggravate:
			return attackerProps[attributes.AggravateBonus]
		case info.ReactionTypeSpread:
			return attackerProps[attributes.SpreadBonus]
		}
	}

	// transformative reactions
	switch atk.AttackTag {
	case info.AttackTagOverloadDamage:
		return attackerProps[attributes.OverloadBonus]
	case info.AttackTagSuperconductDamage:
		return attackerProps[attributes.SuperconductBonus]
	case info.AttackTagECDamage:
		return attackerProps[attributes.ElectroChargedBonus]
	case info.AttackTagShatter:
		return attackerProps[attributes.ShatterBonus]
	case info.AttackTagSwirlPyro:
		return attackerProps[attributes.SwirlPyroBonus]
	case info.AttackTagSwirlHydro:
		return attackerProps[attributes.SwirlHydroBonus]
	case info.AttackTagSwirlCryo:
		return attackerProps[attributes.SwirlCryoBonus]
	case info.AttackTagSwirlElectro:
		return attackerProps[attributes.SwirlElectroBonus]
	case info.AttackTagBurningDamage:
		return attackerProps[attributes.BurningBonus]
	case info.AttackTagBloom:
		return attackerProps[attributes.BloomBonus]
	case info.AttackTagBountifulCore:
		return attackerProps[attributes.BountifulCoreBonus]
	case info.AttackTagBurgeon:
		return attackerProps[attributes.BurgeonBonus]
	case info.AttackTagHyperbloom:
		return attackerProps[attributes.HyperbloomBonus]
	}

	return 0
}
