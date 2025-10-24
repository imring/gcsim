package event

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type TargetDiedEventHandler = EventHandler[TargetDiedEvent]
type TargetDiedEvent struct {
	Target      keys.Target
	AttackEvent *info.AttackEvent
}
