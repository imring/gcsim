package sourcewaterdroplet

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/target/gadget"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

type Gadget struct {
	*gadget.Gadget
}

func New(core core.Core, pos geometry.Point, typ info.GadgetTyp) *Gadget {
	p := &Gadget{}
	p.Gadget = gadget.New(core, pos, 1, typ)
	p.Gadget.SetDuration(878)
	core.AddGadget(p)
	return p
}

func (s *Gadget) HandleAttack(*info.AttackEvent) float64 { return 0 }
func (s *Gadget) SetDirection(trg geometry.Point)        {}
func (s *Gadget) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}

func (s *Gadget) Type() info.TargettableType                           { return info.TargettableGadget }
func (s *Gadget) Attack(*info.AttackEvent, glog.Event) (float64, bool) { return 0, false }
