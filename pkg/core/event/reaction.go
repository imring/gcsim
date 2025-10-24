package event

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type ReactionEventHandler = EventHandler[ReactionEvent]
type ReactionEvent struct {
	Target      keys.Target
	AttackEvent *info.AttackEvent
}

type AuraDurabilityAddedEventHandler = EventHandler[AuraDurabilityAddedEvent]
type AuraDurabilityAddedEvent struct {
	Target     keys.Target
	Element    attributes.ElementType
	Durability info.Durability
}

type AuraDurabilityDepletedEventHandler = EventHandler[AuraDurabilityDepletedEvent]
type AuraDurabilityDepletedEvent struct {
	Target  keys.Target
	Element attributes.ElementType
}

type DendroCoreSpawnedEventHandler = EventHandler[DendroCoreSpawnedEvent]
type DendroCoreSpawnedEvent struct {
	Target      keys.Target
	AttackEvent *info.AttackEvent
}
