package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func init() {

}

func (c *Core) pyroResonance() {
	// val := make([]float64, attributes.EndStatType)
	// val[attributes.ATKP] = 0.25
	// f := func() ([]float64, bool) {
	// 	return val, true
	// }
	// for _, c := range chars {
	// 	c.AddStatMod(character.StatMod{
	// 		Base:         modifier.NewBase("pyro-res", -1),
	// 		AffectedStat: attributes.NoStat,
	// 		Amount:       f,
	// 	})
	// }
}

func (c *Core) hydroResonance() {
	// TODO: reduce pyro duration not implemented; may affect bennett Q?
	// val := make([]float64, attributes.EndStatType)
	// val[attributes.HPP] = 0.25
	// for _, c := range chars {
	// 	c.AddStatMod(character.StatMod{
	// 		Base:         modifier.NewBase("hydro-res-hpp", -1),
	// 		AffectedStat: attributes.HPP,
	// 		Amount: func() ([]float64, bool) {
	// 			return val, true
	// 		},
	// 	})
	// }
}

func (c *Core) cryoResonance() {
	// val := make([]float64, attributes.EndStatType)
	// val[attributes.CR] = .15
	// f := func(ae *info.AttackEvent, t info.Target) ([]float64, bool) {
	// 	r, ok := t.(*enemy.Enemy)
	// 	if !ok {
	// 		return nil, false
	// 	}
	// 	if r.AuraContains(attributes.Cryo) || r.AuraContains(attributes.Frozen) {
	// 		return val, true
	// 	}
	// 	return nil, false
	// }
	// for _, c := range chars {
	// 	c.AddAttackMod(character.AttackMod{
	// 		Base:   modifier.NewBase("cryo-res", -1),
	// 		Amount: f,
	// 	})
	// }
}

func (c *Core) electroResonance() {
	// last := 0
	// //nolint:unparam // ignoring for now, event refactor should get rid of bool return of event sub
	// recoverParticle := func(_ ...any) bool {
	// 	if s.F-last < 300 && last != 0 { // every 5 seconds
	// 		return false
	// 	}
	// 	s.Player.DistributeParticle(character.Particle{
	// 		Source: "electro-res",
	// 		Num:    1,
	// 		Ele:    attributes.Electro,
	// 	})
	// 	last = s.F
	// 	return false
	// }

	// recoverNoGadget := func(args ...any) bool {
	// 	if _, ok := args[0].(*enemy.Enemy); !ok {
	// 		return false
	// 	}
	// 	return recoverParticle(args...)
	// }
	// s.Events.Subscribe(event.OnOverload, recoverNoGadget, "electro-res")
	// s.Events.Subscribe(event.OnSuperconduct, recoverNoGadget, "electro-res")
	// s.Events.Subscribe(event.OnElectroCharged, recoverNoGadget, "electro-res")
	// s.Events.Subscribe(event.OnQuicken, recoverNoGadget, "electro-res")
	// s.Events.Subscribe(event.OnAggravate, recoverNoGadget, "electro-res")
	// s.Events.Subscribe(event.OnHyperbloom, recoverParticle, "electro-res")
}

func (c *Core) geoResonance() {
	// // Increases shield strength by 15%. Additionally, characters protected by a shield will have the
	// // following special characteristics:

	// //	DMG dealt increased by 15%, dealing DMG to enemies will decrease their Geo RES by 20% for 15s.
	// f := func() (float64, bool) { return 0.15, true }
	// s.Player.Shields.AddShieldBonusMod("geo-res", -1, f)

	// // shred geo res of target
	// s.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
	// 	t, ok := args[0].(*enemy.Enemy)
	// 	if !ok {
	// 		return false
	// 	}
	// 	atk := args[1].(*info.AttackEvent)
	// 	if s.Player.Shields.CharacterIsShielded(atk.Info.ActorIndex, s.Player.Active()) {
	// 		t.AddResistMod(info.ResistMod{
	// 			Base:  modifier.NewBaseWithHitlag("geo-res", 15*60),
	// 			Ele:   attributes.Geo,
	// 			Value: -0.2,
	// 		})
	// 	}
	// 	return false
	// }, "geo res")

	// val := make([]float64, attributes.EndStatType)
	// val[attributes.DmgP] = .15
	// atkf := func(ae *info.AttackEvent, t info.Target) ([]float64, bool) {
	// 	if s.Player.Shields.CharacterIsShielded(ae.Info.ActorIndex, s.Player.Active()) {
	// 		return val, true
	// 	}
	// 	return nil, false
	// }
	// for _, c := range chars {
	// 	c.AddAttackMod(character.AttackMod{
	// 		Base:   modifier.NewBase("geo-res", -1),
	// 		Amount: atkf,
	// 	})
	// }
}

