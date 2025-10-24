package simulation

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/geometry"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type stateFn func(*Core) (stateFn, error)

func (c *Core) resFromCurrentState() stats.Result {
	return stats.Result{Seed: uint64(c.Seed), Duration: c.f + 1}
}

func (c *Core) run() (stats.Result, error) {
	// core loop roughly as follows:
	//  - initialize:
	//		- setup
	//		- advance frame by 1
	//		- move to queue phase
	//  - queue phase:
	//		- ask for next action
	//		- move to ready check phase
	//	- ready check phase
	//		- check if action ready (both animation + player); if not ready advance frame until ready
	//		- move to execute action phase
	//	- execute action phase:
	//		- if action has pre-action wait; advance frame until wait is consumed
	//		- execute action and empty queue
	//		- if executed action is no-op, move directly to queue phase
	//		- else advance frame until CanQueueAfter then move to queue phase
	//
	// frame advance will perform the following;
	//	- increment frame counter by 1
	//  - execute any ticks
	//  - check for eneryg procs
	//  - emit OnTick
	//  - perform exit check
	//
	// exit check checks for:
	//	- frame limit
	//  - all enemies dead
	//  - no more actions left

	// TODO: do we need to catch panic here still? or can it be done outside in the worker
	var err error
	for state := initialize; state != nil; {
		state, err = state(c)
		if err != nil {
			return c.resFromCurrentState(), err
		}
	}

	// err = c.eval.Exit()
	// if err != nil {
	// 	fmt.Println("evaluator already closed")
	// 	return c.resFromCurrentState(), err
	// }
	c.Eval.Exit()

	err = c.Eval.Err()
	if err != nil {
		return c.resFromCurrentState(), err
	}

	c.events.SimEndedSuccessfully.Emit(nil)

	return c.gatherResult(), nil
}

func (c *Core) gatherResult() stats.Result {
	res := stats.Result{
		Seed:     uint64(c.Seed),
		Duration: c.f,
		// TotalDamage: c.combat.TotalDamage,                     // TODO: get from collectors
		// DPS:         c.combat.TotalDamage * 60 / float64(c.f), // TODO: get from collectors
		Characters: make([]stats.CharacterResult, c.player.CharCount()),
		Enemies:    make([]stats.EnemyResult, c.combat.EnemyCount()),
		EndStats:   make([]stats.EndStats, c.player.CharCount()),
	}

	for i, char := range c.Config.Characters {
		res.Characters[i].Name = char.Base.Key.String()
	}

	for _, collector := range c.collectors {
		collector.Flush(c, &res)
	}

	return res
}

func (c *Core) popQueue() int {
	switch len(c.queue) {
	case 0:
	case 1:
		c.queue = c.queue[:0]
	default:
		c.queue = c.queue[1:]
	}
	return len(c.queue)
}

func initialize(c *Core) (stateFn, error) {
	// setup list
	//	- resonance
	//	- on hit energy
	//	- base stats
	//	- char inits
	//	- init call backs
	var err error

	startPoint := geometry.Point{X: c.Config.InitialPlayerPos.X, Y: c.Config.InitialPlayerPos.Y}
	err = c.setupCharacters(startPoint, c.Config.InitialPlayerPos.R, c.Config.Characters, c.Config.InitialChar)
	if err != nil {
		return nil, err
	}

	err = c.setupTargets(c.Config.Targets)
	if err != nil {
		return nil, err
	}

	// c.SetupOnNormalHitEnergy()
	err = c.player.InitializeTeam()
	if err != nil {
		return nil, err
	}

	go c.Eval.Run(c)
	// run sim for 90s if no duration set
	if c.Config.Settings.Duration == 0 {
		// fmt.Println("no duration set, running for 90s")
		c.Config.Settings.Duration = 90
	}
	c.flags.DamageMode = c.Config.Settings.DamageMode

	return c.advanceFrames(1, queuePhase)
}

func queuePhase(c *Core) (stateFn, error) {
	if c.noMoreActions {
		return c.advanceFrames(1, queuePhase)
	}
	c.Eval.Continue()
	next, err := c.Eval.NextAction()
	if err != nil {
		return nil, err
	}

	// skip a frame and come back to queue phase if eval does not have any more actions
	// relying on advance frame to exit if need be
	if next == nil {
		c.noMoreActions = true
		// we do the same skip here as if eval doesn't have any more ations
		return c.advanceFrames(1, queuePhase)
	}

	// handle sleep here since it's just a frame skip before requeing next
	if next.Action == info.ActionWait {
		return c.handleWait(next)
	}

	// if next action is delay, we can just queue up the action after that right now
	if next.Action == info.ActionDelay {
		// append here because we can have multiple delay chained
		delay := next.Param["f"]
		c.preActionDelay += delay
		c.log.NewEvent(fmt.Sprintf("delay added %v, total: %v", delay, c.preActionDelay), glog.LogActionEvent, c.player.ActiveCharacter()).
			Write("added", delay).
			Write("total", c.preActionDelay)
		return queuePhase, nil
	}

	// IMPORTANT: evaluator should handle adding in implicit swaps if next char is not active
	// we add a sanity check here just to guard against evaluator error
	active := c.player.PlayerTarget()
	if next.Char != active.GetBase().Key && next.Action != info.ActionSwap {
		return nil, fmt.Errorf("internal error: requested next char %v is not active and next action is not swap", next.Char)
	}

	// TODO: consider changing queue to single item. no need for slice without swap in here
	c.queue = append(c.queue, next)
	return actionReadyCheckPhase, nil
}

