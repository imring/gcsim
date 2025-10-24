package combat

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (h *Handler) ApplyAttack(ae *info.AttackEvent) {
	h.events.ApplyAttack.Emit(ae)

	var landed bool

	for _, v := range h.target.Targets() {
		if v == nil {
			continue
		}
		if ae.Pattern.SkipTargets[v.Type()] || !v.IsAlive() {
			continue
		}
		_, land := h.attack(v, ae)
		if land && v.Type() == info.TargettableEnemy {
			landed = true
		}
	}

	// add hitlag to actor but ignore if this is deployable
	if h.hitlag && landed && !ae.Info.IsDeployable {
		dur := ae.Info.HitlagHaltFrames
		if h.defHalt && ae.Info.CanBeDefenseHalted {
			dur += 0.06 * 60
		}
		if dur > 0 {
			h.player.ApplyHitlag(ae.Info.ActorIndex, ae.Info.HitlagFactor, dur)
			h.log.NewEvent(fmt.Sprintf("%v applying hitlag: %.3f", ae.Info.Abil, dur), glog.LogHitlagEvent, ae.Info.ActorIndex).
				Write("duration", dur).
				Write("factor", ae.Info.HitlagFactor)
		}
	}
}

func (h *Handler) QueueAttackWithSnap(snap info.Snapshot, qa info.QueueAttack) {
	if qa.DmgDelay < 0 {
		panic("dmgDelay cannot be less than 0")
	}
	ae := info.AttackEvent{
		Info:        qa.Info,
		Pattern:     qa.Pattern,
		Attacker:    snap,
		SourceFrame: *h.f,
	}
	// add callbacks only if not nil
	for _, f := range qa.Callbacks {
		if f != nil {
			ae.Callbacks = append(ae.Callbacks, f)
		}
	}
	h.queueDmg(&ae, qa.DmgDelay)
}

func (h *Handler) QueueAttackEvent(ae *info.AttackEvent, dmgDelay int) {
	h.queueDmg(ae, dmgDelay)
}

func (h *Handler) QueueAttack(qa info.QueueAttack) {
	// panic if dmgDelay < snapshotDelay; this should not happen. if it happens then there's something wrong with the
	// character's code
	if qa.DmgDelay < qa.SnapshotDelay {
		panic("dmgDelay cannot be less than snapshotDelay")
	}
	if qa.DmgDelay < 0 {
		panic("dmgDelay cannot be less than 0")
	}
	// create attackevent
	ae := info.AttackEvent{
		Info:        qa.Info,
		Pattern:     qa.Pattern,
		SourceFrame: *h.f,
	}

	// add callbacks only if not nil
	for _, f := range qa.Callbacks {
		if f != nil {
			ae.Callbacks = append(ae.Callbacks, f)
		}
	}

	switch {
	case qa.SnapshotDelay < 0:
		// snapshotDelay < 0 means we don't need a snapshot; optimization for reaction
		// damage essentially
		h.queueDmg(&ae, qa.DmgDelay)
	case qa.SnapshotDelay == 0:
		h.generateSnapshot(&ae)
		h.queueDmg(&ae, qa.DmgDelay)
	default:
		// use add task ctrl to queue; no need to track here
		h.task.Add(func() {
			h.generateSnapshot(&ae)
			h.queueDmg(&ae, qa.DmgDelay-qa.SnapshotDelay)
		}, qa.SnapshotDelay)
	}

}

// This code here should probably be handled in player not core
// since it's a convenience function wrapped around queuedamage
//
// does it make sense for core to have any knowledge of teams? probably not??
func (h *Handler) generateSnapshot(ae *info.AttackEvent) {
	ae.Attacker = h.player.ByIndex(ae.Info.ActorIndex).Snapshot(&ae.Info)
}

func (h *Handler) queueDmg(ae *info.AttackEvent, delay int) {
	if delay == 0 {
		h.ApplyAttack(ae)
		return
	}
	h.task.Add(func() {
		h.ApplyAttack(ae)
	}, delay)
}

func (h *Handler) attack(t core.Target, ae *info.AttackEvent) (float64, bool) {
	willHit, _ := t.AttackWillLand(ae.Pattern)
	if !willHit {
		// Move target logs into the "Sim" event log to avoid cluttering main display for stuff like Guoba
		// And obvious things like "Fischl A4 is single target so it didn't hit targets 2-4"
		// TODO: Maybe want to add a separate set of log events for this?
		// if h.Debug && t.Type() != targets.TargettablePlayer {
		// 	h.Log.NewEventBuildMsg(glog.LogDebugEvent, a.Info.ActorIndex, "skipped ", a.Info.Abil, " ", reason).
		// 		Write("attack_tag", a.Info.AttackTag).
		// 		Write("applied_ele", a.Info.Element).
		// 		Write("dur", a.Info.Durability).
		// 		Write("target", t.Key()).
		// 		Write("geometry.Shape", a.Pattern.Shape.String())
		// }
		return 0, false
	}

	dmg := t.HandleAttack(ae)
	return dmg, true
}
