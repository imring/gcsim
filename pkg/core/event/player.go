package event

import "github.com/genshinsim/gcsim/pkg/core/info"

type EnergyChangeEventHandler = EventHandler[EnergyChangeEvent]
type EnergyChangeEvent struct {
	CharIndex  int
	PreEnergy  float64
	Amount     float64
	Source     string
	IsParticle bool
}

type HPDebtEventHandler = EventHandler[HPDebtEvent]
type HPDebtEvent struct {
	CharIndex int
	Amount    float64 // positive = clear, negative = debt
}

type CharacterActionEventHandler = EventHandler[CharacterActionEvent]
type CharacterActionEvent struct {
	CharIndex int
	Action    info.Action
	Params    map[string]int
}

// TODO: unusable?
type StamUseEventHandler = EventHandler[StamUseEvent]
type StamUseEvent struct {
	Action info.Action
}

type StateChangeEventHandler = EventHandler[StateChangeEvent]
type StateChangeEvent struct {
	Prev info.AnimationState
	Next info.AnimationState
}

type ActionFailedEventHandler = EventHandler[ActionFailedEvent]
type ActionFailedEvent struct {
	CharIndex int
	Action    info.Action
	Params    map[string]int
	Reason    info.Failure
}

type CharacterSwapEventHandler = EventHandler[CharacterSwapEvent]
type CharacterSwapEvent struct {
	Prev int
	Next int
}

type ActionExecEventHandler = EventHandler[ActionExecEvent]
type ActionExecEvent struct {
	CharIndex int
	Action    info.Action
	Params    map[string]int
}

type ActionEventHandler = EventHandler[ActionEvent]
type ActionEvent struct {
	CharIndex int
	Params    map[string]int
}
