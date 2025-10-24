package mistsplitter

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/modifier"
)

const (
	stacksKey     = "mistsplitter-stacks"
	normalBuffKey = "mistsplitter-normal"
	burstBuffKey  = "mistsplitter-burst"
	energyBuffKey = "mistsplitter-energy"
)

type buffState struct {
	stack    float64
	maxBonus float64
	prop     attributes.Prop
}

func init() {
	modifier.Register(stacksKey, modifier.Config{
		ElementDurability: 100,
		Listeners: modifier.Listeners{
			OnStat: func(mod *modifier.Instance) {
				stacksUpdateEnergyStack(mod)
				stacksUpdateBuff(mod)
			},
			OnCharacterAction: stacksUpdateBurstStack,
			OnEnemyHit:        stacksUpdateNormalStack,
		},
	})

	modifier.Register(energyBuffKey, modifier.Config{
		ElementDurability: 100,
		Listeners: modifier.Listeners{
			OnAdd:    onAddStackCB(1),
			OnRemove: onAddStackCB(-1),
		},
	})
	modifier.Register(normalBuffKey, modifier.Config{
		ElementDurability: 100,
		Duration:          5 * 60,
		AffectedByHitlag:  true,
		Listeners: modifier.Listeners{
			OnAdd:    onAddStackCB(1),
			OnRemove: onAddStackCB(-1),
		},
	})
	modifier.Register(burstBuffKey, modifier.Config{
		ElementDurability: 100,
		Duration:          10 * 60,
		AffectedByHitlag:  true,
		Listeners: modifier.Listeners{
			OnAdd:    onAddStackCB(1),
			OnRemove: onAddStackCB(-1),
		},
	})
}

func stacksUpdateEnergyStack(mod *modifier.Instance) {
	char := mod.Core().GetCharacterByTarget(mod.Owner())
	if char.Energy() < char.EnergyMax() || char.EnergyMax() == 0 {
		char.AddModifier(info.Modifier{Key: energyBuffKey})
	} else {
		char.RemoveModifier(energyBuffKey)
	}
}

func stacksUpdateBuff(mod *modifier.Instance) {
	char := mod.Core().GetCharacterByTarget(mod.Owner())
	state := mod.State().(buffState)

	stacks := char.Tag(stacksKey)
	bonus := state.stack * float64(stacks)
	if stacks >= 3 {
		bonus += state.maxBonus
	}
	mod.SetProperty(state.prop, bonus)
}

func stacksUpdateBurstStack(mod *modifier.Instance, char int, action info.Action, p map[string]int) {
	c := mod.Core().GetCharacterByTarget(mod.Owner())
	if c.GetIndex() != char {
		return
	}
	if action != info.ActionBurst {
		return
	}
	c.AddModifier(info.Modifier{Key: burstBuffKey})
}

func stacksUpdateNormalStack(mod *modifier.Instance, target keys.Target, attack *info.AttackEvent) {
	c := mod.Core().GetCharacterByTarget(mod.Owner())
	if attack.Info.ActorIndex != c.GetIndex() {
		return
	}
	if attack.Info.AttackTag != info.AttackTagNormal && attack.Info.AttackTag != info.AttackTagExtra {
		return
	}
	if attack.Info.Element == attributes.ElementNone {
		return
	}
	c.AddModifier(info.Modifier{Key: normalBuffKey})
}

func onAddStackCB(n int) func(mod *modifier.Instance) {
	return func(mod *modifier.Instance) {
		char := mod.Core().GetCharacterByTarget(mod.Owner())
		stacks := min(max(char.Tag(stacksKey)+n, 0), 3)
		char.SetTag(stacksKey, stacks)
	}
}
