package mika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// based on sayu frames
// TODO: update frames & hitlags
var burstFrames []int

const (
	burstHitmark = 12
	healKey      = "eagleplume"
	healIcdKey   = "eagleplume-icd"
)

func init() {
	burstFrames = frames.InitAbilSlice(65) // Q -> N1/E/J
	burstFrames[action.ActionDash] = 64    // Q -> D
	burstFrames[action.ActionSwap] = 64    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// first heal
	c.QueueCharTask(func() {
		heal := burstHealFirstF[c.TalentLvlBurst()] + burstHealFirstP[c.TalentLvlBurst()]*c.MaxHP()
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Skyfeather Song",
			Src:     heal,
			Bonus:   c.Stat(attributes.Heal),
		})

		if c.Base.Cons >= 4 {
			c.c4Count = 5
		}
		c.AddStatus(healKey, 15*60, true)
		c.DeleteStatus(healIcdKey)
	}, burstHitmark)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(7)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) onBurstHeal() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		if !c.StatusIsActive(healKey) {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}
		if c.StatusIsActive(healIcdKey) {
			return false
		}

		heal := burstHealF[c.TalentLvlBurst()] + burstHealP[c.TalentLvlBurst()]*c.MaxHP()
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Eagleplume",
			Src:     heal,
			Bonus:   c.Stat(attributes.Heal),
		})

		// When Mika's own Skyfeather Song's Eagleplume state heals party members, this will restore 3 Energy to Mika.
		// This form of Energy restoration can occur 5 times during the Eagleplume state created by 1 use of Skyfeather Song.
		if c.Base.Cons >= 4 && c.c4Count > 0 {
			c.AddEnergy("mika-c4", 3)
			c.c4Count--
		}

		c.AddStatus(healIcdKey, c.healIcd, true)

		return false
	}, "mika-eagleplume")
}
