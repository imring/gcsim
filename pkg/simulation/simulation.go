package simulation

import (
	"math/rand"
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/status"
	"github.com/genshinsim/gcsim/pkg/core/target"
	"github.com/genshinsim/gcsim/pkg/core/task"
	"github.com/genshinsim/gcsim/pkg/gcs/eval"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type Opt struct {
	Config            *info.ActionList
	Eval              *eval.Eval
	Seed              int64
	Debug             bool
	EnableHitlag      bool
	DefHalt           bool
	DamageMode        bool
	IgnoreBurstEnergy bool
	Delays            info.Delays
}

type Flags struct {
	LogDebug          bool // Used to determine logging level
	DamageMode        bool // for hp mode
	DefHalt           bool // for hitlag
	EnableHitlag      bool // hitlag enabled
	IgnoreBurstEnergy bool // for ignoring energy when using burst
	Custom            map[string]float64
}

type Core struct {
	Opt

	f     int
	flags Flags

	// basic systems
	log    glog.Logger
	rand   *rand.Rand
	events *event.System
	tasks  *task.Handler
	status *status.Handler

	// game systems
	target    *target.Handler
	player    *player.Handler
	combat    *combat.Handler
	construct *construct.Handler

	// action stuff
	noMoreActions  bool
	queue          []*info.ActionEval
	preActionDelay int

	collectors []stats.Collector
}

func New(opt Opt) (*Core, error) {
	c := &Core{}
	c.Opt = opt
	c.rand = rand.New(rand.NewSource(opt.Seed))

	c.flags.Custom = make(map[string]float64)
	c.flags.DamageMode = opt.DamageMode
	c.flags.DefHalt = opt.DefHalt
	c.flags.EnableHitlag = opt.EnableHitlag
	c.flags.IgnoreBurstEnergy = opt.IgnoreBurstEnergy

	if opt.Debug {
		c.log = glog.New(&c.f, 500)
		c.flags.LogDebug = true
	} else {
		c.log = &glog.NilLogger{}
	}

	c.events = &event.System{}
	c.tasks = task.New(&c.f)
	c.status = status.New(&c.f, c.log)

	c.target = target.New()

	c.player = player.New(player.Opt{
		F:      &c.f,
		Log:    c.log,
		Events: c.events,
		Target: c.target,
	})

	c.combat = combat.New(combat.Opt{
		F:      &c.f,
		Task:   c.tasks,
		Events: c.events,
		Player: c.player,
		Target: c.target,
		Log:    c.log,
		Rand:   c.rand,
	})

	c.construct = construct.New(&c.f, c.log, c.events)

	c.queue = make([]*info.ActionEval, 0)

	c.collectors = make([]stats.Collector, 0)
	for _, collector := range stats.Collectors() {
		enabled := c.Config.Settings.CollectStats
		if len(enabled) > 0 && !slices.Contains(enabled, collector.Name) {
			continue
		}
		stat, err := collector.New(c)
		if err != nil {
			return nil, err
		}
		c.collectors = append(c.collectors, stat)
	}

	return c, nil
}

func (c *Core) F() int                { return c.f }
func (c *Core) Events() *event.System { return c.events }
func (c *Core) Rand() *rand.Rand      { return c.rand }
func (c *Core) Tasks() task.Tasker    { return c.tasks }
func (c *Core) Log() glog.Logger      { return c.log }

func (c *Core) Tick() error {
	// things to tick:
	//	- targets
	//	- constructs
	//	- player (stamina, swap, animation, etc...)
	//		- character
	//		- shields
	//		- animation
	//		- stamina
	//		- swap
	//	- tasks
	// TODO: check for errors here?
	c.combat.Tick()
	c.construct.Tick()
	c.player.Tick()
	c.tasks.Run()
	return nil
}
