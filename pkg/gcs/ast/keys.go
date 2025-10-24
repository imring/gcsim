package ast

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var key = map[string]TokenType{
	".":           ItemDot,
	"let":         KeywordLet,
	"while":       KeywordWhile,
	"if":          KeywordIf,
	"else":        KeywordElse,
	"fn":          KeywordFn,
	"switch":      KeywordSwitch,
	"case":        KeywordCase,
	"default":     KeywordDefault,
	"break":       KeywordBreak,
	"continue":    KeywordContinue,
	"fallthrough": KeywordFallthrough,
	"return":      KeywordReturn,
	"for":         KeywordFor,
	// genshin specific keywords
	"options":             KeywordOptions,
	"add":                 KeywordAdd,
	"char":                KeywordChar,
	"stats":               KeywordStats,
	"weapon":              KeywordWeapon,
	"set":                 KeywordSet,
	"lvl":                 KeywordLvl,
	"refine":              KeywordRefine,
	"cons":                KeywordCons,
	"talent":              KeywordTalent,
	"count":               KeywordCount,
	"params":              KeywordParams,
	"label":               KeywordLabel,
	"until":               KeywordUntil,
	"active":              KeywordActive,
	"target":              KeywordTarget,
	"particle_threshold":  KeywordParticleThreshold,
	"particle_drop_count": KeywordParticleDropCount,
	"particle_element":    KeywordParticleElement,
	"resist":              KeywordResist,
	"energy":              KeywordEnergy,
	"hurt":                KeywordHurt,
	// commands
	// team keywords
	// flags
	// ??
	// energy/hurt event related
	// target related
}

var StatKeys = map[string]attributes.Prop{
	"def%":     attributes.DEFPercent,
	"def":      attributes.DEF,
	"hp":       attributes.HP,
	"hp%":      attributes.HPPercent,
	"atk":      attributes.ATK,
	"atk%":     attributes.ATKPercent,
	"er":       attributes.ER,
	"em":       attributes.EM,
	"cr":       attributes.CR,
	"cd":       attributes.CD,
	"heal":     attributes.HealBonus,
	"pyro%":    attributes.PyroP,
	"hydro%":   attributes.HydroP,
	"cryo%":    attributes.CryoP,
	"electro%": attributes.ElectroP,
	"anemo%":   attributes.AnemoP,
	"geo%":     attributes.GeoP,
	"phys%":    attributes.PhyP,
	// "ele%":     attributes.ElementalP,
	"dendro%": attributes.DendroP,
	"atkspd%": attributes.AtkSpd,
	"dmg%":    attributes.DmgP,
}

var EleKeys = map[string]attributes.ElementType{
	"electro":  attributes.ElementElectric,
	"pyro":     attributes.ElementFire,
	"cryo":     attributes.ElementIce,
	"hydro":    attributes.ElementWater,
	"frozen":   attributes.ElementFrozen,
	"anemo":    attributes.ElementWind,
	"dendro":   attributes.ElementGrass,
	"geo":      attributes.ElementRock,
	"physical": attributes.ElementNone, // TODO: is there physical element?
	"none":     attributes.ElementNone,
}

var actionKeys = map[string]info.Action{
	"skill":       info.ActionSkill,
	"burst":       info.ActionBurst,
	"attack":      info.ActionAttack,
	"charge":      info.ActionCharge,
	"high_plunge": info.ActionHighPlunge,
	"low_plunge":  info.ActionLowPlunge,
	"aim":         info.ActionAim,
	"dash":        info.ActionDash,
	"jump":        info.ActionJump,
	"walk":        info.ActionWalk,
	"swap":        info.ActionSwap,
}
