package simulation

func (c *Core) IsDamageMode() bool {
	return c.flags.DamageMode
}

func (c *Core) IsHitlagEnabled() bool {
	return c.flags.EnableHitlag
}

func (c *Core) CanBeDefenseHalt() bool {
	return c.flags.DefHalt
}

func (c *Core) IgnoreBurstEnergy() bool {
	return c.flags.IgnoreBurstEnergy
}
