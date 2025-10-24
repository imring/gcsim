package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/modifier"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

const MaxTeamSize = 4

type Target struct {
	core core.Core

	key             keys.Target
	Hitbox          geometry.Circle
	Tags            map[string]int
	CollidableTypes [info.TargettableTypeCount]bool
	OnCollision     func(core.Target)

	BaseStats attributes.Props
	Modifiers *modifier.Handler

	Alive bool

	// icd related
	icdTagOnTimer       [MaxTeamSize][info.ICDTagLength]bool
	icdTagCounter       [MaxTeamSize][info.ICDTagLength]int
	icdDamageTagOnTimer [MaxTeamSize][info.ICDTagLength]bool
	icdDamageTagCounter [MaxTeamSize][info.ICDTagLength]int

	direction geometry.Point
}

func NewTarget(core core.Core, p geometry.Point, r float64) *Target {
	t := &Target{}
	t.core = core
	t.Hitbox = *geometry.NewCircle(p, r, geometry.DefaultDirection(), 360)
	t.direction = geometry.DefaultDirection()
	t.Tags = make(map[string]int)
	t.Alive = true

	t.Modifiers = modifier.New(core, t.key)

	return t
}

func (t *Target) Collidable() bool                           { return t.OnCollision != nil }
func (t *Target) CollidableWith(x info.TargettableType) bool { return t.CollidableTypes[x] }
func (t *Target) CollidedWith(x core.Target) {
	if t.OnCollision != nil {
		t.OnCollision(x)
	}
}

func (t *Target) Key() keys.Target { return t.key }
func (t *Target) SetKey(x keys.Target) {
	t.key = x
	t.Modifiers.SetOwner(x)
}

func (t *Target) Type() info.TargettableType { return info.TargettableEnemy }

func (t *Target) Shape() geometry.Shape      { return &t.Hitbox }
func (t *Target) SetShape(s geometry.Circle) { t.Hitbox = s }

func (t *Target) Pos() geometry.Point { return t.Hitbox.Pos() }
func (t *Target) SetPos(p geometry.Point) {
	t.Hitbox.SetPos(p)
	t.core.Events().TargetMoved.Emit(event.TargetMovedEvent{
		Target: t.key,
		Point:  p,
	})
	t.core.Log().NewEvent("target position changed", glog.LogSimEvent, -1).
		Write("key", t.key).
		Write("x", p.X).
		Write("y", p.Y)
}

func (t *Target) IsAlive() bool { return t.Alive }
func (t *Target) Kill(atk *info.AttackEvent) {
	t.Alive = false
	t.Core().Events().TargetDied.Emit(event.TargetDiedEvent{
		Target:      t.Key(),
		AttackEvent: atk,
	})
}

func (t *Target) HandleAttack(a *info.AttackEvent) float64 { return 0 }
func (t *Target) Tick()                                    {}

func (t *Target) SetTag(key string, val int) {
	t.Tags[key] = val
}

func (t *Target) GetTag(key string) int {
	return t.Tags[key]
}

func (t *Target) RemoveTag(key string) {
	delete(t.Tags, key)
}

func (t *Target) WillCollide(s geometry.Shape) bool {
	if !t.Alive {
		return false
	}
	switch v := s.(type) {
	case *geometry.Circle:
		return t.Shape().IntersectCircle(*v)
	case *geometry.Rectangle:
		return t.Shape().IntersectRectangle(*v)
	default:
		return false
	}
}

func (t *Target) AttackWillLand(a info.AttackPattern) (bool, string) {
	// shape shouldn't be nil; panic here
	if a.Shape == nil {
		panic("unexpected nil shape")
	}
	if !t.Alive {
		return false, "target dead"
	}
	// shape can't be nil now, check if type matches
	// if !a.Targets[t.typ] {
	// 	return false, "wrong type"
	// }
	// swirl aoe shouldn't hit the src of the aoe
	for _, v := range a.IgnoredKeys {
		if t.Key() == v {
			return false, "no self harm"
		}
	}

	// check if shape matches
	switch v := a.Shape.(type) {
	case *geometry.Circle:
		return t.Shape().IntersectCircle(*v), "intersect circle"
	case *geometry.Rectangle:
		return t.Shape().IntersectRectangle(*v), "intersect rectangle"
	case *SingleTarget:
		// only true if
		return v.Target == t.key, "target"
	default:
		return false, "unknown shape"
	}
}

func (t *Target) IsWithinArea(a info.AttackPattern) bool {
	return a.Shape.PointInShape(t.Pos())
}

func (t *Target) Direction() geometry.Point { return t.direction }
func (t *Target) SetDirection(trg geometry.Point) {
	src := t.Pos()
	t.direction = geometry.CalcDirection(src, trg)
	// t.Core.Combat.Log.NewEvent("set target direction", glog.LogDebugEvent, -1).
	// 	Write("src target key", t.key).
	// 	Write("srcX", src.X).
	// 	Write("srcY", src.Y).
	// 	Write("trgX", trg.X).
	// 	Write("trgY", trg.Y).
	// 	Write("direction", t.direction)
}

func (t *Target) CalcTempDirection(trg geometry.Point) geometry.Point {
	src := t.Pos()
	direction := geometry.CalcDirection(src, trg)
	// t.Core.Combat.Log.NewEvent("using temporary target direction", glog.LogDebugEvent, -1).
	// 	Write("src target key", t.key).
	// 	Write("srcX", src.X).
	// 	Write("srcY", src.Y).
	// 	Write("trgX", trg.X).
	// 	Write("trgY", trg.Y).
	// 	Write("direction", t.direction).
	// 	Write("temporary direction", direction)
	return direction
}

func (t *Target) Stats() attributes.Stats {
	return attributes.NewStats(t.BaseStats, t.Modifiers.Stats())
}

func (t *Target) HasModifier(key string) bool {
	return t.Modifiers.HasModifier(key)
}

func (t *Target) Core() core.Core {
	return t.core
}
