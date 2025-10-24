package modifier

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Instance struct {
	key               string
	source            keys.Target
	duration          int
	thinkInterval     int
	elementDurability info.Durability
	decayRate         info.Durability
	hitlagAffected    bool
	count             int
	maxCount          int
	state             any
	props             attributes.Props

	currentThinkInterval int // TODO: use info.Durability?
	handler              *Handler
	listeners            Listeners
}

func (i *Instance) Key() string {
	return i.key
}

func (i *Instance) Owner() keys.Target {
	return i.handler.Owner()
}

func (i *Instance) Source() keys.Target {
	return i.source
}

func (i *Instance) Duration() int {
	if i.decayRate <= 0 {
		return -1
	}
	return int(i.elementDurability / i.decayRate)
}

func (i *Instance) HitlagAffected() bool {
	return i.hitlagAffected
}

func (i *Instance) Count() int {
	return i.count
}

func (i *Instance) MaxCount() int {
	return i.maxCount
}

func (i *Instance) State() any {
	return i.state
}

func (i *Instance) Props() attributes.Props {
	return i.props
}

func (i *Instance) Handler() *Handler {
	return i.handler
}

func (i *Instance) SetProperty(prop attributes.Prop, value float64) {
	i.props[prop] = value
}

func (i *Instance) SetProperties(props attributes.Props) {
	i.props = props
}

func (i *Instance) Core() core.Core {
	return i.handler.Core()
}

func (i *Instance) updateDuration(dur int) {
	i.duration = dur
	i.elementDurability = info.Durability(dur) * i.decayRate
}
