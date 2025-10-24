package attributes

import (
	"strconv"
	"strings"
)

type (
	Prop    int
	Props   [EndPropType]float64
	PropMap map[Prop]float64
)

// stat types
const (
	NoProp Prop = iota

	// HP = HPBase * (1 + HPPercent) + HP + HPExtra
	HPBase
	HPPercent
	HP
	HPExtra

	// DEF = DEFBase * (1 + DEFPercent) + DEF + DEFExtra
	DEFBase
	DEFPercent
	DEF
	DEFExtra

	// ATK = ATKBase * (1 + ATKPercent) + ATK + ATKExtra
	ATKBase
	ATKPercent
	ATK
	ATKExtra

	// Crit
	CR
	CRExtra
	CD
	CDExtra

	// Energy
	ER
	ERExtra

	// Ignores target defense
	DEFIgnore

	// Elemental Mastery
	EM
	EMExtra

	// Healing Bonus
	HealBonus
	HealTaken

	// Cooldown Reduction
	CooldownR

	// DMG% Bonus = DmgP + (element)P + (element)PExtra
	DmgP
	PyroP
	PyroPExtra
	HydroP
	HydroPExtra
	CryoP
	CryoPExtra
	ElectroP
	ElectroPExtra
	AnemoP
	AnemoPExtra
	GeoP
	GeoPExtra
	DendroP
	DendroPExtra
	PhyP
	PhyPExtra

	// DMG% Resistance = DmgRes + (element)Res + (element)ResExtra
	DmgRes
	PyroRes
	PyroResExtra
	HydroRes
	HydroResExtra
	CryoRes
	CryoResExtra
	ElectroRes
	ElectroResExtra
	AnemoRes
	AnemoResExtra
	GeoRes
	GeoResExtra
	DendroRes
	DendroResExtra
	PhyRes
	PhyResExtra

	// Reaction Bonus
	OverloadBonus
	SuperconductBonus
	MeltBonus
	VaporizeBonus
	SwirlElectroBonus
	SwirlHydroBonus
	SwirlPyroBonus
	SwirlCryoBonus
	ElectroChargedBonus
	ShatterBonus
	BurningBonus
	AggravateBonus
	SpreadBonus
	BloomBonus
	BountifulCoreBonus // special stat for nilou
	BurgeonBonus
	HyperbloomBonus

	// Speed
	OverallSpd
	OverallSpdMult
	// Attack Speed = (OverallSpd + AtkSpd) * OverallSpdMult
	AtkSpd
	// Move Speed = (OverallSpd + MoveSpd) * OverallSpdMult
	MoveSpd

	// Stamina
	CostStamP

	// Shield
	ShieldStrength

	// delim
	EndPropType
)

func (p Prop) String() string {
	return PropTypeString[p]
}

func (p Props) MaxHP() float64    { return p[HPBase]*(1+p[HPPercent]) + p[HP] + p[HPExtra] }
func (p Props) TotalATK() float64 { return p[ATKBase]*(1+p[ATKPercent]) + p[ATK] + p[ATKExtra] }
func (p Props) TotalDEF() float64 {
	percent := max(-0.9, p[DEFPercent])
	return p[DEFBase]*(1+percent) + p[DEF] + p[DEFExtra]
}

func (p Props) CritRate() float64 { return p[CR] + p[CRExtra] }
func (p Props) CritDMG() float64  { return p[CD] + p[CDExtra] }

func (p Props) ER() float64 { return p[ER] + p[ERExtra] }
func (p Props) EM() float64 { return p[EM] + p[EMExtra] }

func (p Props) DMG(ele ElementType) float64 {
	result := p[DmgP]
	switch ele {
	case ElementWind:
		result += p[AnemoP] + p[AnemoPExtra]
	case ElementIce:
		result += p[CryoP] + p[CryoPExtra]
	case ElementElectric:
		result += p[ElectroP] + p[ElectroPExtra]
	case ElementRock:
		result += p[GeoP] + p[GeoPExtra]
	case ElementWater:
		result += p[HydroP] + p[HydroPExtra]
	case ElementFire:
		result += p[PyroP] + p[PyroPExtra]
	case ElementGrass:
		result += p[DendroP] + p[DendroPExtra]
	default: // TODO: is there physical element?
		result += p[PhyP] + p[PhyPExtra]
	}
	return result
}

