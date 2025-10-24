package frames

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func InitNormalCancelSlice(hitmark, animation int) []int {
	t := make([]int, info.EndActionType)
	for i := range t {
		t[i] = animation
	}
	t[info.ActionAim] = hitmark
	t[info.ActionSkill] = hitmark
	t[info.ActionBurst] = hitmark
	t[info.ActionDash] = hitmark
	t[info.ActionJump] = hitmark
	t[info.ActionSwap] = hitmark
	return t
}

func InitAbilSlice(animation int) []int {
	t := make([]int, info.EndActionType)
	for i := range t {
		t[i] = animation
	}
	return t
}

func AtkSpdAdjust(f int, atkspd float64) int {
	if atkspd > 0.6 {
		atkspd = 0.6
	}
	return f - int(min(atkspd, 0.1+(atkspd-0.1)/2)*float64(f))
}

func NewAttackFunc(c core.Character, slice [][]int) func(info.Action) int {
	n := c.NormalCounter()
	atkspd := c.Stat(attributes.OverallSpd) + c.Stat(attributes.AtkSpd)
	spd := c.Stat(attributes.OverallSpdMult)
	return func(next info.Action) int {
		return int(float64(AtkSpdAdjust(slice[n][next], atkspd)) * spd)
	}
}

func NewAbilFunc(slice []int) func(info.Action) int {
	return func(next info.Action) int {
		return slice[next]
	}
}
