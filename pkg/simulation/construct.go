package simulation

import "github.com/genshinsim/gcsim/pkg/core/info"

func (c *Core) ConstructCountByType(t info.GeoConstructType) int {
	return c.construct.CountByType(t)
}

func (c *Core) ConstructExpiry(t info.GeoConstructType) int {
	return c.construct.Expiry(t)
}
