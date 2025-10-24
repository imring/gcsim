package crimson

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/modifier"
)

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

const (
	cw2pc      = "crimson-2pc"
	cw4pc      = "crimson-4pc"
	cw4pcStack = "crimson-4pc-stack"
)

func init() {
	core.RegisterSetFunc(keys.CrimsonWitchOfFlames, NewSet)

	modifier.Register(cw2pc, modifier.Config{
		ElementDurability: 100,
	})

	modifier.Register(cw4pc, modifier.Config{
		ElementDurability: 100,
	})

	modifier.Register(cw4pcStack, modifier.Config{
		Stacking:          modifier.MultipleAllRefresh,
		MaxCount:          3,
		Duration:          10 * 60,
		ElementDurability: 100,
		AffectedByHitlag:  true,
	})
}

func NewSet(c core.Core, char core.Character, count int, param map[string]int) (core.Set, error) {
	s := Set{}

	if count >= 2 {
		s.pc2(char)
	}
	if count >= 4 {
		s.pc4(c, char)
	}

	return &s, nil
}

func (s *Set) pc2(char core.Character) {
	char.AddModifier(info.Modifier{
		Key: cw2pc,
		Props: attributes.PropMap{
			attributes.PyroP: 0.15,
		},
	})
}

func (s *Set) pc4(c core.Core, char core.Character) {
	char.AddModifier(info.Modifier{
		Key: cw4pc,
		Props: attributes.PropMap{
			attributes.OverloadBonus: 0.4,
			attributes.BurningBonus:  0.4,
			attributes.BurgeonBonus:  0.4,
			attributes.VaporizeBonus: 0.15,
			attributes.MeltBonus:     0.15,
		},
	})

	c.Events().Skill.Subscribe(func(event event.ActionEvent) {
		if event.CharIndex != char.GetIndex() {
			return
		}

		char.AddModifier(info.Modifier{
			Key: cw4pcStack,
			Props: attributes.PropMap{
				attributes.PyroP: 0.15 * 0.5,
			},
		})
	})
}
