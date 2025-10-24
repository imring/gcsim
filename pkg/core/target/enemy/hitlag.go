package enemy

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (e *Enemy) ApplyHitlag(factor, dur float64) {
	// TODO: extend all hitlag affected buff expiry by dur * (1 - factor) i think
	ext := int(math.Ceil(dur * (1 - factor)))
	e.frozenFrames += ext

	e.Core().Log().NewEvent("enemy hitlag - extending mods", glog.LogHitlagEvent, -1).
		Write("target", e.Key()).
		Write("duration", dur).
		Write("factor", factor).
		Write("frozen_frames", e.frozenFrames).
		SetEnded(e.Core().F() + int(math.Ceil(dur)))

	e.Modifiers.ExtendByHitlag(ext)
}

func (e *Enemy) QueueEnemyTask(f func(), delay int) {
	if delay == 0 {
		f()
		return
	}
	e.queue.Add(f, delay)
}

func (e *Enemy) Tick() {
	// dead enemy don't tick
	if !e.IsAlive() {
		return
	}

	// decrement frozen time first
	e.frozenFrames -= 1
	left := 0
	if e.frozenFrames < 0 {
		left = -e.frozenFrames
		e.frozenFrames = 0
	}
	// if any left then increase time passed
	if left <= 0 {
		e.Core().Log().NewEvent("enemy skipping tick", glog.LogHitlagEvent, -1).
			Write("target", e.Key()).
			Write("frozen_for", e.frozenFrames)
		// do nothing this tick
		return
	}
	e.timePassed += left

	e.queue.Run()
	e.Reactable.Tick()
}
