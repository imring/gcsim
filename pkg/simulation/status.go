package simulation

func (c *Core) StatusDuration(status string) int {
	return c.status.Duration(status)
}
