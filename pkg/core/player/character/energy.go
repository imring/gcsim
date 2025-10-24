package character

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *Character) Energy() float64 {
	return c.energy
}

func (c *Character) EnergyMax() float64 {
	return c.energyMax
}

func (c *Character) AddEnergy(src string, amt float64) {
	preEnergy := c.energy
	c.energy += amt
	if c.energy > c.energyMax {
		c.energy = c.energyMax
	}
	if c.energy < 0 {
		c.energy = 0
	}

	c.core.Events().EnergyChange.Emit(event.EnergyChangeEvent{
		CharIndex:  c.Index,
		PreEnergy:  preEnergy,
		Amount:     amt,
		Source:     src,
		IsParticle: false,
	})

	c.core.Log().NewEvent("adding energy", glog.LogEnergyEvent, c.Index).
		Write("rec'd", amt).
		Write("pre_recovery", preEnergy).
		Write("post_recovery", c.energy).
		Write("source", src).
		Write("max_energy", c.energyMax)
}

func (c *Character) SetParticleDelay(delay int) {
	c.ParticleDelay = delay
}
