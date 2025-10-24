package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryAggravate(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}

	if r.Durability[attributes.ElementOverdose] < ZeroDur {
		return false
	}

	r.core.Events().Aggravate.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// em isn't snapshot
	char := r.core.GetCharacter(a.Info.ActorIndex)
	a.Info.Catalyzed = true
	a.Info.CatalyzedType = info.ReactionTypeAggravate
	a.Info.FlatDmg += 1.15 * r.calcCatalyzeDmg(&a.Info, char)
	return true
}

func (r *Reactable) TrySpread(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}

	if r.Durability[attributes.ElementOverdose] < ZeroDur {
		return false
	}

	r.core.Events().Spread.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// em isn't snapshot
	char := r.core.GetCharacter(a.Info.ActorIndex)
	a.Info.Catalyzed = true
	a.Info.CatalyzedType = info.ReactionTypeSpread
	a.Info.FlatDmg += 1.25 * r.calcCatalyzeDmg(&a.Info, char)
	return true
}

func (r *Reactable) TryQuicken(a *info.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}

	var consumed info.Durability
	switch a.Info.Element {
	case attributes.ElementGrass:
		if r.Durability[attributes.ElementElectric] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.ElementElectric, a.Info.Durability, 1)
	case attributes.ElementElectric:
		if r.Durability[attributes.ElementGrass] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.ElementGrass, a.Info.Durability, 1)
	default:
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	r.core.Events().Quicken.Emit(event.ReactionEvent{
		Target:      r.self.Key(),
		AttackEvent: a,
	})

	// attach quicken aura; special amount
	r.attachQuicken(consumed)

	if r.Durability[attributes.ElementWater] >= ZeroDur {
		r.core.Tasks().Add(func() {
			r.tryQuickenBloom(a)
		}, 0)
	}

	return true
}

func (r *Reactable) attachQuicken(dur info.Durability) {
	r.attachOverlapRefreshDuration(attributes.ElementOverdose, dur, 12*dur+360)
}
