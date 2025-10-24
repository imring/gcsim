package event

type ConstructSpawnedEventHandler = EventHandler[ConstructSpawnedEvent]
type ConstructSpawnedEvent struct {
	Index int
}
