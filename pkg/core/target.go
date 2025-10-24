package core

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

type Target interface {
	Key() keys.Target           // unique key for the target
	SetKey(k keys.Target)       // update key
	Type() info.TargettableType // type of target
	Shape() geometry.Shape      // geometry.Shape of target
	SetShape(s geometry.Circle) // update shape
	Pos() geometry.Point        // center of target
	SetPos(p geometry.Point)    // move target
	IsAlive() bool
	SetTag(key string, val int)
	GetTag(key string) int
	RemoveTag(key string)
	HandleAttack(a *info.AttackEvent) float64
	AttackWillLand(a info.AttackPattern) (bool, string) // hurtbox collides with AttackPattern
	IsWithinArea(a info.AttackPattern) bool             // center is in AttackPattern
	Tick()                                              // called every tick
	Kill(atk *info.AttackEvent)
	// for collision check
	CollidableWith(info.TargettableType) bool
	CollidedWith(t Target)
	WillCollide(geometry.Shape) bool
	// direction related
	Direction() geometry.Point                           // returns viewing direction as a geometry.Point
	SetDirection(trg geometry.Point)                     // calculates viewing direction relative to default direction (0, 1)
	CalcTempDirection(trg geometry.Point) geometry.Point // used for stuff like Bow CA
	// stats
	Stats() attributes.Stats
	HasModifier(key string) bool
}

type TargetWithAura interface {
	Target
	AuraContains(e ...attributes.ElementType) bool
}
