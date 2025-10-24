package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/reactable"
	"github.com/genshinsim/gcsim/pkg/core/target"
	"github.com/genshinsim/gcsim/pkg/core/task"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

type Enemy struct {
	*target.Target
	*reactable.Reactable

	prof    info.EnemyProfile
	hpRatio float64

	// hitlag stuff
	timePassed   int
	frozenFrames int
	queue        *task.Handler
}

func New(core core.Core, p info.EnemyProfile) *Enemy {
	e := &Enemy{}
	e.prof = p
	e.Target = target.NewTarget(core, geometry.Point{X: p.Pos.X, Y: p.Pos.Y}, p.Pos.R)

	e.BaseStats[attributes.DEFBase] = 5*float64(p.Level) + 500

	react := &reactable.Reactable{}
	react.Init(e, core)
	react.FreezeResist = e.prof.FreezeResist
	e.Reactable = react

	e.queue = task.New(&e.timePassed)

	if core.IsDamageMode() {
		e.hpRatio = 1.0
		e.BaseStats[attributes.HPBase] = p.HP
	}

	return e
}

func (e *Enemy) Type() info.TargettableType      { return info.TargettableEnemy }
func (e *Enemy) MaxHP() float64                  { return e.Stats().Props.MaxHP() }
func (e *Enemy) HP() float64                     { return e.hpRatio * e.MaxHP() }
func (e *Enemy) SetDirection(trg geometry.Point) {}
func (e *Enemy) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}
