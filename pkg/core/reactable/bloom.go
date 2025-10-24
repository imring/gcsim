package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/target/gadget"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

const DendroCoreDelay = 30

func (r *Reactable) TryBloom(a *info.AttackEvent) bool {
	// can be hydro bloom, dendro bloom, or quicken bloom
	if a.Info.Durability < ZeroDur {
		return false
	}
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.ElementWater:
		// this part is annoying. bloom will happen if any of the dendro like aura is present
		// so we gotta check for all 3...
		switch {
		case r.Durability[attributes.ElementGrass] > ZeroDur:
		case r.Durability[attributes.ElementOverdose] > ZeroDur:
		case r.Durability[attributes.ElementBurningFuel] > ZeroDur:
		default:
			return false
		}
		// reduce only check for one element so have to call twice to check for quicken as well
		consumed = r.reduce(attributes.ElementGrass, a.Info.Durability, 0.5)
		f := r.reduce(attributes.ElementOverdose, a.Info.Durability, 0.5)
		if f > consumed {
			consumed = f
		}
	case attributes.ElementGrass:
		if r.Durability[attributes.ElementWater] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.ElementWater, a.Info.Durability, 2)
	default:
		return false
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	r.addBloomGadget(a)
	r.core.Events().Bloom.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})
	return true
}

// this function should only be called after a catalyze reaction (queued to the end of current frame)
// this reaction will check if any hydro exists and if so trigger a bloom reaction
func (r *Reactable) tryQuickenBloom(a *info.AttackEvent) {
	if r.Durability[attributes.ElementOverdose] < ZeroDur {
		// this should be a sanity check; should not happen realistically unless something wipes off
		// the quicken immediately (same frame) after catalyze
		return
	}
	if r.Durability[attributes.ElementWater] < ZeroDur {
		return
	}
	avail := r.Durability[attributes.ElementOverdose]
	consumed := r.reduce(attributes.ElementWater, avail, 2)
	r.Durability[attributes.ElementOverdose] -= consumed

	r.addBloomGadget(a)
	r.core.Events().Bloom.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})
}

type DendroCore struct {
	*gadget.Gadget
	core      core.Core
	srcFrame  int
	CharIndex int
}

func (r *Reactable) addBloomGadget(a *info.AttackEvent) {
	r.core.Tasks().Add(func() {
		t := NewDendroCore(r.core, r.self.Shape(), a)
		r.core.AddGadget(t)
		r.core.Events().DendroCoreSpawned.Emit(event.DendroCoreSpawnedEvent{
			Target:      r.self.Key(),
			AttackEvent: a,
		})
		r.core.Log().NewEvent(
			"dendro core spawned",
			glog.LogElementEvent,
			a.Info.ActorIndex,
		).
			Write("src", t.Src()).
			Write("expiry", r.core.F()+t.Duration())
	}, DendroCoreDelay)
}

func NewDendroCore(c core.Core, shp geometry.Shape, a *info.AttackEvent) *DendroCore {
	s := &DendroCore{
		core:      c,
		srcFrame:  c.F(),
		CharIndex: a.Info.ActorIndex,
	}

	circ, ok := shp.(*geometry.Circle)
	if !ok {
		panic("rectangle target hurtbox is not supported for dendro core spawning")
	}

	// for simplicity, seeds spawn randomly at radius + 0.5
	r := circ.Radius() + 0.5
	gadget := gadget.New(c, geometry.CalcRandomPointFromCenter(circ.Pos(), r, r, c.Rand()), 2, info.GadgetTypDendroCore)
	gadget.SetDuration(300) // ???

	char := c.GetCharacter(a.Info.ActorIndex)

	explode := func(reason string) func() {
		return func() {
			c.Tasks().Add(func() {
				ai, snap := NewBloomAttack(char, s, nil)
				ap := combat.NewCircleHitOnTarget(s, nil, 5)
				c.QueueAttackWithSnap(snap, info.QueueAttack{
					Info:    ai,
					Pattern: ap,
				})

				// self damage
				ai.Abil += info.SelfDamageSuffix
				ai.FlatDmg = 0.05 * ai.FlatDmg
				ap.SkipTargets[info.TargettablePlayer] = false
				ap.SkipTargets[info.TargettableEnemy] = true
				ap.SkipTargets[info.TargettableGadget] = true
				c.QueueAttackWithSnap(snap, info.QueueAttack{
					Info:    ai,
					Pattern: ap,
				})

				c.Log().NewEvent(
					"dendro core "+reason,
					glog.LogElementEvent,
					char.GetIndex(),
				).Write("src", s.Src())
			}, 1)
		}
	}
	gadget.OnExpiry = explode("expired")
	gadget.OnKill = explode("killed")
	s.Gadget = gadget

	return s
}

func (s *DendroCore) Tick() {
	// this is needed since gadget tick
	s.Gadget.Tick()
}

func (s *DendroCore) HandleAttack(atk *info.AttackEvent) float64 {
	s.core.Events().GadgetHit.Emit(event.TargetHitEvent{
		Target:      s.Key(),
		AttackEvent: atk,
	})
	s.Attack(atk, nil)
	return 0
}

