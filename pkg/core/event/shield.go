package event

import "github.com/genshinsim/gcsim/pkg/core/info"

type ShieldedEventHandler = EventHandler[ShieldedEvent]
type ShieldedEvent struct {
	Type info.ShieldType
	// TODO: interface/index?
}

type ShieldBreakEventHandler = EventHandler[ShieldBreakEvent]
type ShieldBreakEvent struct {
	Type info.ShieldType
	// TODO: interface/index?
}