func (c *Core) anemoResonance() {
	// s.Player.AddStamPercentMod("anemo-res-stam", -1, func(a action.Action) (float64, bool) {
	// 	return -0.15, false
	// })
	// // TODO: movement spd increase?
	// for _, c := range chars {
	// 	c.AddCooldownMod(character.CooldownMod{
	// 		Base:   modifier.NewBase("anemo-res-cd", -1),
	// 		Amount: func(a action.Action) float64 { return -0.05 },
	// 	})
	// }
}

func (c *Core) dendroResonance() {
	// val := make([]float64, attributes.EndStatType)
	// val[attributes.EM] = 50
	// for _, c := range chars {
	// 	c.AddStatMod(character.StatMod{
	// 		Base:         modifier.NewBase("dendro-res-50", -1),
	// 		AffectedStat: attributes.EM,
	// 		Amount: func() ([]float64, bool) {
	// 			return val, true
	// 		},
	// 	})
	// }

	// twoBuff := make([]float64, attributes.EndStatType)
	// twoBuff[attributes.EM] = 30
	// twoEl := func(args ...any) bool {
	// 	if _, ok := args[0].(*enemy.Enemy); !ok {
	// 		return false
	// 	}
	// 	for _, c := range chars {
	// 		c.AddStatMod(character.StatMod{
	// 			Base:         modifier.NewBaseWithHitlag("dendro-res-30", 6*60),
	// 			AffectedStat: attributes.EM,
	// 			Amount: func() ([]float64, bool) {
	// 				return twoBuff, true
	// 			},
	// 		})
	// 	}
	// 	return false
	// }
	// s.Events.Subscribe(event.OnBurning, twoEl, "dendro-res")
	// s.Events.Subscribe(event.OnBloom, twoEl, "dendro-res")
	// s.Events.Subscribe(event.OnQuicken, twoEl, "dendro-res")

	// threeBuff := make([]float64, attributes.EndStatType)
	// threeBuff[attributes.EM] = 20
	// threeEl := func(_ ...any) bool {
	// 	for _, c := range chars {
	// 		c.AddStatMod(character.StatMod{
	// 			Base:         modifier.NewBaseWithHitlag("dendro-res-20", 6*60),
	// 			AffectedStat: attributes.EM,
	// 			Amount: func() ([]float64, bool) {
	// 				return threeBuff, true
	// 			},
	// 		})
	// 	}
	// 	return false
	// }
	// threeElNoGadget := func(args ...any) bool {
	// 	if _, ok := args[0].(*enemy.Enemy); !ok {
	// 		return false
	// 	}
	// 	return threeEl(nil)
	// }
	// s.Events.Subscribe(event.OnAggravate, threeElNoGadget, "dendro-res")
	// s.Events.Subscribe(event.OnSpread, threeElNoGadget, "dendro-res")
	// s.Events.Subscribe(event.OnHyperbloom, threeEl, "dendro-res")
	// s.Events.Subscribe(event.OnBurgeon, threeEl, "dendro-res")
}

func (c *Core) setupResonance() {
	chars := c.player.Chars()
	if len(chars) < 4 {
		return // no resonance if less than 4 chars
	}
	// count number of ele first
	count := make(map[attributes.ElementType]int)
	for _, char := range chars {
		count[char.GetBase().Element]++
	}

	for k, v := range count {
		if v < 2 {
			continue
		}
		switch k {
		case attributes.ElementFire:
			c.pyroResonance()
		case attributes.ElementWater:
			c.hydroResonance()
		case attributes.ElementIce:
			c.cryoResonance()
		case attributes.ElementElectric:
			c.electroResonance()
		case attributes.ElementRock:
			c.geoResonance()
		case attributes.ElementWind:
			c.anemoResonance()
		case attributes.ElementGrass:
			c.dendroResonance()
		}
	}
}
