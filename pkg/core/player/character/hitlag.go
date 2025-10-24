package character

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *Character) QueueCharTask(f func(), delay int) {
	if delay == 0 {
		f()
		return
	}
	c.queue.Add(f, delay)
}

func (c *Character) Tick() {
	// decrement frozen time first
	c.frozenFrames -= 1
	left := 0
	if c.frozenFrames < 0 {
		left = -c.frozenFrames
		c.frozenFrames = 0
	}

	c.Modifiers.Tick(c.FramePausedOnHitlag())

	// if any left then increase time passed
	if left <= 0 {
		// do nothing this tick
		return
	}
	c.TimePassed += left

	// check char queue for any executable actions
	c.queue.Run()
}

func (c *Character) FramePausedOnHitlag() bool {
	return c.frozenFrames > 0
}

// ApplyHitlag adds hitlag to the character for specified duration
func (c *Character) ApplyHitlag(factor, dur float64) {
	// number of frames frozen is total duration * (1 - factor)
	ext := int(math.Ceil(dur * (1 - factor)))
	c.frozenFrames += ext
	c.core.Log().NewEvent(
		fmt.Sprintf("hitlag applied to char: %.3f", dur),
		glog.LogHitlagEvent, c.GetIndex(),
	).
		Write("duration", dur).
		Write("factor", factor).
		Write("frozen_frames", c.frozenFrames).
		SetEnded(c.core.F() + int(math.Ceil(dur)))

}
