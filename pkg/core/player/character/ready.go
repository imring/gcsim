package character

import (
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (c *Character) ActionReady(a info.Action, p map[string]int) (bool, info.Failure) {
	// for dash and charge need to check for stam usage as well

	switch a {
	case info.ActionBurst:
		if !c.core.IgnoreBurstEnergy() && c.Energy() != c.EnergyMax() {
			return false, info.InsufficientEnergy
		}
		if c.AvailableCDCharge[a] <= 0 {
			return false, info.BurstCD
		}
	case info.ActionSkill:
		if c.AvailableCDCharge[a] <= 0 {
			return false, info.SkillCD
		}
	case info.ActionCharge:
		req := c.AbilStamCost(a, p)
		if c.core.PlayerStam() < req {
			c.core.Log().NewEvent("insufficient stam: charge attack", glog.LogWarnings, -1).
				Write("have", c.core.PlayerStam())
			return false, info.InsufficientStamina
		}
	case info.ActionDash:
		req := c.AbilStamCost(a, p)
		if c.core.PlayerStam() < req {
			c.core.Log().NewEvent("insufficient stam: dash", glog.LogWarnings, -1).
				Write("have", c.core.PlayerStam())
			return false, info.InsufficientStamina
		}
		if c.core.ActiveCharacter() == c.GetIndex() && c.core.PlayerDashLockout() && c.core.PlayerRemainingDashCD() > 0 {
			c.core.Log().NewEvent("dash on cooldown", glog.LogWarnings, -1).
				Write("dash_cd_expiration", c.core.PlayerRemainingDashCD())
			return false, info.DashCD
		}
		if c.core.ActiveCharacter() != c.GetIndex() && c.DashLockout() && c.RemainingDashCD() > 0 {
			return false, info.DashCD
		}
	}
	return true, info.NoFailure
}
