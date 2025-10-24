package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryMelt(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.ElementFire:
		if r.Durability[attributes.ElementIce] < ZeroDur && r.Durability[attributes.ElementFrozen] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.ElementIce, a.Info.Durability, 2)
		f := r.reduce(attributes.ElementFrozen, a.Info.Durability, 2)
		if f > consumed {
			consumed = f
		}
		a.Info.AmpMult = 2.0
	case attributes.ElementIce:
		if r.Durability[attributes.ElementFire] < ZeroDur && r.Durability[attributes.ElementBurning] < ZeroDur {
			return false
		}
		r.reduce(attributes.ElementFire, a.Info.Durability, 0.5)
		a.Info.AmpMult = 1.5
		r.burningCheck()
	default:
		// should be here
		return false
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true
	a.Info.Amped = true
	a.Info.AmpType = info.ReactionTypeMelt
	r.core.Events().Melt.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})
	return true
}
