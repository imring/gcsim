package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

func (c *Core) PlayerTarget() core.Character {
	return c.player.PlayerTarget()
}

func (c *Core) ActiveCharacter() int {
	return c.player.ActiveCharacter()
}

func (c *Core) GetCharacter(index int) core.Character {
	return c.player.GetCharacter(index)
}

func (c *Core) GetCharacterByKey(key keys.Char) core.Character {
	return c.player.GetCharacterByKey(key)
}

func (c *Core) GetCharacterByTarget(key keys.Target) core.Character {
	return c.player.GetCharacterByTarget(key)
}

func (c *Core) SetPlayerPos(pos geometry.Point) {
	c.player.SetPlayerPos(pos)
}

func (c *Core) SetPlayerDirectionToClosestEnemy() {
	c.player.SetPlayerDirectionToClosestEnemy()
}

func (c *Core) AddPlayerShield(shd core.Shield) {
	c.player.Shield.Add(shd)
}

func (c *Core) PlayerSwapCD() int {
	return c.player.SwapCD()
}

func (c *Core) SetPlayerSwapICD(cd int) {
	c.player.SetSwapICD(cd)
}

func (c *Core) PlayerDashLockout() bool {
	return c.player.DashLockout()
}

func (c *Core) PlayerRemainingDashCD() int {
	return c.player.RemainingDashCD()
}

func (c *Core) SetPlayerDashCD(lockout bool, cd int) {
	c.player.SetDashCD(lockout, cd)
}

func (c *Core) PlayerUseStam(amount float64, a info.Action) {
	c.player.UseStam(amount, a)
}

func (c *Core) PlayerStam() float64 {
	return c.player.Stam()
}

func (c *Core) PlayerCurrentState() info.AnimationState {
	return c.player.Animation.CurrentState()
}

func (c *Core) GetPlayerLastAction() info.LastAction {
	return c.player.LastAction()
}

func (c *Core) PlayerAirborne() info.AirborneSource {
	return c.player.Airborne()
}

func (c *Core) SetPlayerAirborne(a info.AirborneSource) {
	c.player.SetAirborne(a)
}
