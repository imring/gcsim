package combat

import (
	"math/rand"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/target"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

type Opt struct {
	F      *int
	Task   task.Tasker
	Events *event.System
	Player *player.Handler
	Target *target.Handler
	Log    glog.Logger
	Rand   *rand.Rand

	EnableHitlag bool
	DefHalt      bool
}

type Handler struct {
	f      *int
	task   task.Tasker
	events *event.System
	player *player.Handler
	target *target.Handler
	log    glog.Logger
	rand   *rand.Rand

	hitlag  bool
	defHalt bool

	enemies []core.Enemy
	gadgets []core.Gadget

	defaultTarget keys.Target

	gccount int
}

func New(opt Opt) *Handler {
	h := &Handler{
		f:       opt.F,
		task:    opt.Task,
		events:  opt.Events,
		player:  opt.Player,
		target:  opt.Target,
		log:     opt.Log,
		rand:    opt.Rand,
		hitlag:  opt.EnableHitlag,
		defHalt: opt.DefHalt,
	}
	h.enemies = make([]core.Enemy, 0, 5)
	h.gadgets = make([]core.Gadget, 0, 10)

	return h
}

func (h *Handler) SetDefaultTarget(key keys.Target) { h.defaultTarget = key }
func (h *Handler) DefaultTarget() keys.Target       { return h.defaultTarget }

func (h *Handler) Tick() {
	// collision check happens before each object ticks (as collision may remove the object)
	// enemy and player does not check for collision
	// gadgets check against player and enemy

	player := h.player.PlayerTarget()
	for i := 0; i < len(h.gadgets); i++ {
		if h.gadgets[i] != nil && h.gadgets[i].CollidableWith(info.TargettablePlayer) {
			if h.gadgets[i].WillCollide(player.Shape()) {
				h.gadgets[i].CollidedWith(player)
			}
		}
		// sanity check in case gadget is gone
		if h.gadgets[i] != nil && h.gadgets[i].CollidableWith(info.TargettableEnemy) {
			for j := 0; j < len(h.enemies) && h.gadgets[i] != nil; j++ {
				if h.gadgets[i].WillCollide(h.enemies[j].Shape()) {
					h.gadgets[i].CollidedWith(h.enemies[j])
				}
			}
		}
	}

	for _, v := range h.enemies {
		if v != nil {
			v.Tick()
		}
	}
	for _, v := range h.gadgets {
		if v != nil {
			v.Tick()
		}
	}

	// TODO: clean up every 100 tick reasonable?
	h.gccount++
	if h.gccount > 100 {
		n := 0
		for i, v := range h.gadgets {
			if v != nil {
				h.gadgets[n] = h.gadgets[i]
				n++
			}
		}
		h.gadgets = h.gadgets[:n]
		h.gccount = 0
	}
}