func actionReadyCheckPhase(c *Core) (stateFn, error) {
	// TODO: this sanity check is probably not necessary
	if len(c.queue) == 0 {
		return nil, errors.New("unexpected queue length is 0")
	}
	q := c.queue[0]

	// check if the next queue item is valid
	// example: most sword characters can't do charge if the previous action was not attack
	char := c.player.PlayerTarget()
	if err := char.NextQueueItemIsValid(q.Char, q.Action, q.Param); err != nil {
		switch {
		case errors.Is(err, player.ErrInvalidChargeAction):
			return nil, fmt.Errorf("%v: %w", char.GetBase().Key, player.ErrInvalidChargeAction)
		default:
			return nil, err
		}
	}

	// TODO: this loop should be optimized to skip more than 1 frame at a time
	if err := c.player.ReadyCheck(q.Action, q.Char, q.Param); err != nil {
		// repeat this phase until action is ready
		switch {
		case errors.Is(err, player.ErrActionNotReady):
			c.log.NewEvent(fmt.Sprintf("could not execute %v; action not ready", q.Action), glog.LogSimEvent, c.player.ActiveCharacter())
			return c.advanceFrames(1, actionReadyCheckPhase)
		case errors.Is(err, player.ErrPlayerNotReady):
			return c.advanceFrames(1, actionReadyCheckPhase)
		case errors.Is(err, player.ErrActionNoOp):
			// don't do anything here
		default:
			return nil, err
		}
	}

	return executeActionPhase, nil
}

func skipUntilCanQueue(c *Core) (stateFn, error) {
	if !c.player.Animation.CanQueueNextAction() {
		return c.advanceFrames(1, skipUntilCanQueue)
	}
	return queuePhase, nil
}

// nextFrame moves up the frame by 1, performing
func (c *Core) advanceFrames(f int, next stateFn) (stateFn, error) {
	for range f {
		done, err := c.nextFrame()
		if err != nil {
			return nil, err
		}
		if done {
			return nil, nil
		}
	}
	return next, nil
}

func (c *Core) handleWait(q *info.ActionEval) (stateFn, error) {
	// to maintain existing functionality, wait (alias sleep) is always ready and should cause
	// advanceFrames to be called equal to the param f
	skip := q.Param["f"]
	// log wait(0) differently to make it obvious
	if skip == 0 {
		c.log.NewEvent("executed noop wait(0)", glog.LogActionEvent, c.player.ActiveCharacter()).
			Write("f", skip)
	} else {
		c.log.NewEvent("executed wait", glog.LogActionEvent, c.player.ActiveCharacter()).
			Write("f", skip)
	}
	if l := c.popQueue(); l > 0 {
		// don't go back to queue if there are more actions already queued
		return c.advanceFrames(skip, actionReadyCheckPhase)
	}
	return c.advanceFrames(skip, queuePhase)
}

func executeActionDelay(c *Core) (stateFn, error) {
	if c.preActionDelay > 0 {
		if !c.player.PlayerTarget().FramePausedOnHitlag() {
			c.preActionDelay--
		}
		return c.advanceFrames(1, executeActionDelay)
	}
	// go back to the ready check phase in case an action becomes unavailable after delay
	return actionReadyCheckPhase, nil
}

func executeActionPhase(c *Core) (stateFn, error) {
	// TODO: this sanity check is probably not necessary
	if len(c.queue) == 0 {
		return nil, errors.New("unexpected queue length is 0")
	}
	if c.preActionDelay > 0 {
		delay := c.preActionDelay
		c.log.NewEvent(fmt.Sprintf("pre action delay: %v", delay), glog.LogActionEvent, c.player.ActiveCharacter()).
			Write("delay", delay)
		return executeActionDelay, nil
	}
	q := c.queue[0]
	err := c.player.Exec(q.Action, q.Char, q.Param)
	if err != nil {
		// TODO: this check probably doesn't do anything
		if errors.Is(err, player.ErrActionNoOp) {
			if l := c.popQueue(); l > 0 {
				// don't go back to queue if there are more actions already queued
				return actionReadyCheckPhase, nil
			}
			return queuePhase, nil
		}
		// this is now unexpected since action should be ready now
		// wrap the error for more context
		return nil, fmt.Errorf("error encountered on %v executing %v: %w", q.Char.String(), q.Action.String(), err)
	}
	// TODO: this check here is probably unnecessary
	if l := c.popQueue(); l > 0 {
		// don't go back to queue if there are more actions already queued
		return actionReadyCheckPhase, nil
	}

	return skipUntilCanQueue, nil
}

func (c *Core) nextFrame() (bool, error) {
	c.f++
	err := c.Tick()
	if err != nil {
		return false, err
	}
	c.handleEnergy()
	c.handleHurt()
	c.events.Tick.Emit(nil)
	return c.stopCheck(), nil
}

func (c *Core) stopCheck() bool {
	if c.flags.DamageMode {
		// stop if no more actions
		if c.noMoreActions {
			return true
		}
		// stop if all targets are reporting dead
		allDead := true
		for _, t := range c.combat.Enemies() {
			if t.IsAlive() {
				allDead = false
				break
			}
		}
		return allDead
	}
	return c.f == int(c.Config.Settings.Duration*60)
}

// TODO: remove defer in favour of every function actually returning error
//
//nolint:nonamedreturns // not possible to perform the res, err modification without named return
func (c *Core) Run() (res stats.Result, err error) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if r := recover(); r != nil {
			res = stats.Result{Seed: uint64(c.Seed), Duration: c.f + 1}
			err = fmt.Errorf("simulation panic occured: %v \n"+string(debug.Stack()), r)
		}
	}()
	res, err = c.run()
	return res, err
}
