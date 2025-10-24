package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/model"
)

func ConvertProtoProp(s model.StatType) attributes.Prop {
	switch s {
	case model.StatType_INVALID_STAT_TYPE:
		return attributes.NoProp
	case model.StatType_FIGHT_PROP_DEFENSE_PERCENT:
		return attributes.DEFPercent
	case model.StatType_FIGHT_PROP_DEFENSE:
		return attributes.DEF
	case model.StatType_FIGHT_PROP_HP:
		return attributes.HP
	case model.StatType_FIGHT_PROP_HP_PERCENT:
		return attributes.HPPercent
	case model.StatType_FIGHT_PROP_ATTACK:
		return attributes.ATK
	case model.StatType_FIGHT_PROP_ATTACK_PERCENT:
		return attributes.ATKPercent
	case model.StatType_FIGHT_PROP_CHARGE_EFFICIENCY:
		return attributes.ER
	case model.StatType_FIGHT_PROP_ELEMENT_MASTERY:
		return attributes.EM
	case model.StatType_FIGHT_PROP_CRITICAL:
		return attributes.CR
	case model.StatType_FIGHT_PROP_CRITICAL_HURT:
		return attributes.CD
	case model.StatType_FIGHT_PROP_HEAL_ADD:
		return attributes.HealBonus
	case model.StatType_FIGHT_PROP_FIRE_ADD_HURT:
		return attributes.PyroP
	case model.StatType_FIGHT_PROP_WATER_ADD_HURT:
		return attributes.HydroP
	case model.StatType_FIGHT_PROP_GRASS_ADD_HURT:
		return attributes.DendroP
	case model.StatType_FIGHT_PROP_ELEC_ADD_HURT:
		return attributes.ElectroP
	case model.StatType_FIGHT_PROP_WIND_ADD_HURT:
		return attributes.AnemoP
	case model.StatType_FIGHT_PROP_ICE_ADD_HURT:
		return attributes.CryoP
	case model.StatType_FIGHT_PROP_ROCK_ADD_HURT:
		return attributes.GeoP
	case model.StatType_FIGHT_PROP_PHYSICAL_ADD_HURT:
		return attributes.PhyP
	case model.StatType_FIGHT_PROP_SHIELD_COST_MINUS_RATIO_ADD_HURT:
		//TODO: this is not a stat for gcsim yet
		return attributes.NoProp
	case model.StatType_FIGHT_PROP_HEALED_ADD:
		return attributes.HealTaken
	case model.StatType_FIGHT_PROP_BASE_HP:
		return attributes.HPBase
	case model.StatType_FIGHT_PROP_BASE_ATTACK:
		return attributes.ATKBase
	case model.StatType_FIGHT_PROP_BASE_DEFENSE:
		return attributes.DEFBase
	case model.StatType_FIGHT_PROP_MAX_HP:
		//TODO: this is for maxhp which is not a stat for us
		return attributes.NoProp
	default:
		return attributes.NoProp
	}
}

func ConvertProtoElement(e model.Element) attributes.ElementType {
	switch e {
	case model.Element_Electric:
		return attributes.ElementElectric
	case model.Element_Fire:
		return attributes.ElementFire
	case model.Element_Ice:
		return attributes.ElementIce
	case model.Element_Water:
		return attributes.ElementWater
	case model.Element_Grass:
		return attributes.ElementGrass
	case model.Element_Overdose:
		return attributes.ElementOverdose
	case model.Element_Frozen:
		return attributes.ElementFrozen
	case model.Element_Wind:
		return attributes.ElementWind
	case model.Element_Rock:
		return attributes.ElementRock
	default:
		return attributes.ElementNone
	}
}

func ConvertRegion(z model.ZoneType) ZoneType {
	switch z {
	case model.ZoneType_ASSOC_TYPE_MONDSTADT:
		return ZoneMondstadt
	case model.ZoneType_ASSOC_TYPE_LIYUE:
		return ZoneLiyue
	case model.ZoneType_ASSOC_TYPE_INAZUMA:
		return ZoneInazuma
	case model.ZoneType_ASSOC_TYPE_SUMERU:
		return ZoneSumeru
	case model.ZoneType_ASSOC_TYPE_FONTAINE:
		return ZoneFontaine
	case model.ZoneType_ASSOC_TYPE_NATLAN:
		return ZoneNatlan
	case model.ZoneType_ASSOC_TYPE_RANGER:
		// aloy
		return ZoneUnknown
	case model.ZoneType_ASSOC_TYPE_MAINACTOR:
		// traveler
		return ZoneUnknown
	default:
		return ZoneUnknown
	}
}

func ConvertWeaponClass(w model.WeaponClass) WeaponClass {
	switch w {
	case model.WeaponClass_WEAPON_SWORD_ONE_HAND:
		return WeaponClassSword
	case model.WeaponClass_WEAPON_CLAYMORE:
		return WeaponClassClaymore
	case model.WeaponClass_WEAPON_POLE:
		return WeaponClassSpear
	case model.WeaponClass_WEAPON_BOW:
		return WeaponClassBow
	case model.WeaponClass_WEAPON_CATALYST:
		return WeaponClassCatalyst
	default:
		//TODO: we should have an invalid?
		return WeaponClassSword
	}
}

func ConvertBodyType(b model.BodyType) BodyType {
	switch b {
	case model.BodyType_BODY_BOY:
		return BodyBoy
	case model.BodyType_BODY_GIRL:
		return BodyGirl
	case model.BodyType_BODY_MALE:
		return BodyMale
	case model.BodyType_BODY_LADY:
		return BodyLady
	case model.BodyType_BODY_LOLI:
		return BodyLoli
	default:
		//TODO: no invalid
		return BodyMale
	}
}

func ConvertRarity(q model.QualityType) int {
	switch q {
	case model.QualityType_QUALITY_ORANGE_SP:
		return 5
	case model.QualityType_QUALITY_ORANGE:
		return 5
	case model.QualityType_QUALITY_PURPLE:
		return 4
	default:
		return 4
	}
}
