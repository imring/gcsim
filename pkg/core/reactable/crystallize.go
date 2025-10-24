package reactable

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/target/gadget"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

type CrystallizeShield struct {
	*shield.Tmpl
	emBonus float64
}

func (r *Reactable) TryCrystallizeElectro(a *info.AttackEvent) bool {
	if r.Durability[attributes.ElementElectric] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.ElementElectric, info.ReactionTypeCrystallizeElectro, r.core.Events().CrystallizeElectro)
	}
	return false
}

func (r *Reactable) TryCrystallizeHydro(a *info.AttackEvent) bool {
	if r.Durability[attributes.ElementWater] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.ElementWater, info.ReactionTypeCrystallizeHydro, r.core.Events().CrystallizeHydro)
	}
	return false
}

func (r *Reactable) TryCrystallizeCryo(a *info.AttackEvent) bool {
	if r.Durability[attributes.ElementIce] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.ElementIce, info.ReactionTypeCrystallizeCryo, r.core.Events().CrystallizeCryo)
	}
	return false
}

func (r *Reactable) TryCrystallizePyro(a *info.AttackEvent) bool {
	if r.Durability[attributes.ElementFire] > ZeroDur || r.Durability[attributes.ElementBurning] > ZeroDur {
		reacted := r.tryCrystallizeWithEle(a, attributes.ElementFire, info.ReactionTypeCrystallizePyro, r.core.Events().CrystallizePyro)
		r.burningCheck()
		return reacted
	}
	return false
}

func (r *Reactable) TryCrystallizeFrozen(a *info.AttackEvent) bool {
	if r.Durability[attributes.ElementFrozen] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.ElementFrozen, info.ReactionTypeCrystallizeCryo, r.core.Events().CrystallizeCryo)
	}
	return false
}

func (r *Reactable) tryCrystallizeWithEle(a *info.AttackEvent, ele attributes.ElementType, rt info.ReactionType, evt event.ReactionEventHandler) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.crystallizeGCD != -1 && r.core.F() < r.crystallizeGCD {
		return false
	}
	r.crystallizeGCD = r.core.F() + 60
	char := r.core.GetCharacter(a.Info.ActorIndex)
	r.addCrystallizeShard(char, rt, ele, r.core.F())
	// reduce
	r.reduce(ele, a.Info.Durability, 0.5)
	//TODO: confirm u can only crystallize once
	a.Info.Durability = 0
	a.Reacted = true
	// event
	evt.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})
	// check freeze + ec
	switch {
	case ele == attributes.ElementElectric && r.Durability[attributes.ElementWater] > ZeroDur:
		r.checkEC()
	case ele == attributes.ElementFrozen:
		r.checkFreeze()
	}
	return true
}

func NewCrystallizeShield(index int, typ attributes.ElementType, src, lvl int, em float64, expiry int) *CrystallizeShield {
	s := &CrystallizeShield{}
	tmpl := &shield.Tmpl{}

	// TODO: 100 lvl
	lvl--
	if lvl > 89 {
		lvl = 89
	}
	if lvl < 0 {
		lvl = 0
	}

	tmpl.ActorIndex = index
	tmpl.Target = -1
	tmpl.Ele = typ
	tmpl.ShieldType = info.ShieldCrystallize
	tmpl.Name = "Crystallize " + typ.String()
	tmpl.Src = src
	tmpl.HP = shieldBaseHP[lvl]
	tmpl.Expires = expiry
	s.Tmpl = tmpl

	s.emBonus = (40.0 / 9.0) * (em / (1400 + em))

	return s
}

func (c *CrystallizeShield) OnDamage(dmg float64, ele attributes.ElementType, bonus float64) (float64, bool) {
	bonus += c.emBonus
	return c.OnDamage(dmg, ele, bonus)
}

var shieldBaseHP = []float64{
	91.1791000366211,
	98.7076644897461,
	106.236221313477,
	113.764770507813,
	121.293319702148,
	128.821884155273,
	136.35041809082,
	143.878982543945,
	151.407516479492,
	158.936080932617,
	169.991485595703,
	181.076248168945,
	192.190368652344,
	204.048202514648,
	215.938995361328,
	227.862747192383,
	247.685943603516,
	267.542114257813,
	287.431213378906,
	303.826416015625,
	320.225219726563,
	336.627624511719,
	352.319274902344,
	368.010925292969,
	383.702545166016,
	394.432373046875,
	405.181457519531,
	415.949920654297,
	426.737640380859,
	437.544708251953,
	450.600006103516,
	463.700286865234,
	476.845581054688,
	491.127502441406,
	502.554565429688,
	514.012084960938,
	531.409606933594,
	549.979614257813,
	568.584899902344,
	584.996520996094,
	605.670349121094,
	626.38623046875,
	646.052307128906,
	665.755615234375,
	685.49609375,
	700.839416503906,
	723.333129882813,
	745.865295410156,
	768.435729980469,
	786.791931152344,
	809.538818359375,
	832.329040527344,
	855.162658691406,
	878.039611816406,
	899.484802246094,
	919.361999511719,
	946.039611816406,
	974.764221191406,
	1003.57861328125,
	1030.07702636719,
	1056.63500976563,
	1085.24633789063,
	1113.92443847656,
	1149.25866699219,
	1178.06481933594,
	1200.22375488281,
	1227.66027832031,
	1257.24304199219,
	1284.91735839844,
	1314.7529296875,
	1342.66516113281,
	1372.75244140625,
	1396.32104492188,
	1427.31237792969,
	1458.37451171875,
	1482.33581542969,
	1511.91088867188,
	1541.54931640625,
	1569.15368652344,
	1596.81433105469,
	1622.41967773438,
	1648.07397460938,
	1666.37609863281,
	1684.67822265625,
	1702.98034667969,
	1726.10473632813,
	1754.67150878906,
	1785.86657714844,
	1817.13745117188,
	1851.06030273438,
}

