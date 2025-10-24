package attributes

type (
	ElementType int
	ElementMap  map[ElementType]float64
)

const (
	ElementNone           ElementType = 0
	ElementFire           ElementType = 1
	ElementWater          ElementType = 2
	ElementGrass          ElementType = 3
	ElementElectric       ElementType = 4
	ElementIce            ElementType = 5
	ElementFrozen         ElementType = 6
	ElementWind           ElementType = 7
	ElementRock           ElementType = 8
	ElementAntiFire       ElementType = 9
	ElementVehicleMuteIce ElementType = 10
	ElementMushroom       ElementType = 11
	ElementOverdose       ElementType = 12
	ElementWood           ElementType = 13
	LiquidPhlogiston      ElementType = 14
	SolidPhlogiston       ElementType = 15
	SolidifyPhlogiston    ElementType = 16
	ElementBurning        ElementType = 17 // TODO: should be removed
	ElementBurningFuel    ElementType = 18 // TODO: should be removed
	ElementCOUNT          ElementType = 19
)

var elementoToString = map[ElementType]string{
	ElementElectric:    "electro",
	ElementFire:        "pyro",
	ElementWater:       "hydro",
	ElementGrass:       "dendro",
	ElementWind:        "anemo",
	ElementIce:         "cryo",
	ElementRock:        "geo",
	ElementFrozen:      "frozen",
	ElementOverdose:    "quicken",
	ElementBurning:     "burning",
	ElementBurningFuel: "dendro-fuel",
}

func (e ElementType) String() string {
	if str, ok := elementoToString[e]; ok {
		return str
	}
	return "none"
}

func StringToEle(s string) ElementType {
	for i, v := range elementoToString {
		if v == s {
			return ElementType(i)
		}
	}
	return ElementNone
}
