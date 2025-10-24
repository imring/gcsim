package character

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func (c *Character) Attack(p map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action attack not implemented", c.Base.Key)
}
func (c *Character) Aimed(p map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action aimed not implemented", c.Base.Key)
}
func (c *Character) ChargeAttack(p map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action charge not implemented", c.Base.Key)
}
func (c *Character) HighPlungeAttack(p map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action high_plunge not implemented", c.Base.Key)
}
func (c *Character) LowPlungeAttack(p map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action low_plunge not implemented", c.Base.Key)
}
func (c *Character) Skill(p map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action skill not implemented", c.Base.Key)
}
func (c *Character) Burst(p map[string]int) (action.Info, error) {
	return action.Info{}, fmt.Errorf("%v: action burst not implemented", c.Base.Key)
}

func (c *Character) NextQueueItemIsValid(_ keys.Char, a info.Action, p map[string]int) error {
	if a == info.ActionCharge {
		switch c.Weapon.Class {
		case info.WeaponClassSword, info.WeaponClassSpear:
			// cannot do charge on most sword/polearm characters without attack beforehand
			if c.core.GetPlayerLastAction().Type != info.ActionAttack {
				return info.ErrInvalidChargeAction
			}
		}
	}
	return nil
}

func (c *Character) ResetNormalCounter() {
	c.normalCounter = 0
}

func (c *Character) AdvanceNormalIndex() {
	c.normalCounter++
	if c.normalCounter == c.normalHitNum {
		c.normalCounter = 0
	}
}

func (c *Character) NormalCounter() int     { return c.normalCounter }
func (c *Character) NextNormalCounter() int { return c.normalCounter + 1 }
func (c *Character) DashLockout() bool      { return c.dashLockout }
func (c *Character) RemainingDashCD() int   { return c.remainingDashCD }

func (c *Character) SetDashCD(cd int, lockout bool) {
	c.remainingDashCD = cd
	c.dashLockout = lockout
}