type CrystallizeShard struct {
	*gadget.Gadget
	core core.Core
	// earliest that a shard can be picked up after spawn
	EarliestPickup int
	// captures the shield because em snapshots
	Shield *CrystallizeShield
	// for logging purposes
	src    int
	expiry int
}

func (r *Reactable) addCrystallizeShard(char core.Character, rt info.ReactionType, typ attributes.ElementType, src int) {
	// delay shard spawn
	r.core.Tasks().Add(func() {
		// grab current snapshot for shield
		ai := info.Attack{
			ActorIndex: char.GetIndex(),
			DamageSrc:  r.self.Key(),
			Abil:       string(rt),
		}
		snap := char.Snapshot(&ai)
		lvl := snap.Level
		// shield snapshots em on shard spawn
		em := snap.Stats.Props.EM()
		// expiry will get set properly later
		shd := NewCrystallizeShield(char.GetIndex(), typ, src, lvl, em, -1)
		cs := NewCrystallizeShard(r.core, r.self.Shape(), shd)
		r.core.AddGadget(cs)
		r.core.Log().NewEvent(
			fmt.Sprintf("%v crystallize shard spawned", cs.Shield.Element()),
			glog.LogElementEvent,
			cs.Shield.ShieldOwner(),
		).
			Write("src", cs.src).
			Write("expiry", cs.expiry).
			Write("earliest_pickup", cs.EarliestPickup)
	}, 23)
}

func NewCrystallizeShard(c core.Core, shp geometry.Shape, shd *CrystallizeShield) *CrystallizeShard {
	cs := &CrystallizeShard{
		core: c,
	}

	circ, ok := shp.(*geometry.Circle)
	if !ok {
		panic("rectangle target hurtbox is not supported for crystallize shard spawning")
	}

	// for simplicity, crystallize shards spawn randomly at radius + 0.5
	r := circ.Radius() + 0.5
	// radius 2 is ok
	gadget := gadget.New(c, geometry.CalcRandomPointFromCenter(circ.Pos(), r, r, c.Rand()), 2, info.GadgetTypCrystallizeShard)

	// shard lasts for 15s from shard spawn
	gadget.SetDuration(15 * 60)
	// earliest shard pickup is 54f from crystallize text, so 31f from shard spawn
	cs.EarliestPickup = c.F() + 31
	cs.Shield = shd
	cs.src = c.F()
	cs.expiry = c.F() + gadget.Duration()

	cs.Gadget = gadget

	return cs
}

func (cs *CrystallizeShard) AddShieldKillShard() bool {
	// don't pick up if shard is not available for pick up yet
	if cs.core.F() < cs.EarliestPickup {
		cs.core.Log().NewEvent(
			fmt.Sprintf("%v crystallize shard could not be picked up", cs.Shield.Element()),
			glog.LogElementEvent,
			cs.core.ActiveCharacter(),
		).
			Write("src", cs.src).
			Write("expiry", cs.expiry).
			Write("earliest_pickup", cs.EarliestPickup)
		return false
	}
	cs.core.Log().NewEvent(
		fmt.Sprintf("%v crystallize shard picked up", cs.Shield.Element()),
		glog.LogElementEvent,
		cs.core.ActiveCharacter(),
	).
		Write("src", cs.src).
		Write("expiry", cs.expiry).
		Write("earliest_pickup", cs.EarliestPickup)
	// add shield
	cs.Shield.SetExpiry(cs.core.F() + 15.1*60) // shield lasts for 15.1s from shard pickup
	cs.core.AddPlayerShield(cs.Shield)
	// kill self
	cs.Kill(nil)
	return true
}

func (cs *CrystallizeShard) HandleAttack(atk *info.AttackEvent) float64 {
	cs.core.Events().GadgetHit.Emit(event.TargetHitEvent{
		Target:      cs.Key(),
		AttackEvent: atk,
	})
	return 0
}

func (cs *CrystallizeShard) Attack(*info.AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (cs *CrystallizeShard) SetDirection(trg geometry.Point)                      {}
func (cs *CrystallizeShard) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}