func (p Props) Resist(ele ElementType) float64 {
	result := p[DmgRes]
	switch ele {
	case ElementWind:
		result += p[AnemoRes] + p[AnemoResExtra]
	case ElementIce:
		result += p[CryoRes] + p[CryoResExtra]
	case ElementElectric:
		result += p[ElectroRes] + p[ElectroResExtra]
	case ElementRock:
		result += p[GeoRes] + p[GeoResExtra]
	case ElementWater:
		result += p[HydroRes] + p[HydroResExtra]
	case ElementFire:
		result += p[PyroRes] + p[PyroResExtra]
	case ElementGrass:
		result += p[DendroRes] + p[DendroResExtra]
	default: // TODO: is there physical element?
		result += p[PhyRes] + p[PhyResExtra]
	}
	return result
}

func (p *Props) Add(b Props) {
	for i := range p {
		p[i] += b[i]
	}
}

func PrettyPrintStatsSlice(stats []float64) []string {
	r := make([]string, 0)
	var sb strings.Builder
	for i, v := range stats {
		if v == 0 {
			continue
		}
		sb.WriteString(PropTypeString[i])
		sb.WriteString(": ")
		sb.WriteString(strconv.FormatFloat(v, 'f', 2, 32))
		r = append(r, sb.String())
		sb.Reset()
	}

	return r
}

func PrettyPrintStatsMap(stats PropMap) []string {
	r := make([]string, 0, len(stats))
	var sb strings.Builder
	for k, v := range stats {
		if v == 0 {
			continue
		}
		sb.WriteString(k.String())
		sb.WriteString(": ")
		sb.WriteString(strconv.FormatFloat(v, 'f', 2, 32))
		r = append(r, sb.String())
		sb.Reset()
	}
	return r
}

var PropTypeString = [...]string{
	"n/a",
	"hp_base",
	"hp%",
	"hp",
	"hp_extra",
	"def_base",
	"def%",
	"def",
	"def_extra",
	"atk_base",
	"atk%",
	"atk",
	"atk_extra",
	"cr",
	"cr_extra",
	"cd",
	"cd_extra",
	"er",
	"er_extra",
	"def_ignore",
	"em",
	"em_extra",
	"heal",
	"heal_taken",
	"cd_reduction",
	"dmg%",
	"pyro%",
	"pyro%_extra",
	"hydro%",
	"hydro%_extra",
	"cryo%",
	"cryo%_extra",
	"electro%",
	"electro%_extra",
	"anemo%",
	"anemo%_extra",
	"geo%",
	"geo%_extra",
	"dendro%",
	"dendro%_extra",
	"phy%",
	"phy%_extra",
	"dmg_res",
	"pyro_res",
	"pyro_res_extra",
	"hydro_res",
	"hydro_res_extra",
	"cryo_res",
	"cryo_res_extra",
	"electro_res",
	"electro_res_extra",
	"anemo_res",
	"anemo_res_extra",
	"geo_res",
	"geo_res_extra",
	"dendro_res",
	"dendro_res_extra",
	"phy_res",
	"phy_res_extra",
	"overload_bonus",
	"superconduct_bonus",
	"melt_bonus",
	"vaporize_bonus",
	"swirl_electro_bonus",
	"swirl_hydro_bonus",
	"swirl_pyro_bonus",
	"swirl_cryo_bonus",
	"electro_charged_bonus",
	"shatter_bonus",
	"burning_bonus",
	"aggravate_bonus",
	"spread_bonus",
	"bloom_bonus",
	"bountiful_core_bonus",
	"burgeon_bonus",
	"hyperbloom_bonus",
	"overall_spd",
	"overall_spd_mult",
	"atk_spd",
	"move_spd",
	"cost_stam",
	"shield_strength",
	"",
}

func StrToPropType(s string) Prop {
	for i, v := range PropTypeString {
		if v == s {
			return Prop(i)
		}
	}
	return -1
}

var ElementToPropBonus = map[ElementType]Prop{
	ElementFire:     PyroP,
	ElementWater:    HydroP,
	ElementIce:      CryoP,
	ElementElectric: ElectroP,
	ElementWind:     AnemoP,
	ElementRock:     GeoP,
	ElementGrass:    DendroP,
	ElementNone:     PhyP, // TODO
}
