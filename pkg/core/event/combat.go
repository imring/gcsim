package event

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

type ApplyAttackEventHandler = EventHandler[*info.AttackEvent]

type TargetHitEventHandler = EventHandler[TargetHitEvent]
type TargetHitEvent struct {
	Target      keys.Target
	AttackEvent *info.AttackEvent
}

type EnemyDamageEventHandler = EventHandler[EnemyDamageEvent]
type EnemyDamageEvent struct {
	Target      keys.Target
	AttackEvent *info.AttackEvent
	Damage      float64
	IsCrit      bool
}

type TargetMovedEventHandler = EventHandler[TargetMovedEvent]
type TargetMovedEvent struct {
	Target keys.Target
	Point  geometry.Point
}
