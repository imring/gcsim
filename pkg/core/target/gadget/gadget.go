// gadget provides a skeleton implementation for gadgets and is
// basically a thin wrapper around target.Target and adds some
// helper functionalities
package gadget

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/target"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

type Gadget struct {
	*target.Target

	core            core.Core
	src             int
	gadgetTyp       info.GadgetTyp
	OnKill          func()
	OnExpiry        func() // only called if gadget dies from expiry
	ThinkInterval   int    // should be > 0
	OnThinkInterval func()
	duration        int // how long gadget should live for; use -1 for infinite
	// internal helper
	sinceLastThink int
}

func New(core core.Core, p geometry.Point, r float64, typ info.GadgetTyp) *Gadget {
	g := &Gadget{
		core:      core,
		src:       core.F(),
		gadgetTyp: typ,
	}
	g.Target = target.NewTarget(core, p, r)
	return g
}

func (g *Gadget) Kill(atk *info.AttackEvent) {
	g.Target.Kill(atk)
	if g.OnKill != nil {
		g.OnKill()
	}
	g.core.RemoveGadget(g.Key())
}

func (g *Gadget) Type() info.TargettableType { return info.TargettableGadget }
func (g *Gadget) Src() int                   { return g.src }
func (g *Gadget) GadgetTyp() info.GadgetTyp  { return g.gadgetTyp }

func (g *Gadget) Duration() int     { return g.duration }
func (g *Gadget) SetDuration(d int) { g.duration = d }

func (g *Gadget) Tick() {
	if g.OnThinkInterval != nil && g.ThinkInterval > 0 {
		if g.sinceLastThink < g.ThinkInterval {
			g.sinceLastThink++
		}
		if g.sinceLastThink == g.ThinkInterval {
			g.OnThinkInterval()
			g.sinceLastThink = 0
		}
	}
	if g.duration != -1 && g.duration > 0 {
		g.duration--
		if g.duration == 0 {
			if g.OnExpiry != nil {
				g.OnExpiry()
			}
			g.core.RemoveGadget(g.Key())
		}
	}
}

func (g *Gadget) SetOnKill(fn func())           { g.OnKill = fn }
func (g *Gadget) SetOnExpiry(fn func())         { g.OnExpiry = fn }
func (g *Gadget) SetThinkInterval(interval int) { g.ThinkInterval = interval }
func (g *Gadget) SetOnThinkInterval(fn func())  { g.OnThinkInterval = fn }
