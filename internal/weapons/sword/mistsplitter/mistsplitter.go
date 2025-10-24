package mistsplitter

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/modifier"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	buffKey = "mistsplitter"
)

func init() {
	core.RegisterWeaponFunc(keys.MistsplitterReforged, NewWeapon)

	modifier.Register(buffKey, modifier.Config{
		ElementDurability: 100,
	})
}

// Gain a 12% Elemental DMG Bonus for all elements and receive the might of the
// Mistsplitter's Emblem. At stack levels 1/2/3, the Mistsplitter's Emblem
// provides a 8/16/28% Elemental DMG Bonus for the character's Elemental Type.
// The character will obtain 1 stack of Mistsplitter's Emblem in each of the
// following scenarios: Normal Attack deals Elemental DMG (stack lasts 5s),
// casting Elemental Burst (stack lasts 10s); Energy is less than 100% (stack
// disappears when Energy is full). Each stack's duration is calculated
// independently.
func NewWeapon(c core.Core, char core.Character, p info.WeaponProfile) (core.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	// perm buff
	base := 0.09 + float64(r)*0.03
	char.AddModifier(info.Modifier{
		Key:    buffKey,
		Source: char.Key(),
		Props: attributes.PropMap{
			attributes.PyroP:    base,
			attributes.HydroP:   base,
			attributes.CryoP:    base,
			attributes.ElectroP: base,
			attributes.AnemoP:   base,
			attributes.GeoP:     base,
			attributes.DendroP:  base,
		},
	})

	// stacking buff
	stack := 0.06 + float64(r)*0.02
	maxBonus := 0.03 + float64(r)*0.01
	ele := char.GetBase().Element

	char.AddModifier(info.Modifier{
		Key:    stacksKey,
		Source: char.Key(),
		State: buffState{
			stack:    stack,
			maxBonus: maxBonus,
			prop:     attributes.ElementToPropBonus[ele],
		},
	})

	return w, nil
}