func (s *DendroCore) Attack(atk *info.AttackEvent, evt glog.Event) (float64, bool) {
	if atk.Info.Durability < ZeroDur {
		return 0, false
	}

	char := s.core.GetCharacter(atk.Info.ActorIndex)
	// only contact with pyro/electro to trigger burgeon/hyperbloom accordingly
	switch atk.Info.Element {
	case attributes.ElementElectric:
		// trigger hyperbloom targets the nearest enemy
		// it can also do damage to player in small aoe
		s.core.Tasks().Add(func() {
			ai, snap := NewHyperbloomAttack(char, s)
			// queue dmg nearest enemy within radius 15
			enemy := s.core.ClosestEnemyWithinArea(combat.NewCircleHitOnTarget(s.Gadget, nil, 15), nil)
			if enemy != nil {
				ap := combat.NewCircleHitOnTarget(enemy, nil, 1)
				s.core.QueueAttackWithSnap(snap, info.QueueAttack{
					Info:    ai,
					Pattern: ap,
				})

				// also queue self damage
				ai.Abil += info.SelfDamageSuffix
				ai.FlatDmg = 0.05 * ai.FlatDmg
				ap.SkipTargets[info.TargettablePlayer] = false
				ap.SkipTargets[info.TargettableEnemy] = true
				ap.SkipTargets[info.TargettableGadget] = true
				s.core.QueueAttackWithSnap(snap, info.QueueAttack{
					Info:    ai,
					Pattern: ap,
				})
			}
		}, 60)

		s.Gadget.SetOnKill(nil)
		s.Gadget.Kill(nil)
		s.core.Events().Hyperbloom.Emit(event.ReactionEvent{
			Target:      s.Key(),
			AttackEvent: atk,
		})
		s.core.Log().NewEvent(
			"hyperbloom triggered",
			glog.LogElementEvent,
			char.GetIndex(),
		).
			Write("dendro_core_char", s.CharIndex).
			Write("dendro_core_src", s.Gadget.Src())
	case attributes.ElementFire:
		// trigger burgeon, aoe dendro damage
		// self damage
		s.core.Tasks().Add(func() {
			ai, snap := NewBurgeonAttack(char, s)
			ap := combat.NewCircleHitOnTarget(s, nil, 5)
			s.core.QueueAttackWithSnap(snap, info.QueueAttack{
				Info:    ai,
				Pattern: ap,
			})

			// queue self damage
			ai.Abil += info.SelfDamageSuffix
			ai.FlatDmg = 0.05 * ai.FlatDmg
			ap.SkipTargets[info.TargettablePlayer] = false
			ap.SkipTargets[info.TargettableEnemy] = true
			ap.SkipTargets[info.TargettableGadget] = true
			s.core.QueueAttackWithSnap(snap, info.QueueAttack{
				Info:    ai,
				Pattern: ap,
			})
		}, 1)

		s.Gadget.SetOnKill(nil)
		s.Gadget.Kill(nil)
		s.core.Events().Burgeon.Emit(event.ReactionEvent{
			Target:      s.Key(),
			AttackEvent: atk,
		})
		s.core.Log().NewEvent(
			"burgeon triggered",
			glog.LogElementEvent,
			char.GetIndex(),
		).
			Write("dendro_core_char", s.CharIndex).
			Write("dendro_core_src", s.Gadget.Src())
	default:
		return 0, false
	}

	return 0, false
}

const (
	BloomMultiplier      = 2
	BurgeonMultiplier    = 3
	HyperbloomMultiplier = 3
)

func NewBloomAttack(char core.Character, src core.Target, modify func(*info.Attack)) (info.Attack, info.Snapshot) {
	ai := info.Attack{
		ActorIndex:       char.GetIndex(),
		DamageSrc:        src.Key(),
		Element:          attributes.ElementGrass,
		AttackTag:        info.AttackTagBloom,
		ICDTag:           info.ICDTagBloomDamage,
		ICDGroup:         info.ICDGroupReactionA,
		StrikeType:       info.StrikeTypeDefault,
		Abil:             string(info.ReactionTypeBloom),
		IgnoreDefPercent: 1,
	}
	if modify != nil {
		modify(&ai)
	}
	flatdmg, snap := calcReactionDmg(&ai, char)
	ai.FlatDmg = BloomMultiplier * flatdmg
	return ai, snap
}

func NewBurgeonAttack(char core.Character, src core.Target) (info.Attack, info.Snapshot) {
	ai := info.Attack{
		ActorIndex:       char.GetIndex(),
		DamageSrc:        src.Key(),
		Element:          attributes.ElementGrass,
		AttackTag:        info.AttackTagBurgeon,
		ICDTag:           info.ICDTagBurgeonDamage,
		ICDGroup:         info.ICDGroupReactionA,
		StrikeType:       info.StrikeTypeDefault,
		Abil:             string(info.ReactionTypeBurgeon),
		IgnoreDefPercent: 1,
	}
	flatdmg, snap := calcReactionDmg(&ai, char)
	ai.FlatDmg = BurgeonMultiplier * flatdmg
	return ai, snap
}

func NewHyperbloomAttack(char core.Character, src core.Target) (info.Attack, info.Snapshot) {
	ai := info.Attack{
		ActorIndex:       char.GetIndex(),
		DamageSrc:        src.Key(),
		Element:          attributes.ElementGrass,
		AttackTag:        info.AttackTagHyperbloom,
		ICDTag:           info.ICDTagHyperbloomDamage,
		ICDGroup:         info.ICDGroupReactionA,
		StrikeType:       info.StrikeTypeDefault,
		Abil:             string(info.ReactionTypeHyperbloom),
		IgnoreDefPercent: 1,
	}
	flatdmg, snap := calcReactionDmg(&ai, char)
	ai.FlatDmg = HyperbloomMultiplier * flatdmg
	return ai, snap
}

func (s *DendroCore) SetDirection(trg geometry.Point) {}
func (s *DendroCore) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}
