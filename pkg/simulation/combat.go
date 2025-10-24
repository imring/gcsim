package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

func (c *Core) QueueAttackWithSnap(snap info.Snapshot, qa info.QueueAttack) {
	c.combat.QueueAttackWithSnap(snap, qa)
}

func (c *Core) QueueAttackEvent(ae *info.AttackEvent, dmgDelay int) {
	c.combat.QueueAttackEvent(ae, dmgDelay)
}

func (c *Core) QueueAttack(qa info.QueueAttack) {
	c.combat.QueueAttack(qa)
}

func (c *Core) GetEnemy(i int) core.Enemy {
	return c.combat.Enemy(i)
}

func (c *Core) GetEnemies() []core.Enemy {
	return c.combat.Enemies()
}

func (c *Core) SetEnemyPos(i int, pos geometry.Point) {
	c.combat.SetEnemyPos(i, pos)
}

func (c *Core) KillEnemy(i int) {
	c.combat.KillEnemy(i)
}

func (c *Core) AddGadget(g core.Gadget) {
	c.combat.AddGadget(g)
}

func (c *Core) RemoveGadget(key keys.Target) {
	c.combat.RemoveGadget(key)
}

func (c *Core) GetGadgets() []core.Gadget {
	return c.combat.GetGadgets()
}

func (c *Core) SetDefaultTarget(key keys.Target) {
	c.combat.SetDefaultTarget(key)
}

func (c *Core) ClosestEnemy(pos geometry.Point) core.Enemy {
	return c.combat.ClosestEnemy(pos)
}

func (c *Core) ClosestGadget(pos geometry.Point) core.Gadget {
	return c.combat.ClosestGadget(pos)
}

func (c *Core) ClosestEnemyWithinArea(ap info.AttackPattern, filter func(t core.Enemy) bool) core.Enemy {
	return c.combat.ClosestEnemyWithinArea(ap, filter)
}
