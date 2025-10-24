// package animation provides a simple way of tracking the current
// animation state at any given frame, as well as if the current frame
// is in animation lock or not
package animation

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

type Handler struct { //nolint:revive // cannot just name this Handler because then there is a conflict with Handler in player package
	f      *int
	events *event.System
	log    glog.Logger
	tasks  task.Tasker

	char    int
	started int
	lastAct info.Action
	aniEvt  *action.Info

	state       info.AnimationState
	stateExpiry int

	event glog.Event
}

func New(f *int, log glog.Logger, events *event.System, tasks task.Tasker) *Handler {
	h := &Handler{
		f:      f,
		log:    log,
		events: events,
		tasks:  tasks,
	}
	return h
}

// IsAnimationLocked returns true if the next action can be executed on the
// current frame; false otherwise
func (h *Handler) IsAnimationLocked(next info.Action) bool {
	if h.aniEvt == nil {
		return false
	}
	// actions are always executed after ticks and right before we advance to
	// the next frame i.e. at the end of a frame
	//
	// if an action should take 20 frames of animation then this action would be ready if
	// f >= s + 20
	//
	// i.e the action lasted 20 frames counting the current frame
	// fmt.Printf("animation check; current frame %v, animation duration %v\n", *h.f, h.info.Frames(next))
	return !h.aniEvt.CanUse(next)
}

// CanQueue returns true if we can start looking for the next action to queue
// on the current frame, false otherwise
func (h *Handler) CanQueueNextAction() bool {
	if h.aniEvt == nil {
		return true
	}
	return h.aniEvt.CanQueueNext()
}

func (h *Handler) SetActionUsed(char int, act info.Action, evt *action.Info) {
	// remove previous if still active
	if h.aniEvt != nil {
		if h.aniEvt.OnRemoved != nil {
			h.aniEvt.OnRemoved(evt.State)
		}
		h.logEnded()
	}
	// setup next
	h.char = char
	h.started = *h.f
	h.aniEvt = evt
	h.events.StateChange.Emit(event.StateChangeEvent{
		Prev: h.state,
		Next: evt.State,
	})
	h.state = evt.State
	h.stateExpiry = *h.f + evt.AnimationLength
	h.lastAct = act
	h.event = h.log.NewEvent(fmt.Sprintf("%v started", act.String()), glog.LogHitlagEvent, char).
		Write("AnimationLength", evt.AnimationLength).
		Write("CanQueueAfter", evt.CanQueueAfter).
		Write("State", evt.State.String())
	for i := range info.EndActionType {
		h.event.Write(i.String(), evt.Frames(i))
	}
}

func (h *Handler) CurrentState() info.AnimationState {
	if h.aniEvt == nil {
		return info.AnimationStateIdle
	}
	return h.state
}

func (h *Handler) CurrentStateStart() int {
	return h.started
}

func (h *Handler) Tick() {
	if h.aniEvt != nil && h.aniEvt.Tick() {
		h.logEnded()
		h.events.StateChange.Emit(event.StateChangeEvent{
			Prev: h.state,
			Next: info.AnimationStateIdle,
		})
		h.state = info.AnimationStateIdle
		h.aniEvt = nil
	}
}

func (h *Handler) logEnded() {
	h.event.SetEnded(*h.f)
	h.log.NewEvent(
		fmt.Sprintf("%v from %v ended, time passed: %v (actual: %v)", h.lastAct, h.started, h.aniEvt.TimePassed, h.aniEvt.NormalizedTimePassed),
		glog.LogHitlagEvent,
		h.char,
	)
}
