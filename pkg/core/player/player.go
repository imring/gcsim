package player

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/animation"
	"github.com/genshinsim/gcsim/pkg/core/target"
	"github.com/genshinsim/gcsim/pkg/core/task"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

const (
	MaxStam      = 240
	StamCDFrames = 90
	SwapCDFrames = 60
)

type Handler struct {
	Opt

	// handlers
	Animation *animation.Handler
	Shield    *ShieldHandler

	// tracking
	chars   []core.Character
	active  int
	charPos map[keys.Char]int

	// stam
	stam        float64
	lastStamUse int

	// airborne source
	airborne info.AirborneSource

	// swap
	swapCD  int
	swapICD int

	// dash: dash fails iff lockout && on CD
	dashCDExpirationFrame int
	dashLockout           bool

	// last action
	lastAction info.LastAction
}

type Opt struct {
	F      *int
	Log    glog.Logger
	Events *event.System
	Target *target.Handler
	Tasks  task.Tasker
	Delays info.Delays
}

func New(opts Opt) *Handler {
	h := &Handler{
		Opt:     opts,
		chars:   make([]core.Character, 0, 4),
		charPos: make(map[keys.Char]int),
		stam:    MaxStam,
		swapICD: SwapCDFrames,
	}
	h.Animation = animation.New(opts.F, opts.Log, opts.Events, opts.Tasks)
	h.Shield = NewShieldHandler(opts.F, opts.Log, opts.Events, h)
	return h
}

func (h *Handler) ByIndex(i int) core.Character { return h.chars[i] }
func (h *Handler) ActiveCharacter() int         { return h.active }
func (h *Handler) Chars() []core.Character      { return h.chars }
func (h *Handler) CharCount() int               { return len(h.chars) }

func (h *Handler) AddChar(char core.Character) int {
	h.chars = append(h.chars, char)
	index := len(h.chars) - 1
	char.SetIndex(index)
	h.charPos[char.GetBase().Key] = index

	return index
}

func (h *Handler) SetPlayerPos(pos geometry.Point) {
	for _, c := range h.chars {
		c.SetPos(pos)
	}
}

func (h *Handler) SetPlayerDirectionToClosestEnemy() {
	for _, c := range h.chars {
		c.SetDirectionToClosestEnemy()
	}
}

func (h *Handler) SwapCD() int          { return h.swapCD }
func (h *Handler) SetSwapCD(cd int)     { h.swapCD = cd }
func (h *Handler) SwapICD() int         { return h.swapICD }
func (h *Handler) SetSwapICD(cd int)    { h.swapICD = cd }
func (h *Handler) DashLockout() bool    { return h.dashLockout }
func (h *Handler) RemainingDashCD() int { return h.dashCDExpirationFrame - *h.F }

func (h *Handler) SetDashCD(b bool, cd int) {
	h.dashLockout = b
	h.dashCDExpirationFrame = *h.F + cd
}

func (h *Handler) Stam() float64 { return h.stam }
func (h *Handler) UseStam(amount float64, a info.Action) {
	h.stam -= amount
	// this really shouldn't happen??
	if h.stam < 0 {
		h.stam = 0
	}
	h.lastStamUse = *h.F
	h.Events.StamUse.Emit(event.StamUseEvent{
		Action: a,
	})
}

func (h *Handler) ResetAllNormalCounter() {
	for _, char := range h.chars {
		char.ResetNormalCounter()
	}
}

func (h *Handler) swap(to keys.Char) func() {
	return func() {
		prev := h.active
		h.active = h.charPos[to]

		// still have remaining frames left on dash CD, save in char for when they go on-field again
		if h.dashCDExpirationFrame > *h.F {
			h.chars[prev].SetDashCD(h.dashCDExpirationFrame-*h.F, h.dashLockout)
		}

		// set the new DashCDExpirationFrame and reset character remaining back to 0
		h.dashCDExpirationFrame = *h.F + h.chars[h.active].RemainingDashCD()
		h.dashLockout = h.chars[h.active].DashLockout()
		h.chars[h.active].SetDashCD(0, h.dashLockout)

		h.swapCD = h.swapICD
		h.ResetAllNormalCounter()

		evt := h.Log.NewEvent("executed swap", glog.LogActionEvent, h.active).
			Write("action", "swap").
			Write("target", to.String())

		if h.chars[prev].RemainingDashCD() > 0 {
			evt.Write("prev_dash_cd", h.chars[prev].RemainingDashCD()).
				Write("prev_dash_lockout", h.chars[prev].DashLockout())
		}

		if h.dashCDExpirationFrame > *h.F {
			evt.Write("target_dash_cd", h.dashCDExpirationFrame-*h.F).
				Write("target_dash_expiry_frame", h.dashCDExpirationFrame).
				Write("target_dash_lockout", h.dashLockout)
		}

		h.Events.CharacterSwap.Emit(event.CharacterSwapEvent{
			Prev: prev,
			Next: h.active,
		})
	}
}

func (h *Handler) LastAction() info.LastAction { return h.lastAction }

func (h *Handler) SetActive(i int) {
	if i < 0 || i >= len(h.chars) {
		return
	}
	h.active = i
}

func (h *Handler) Tick() {
	//	- player (stamina, swap, animation, etc...)
	//		- character
	//		- shields
	//		- animation
	//		- stamina
	//		- swap
	// recover stamina
	if h.stam < MaxStam && *h.F-h.lastStamUse > StamCDFrames {
		h.stam += 25.0 / 60
		if h.stam > MaxStam {
			h.stam = MaxStam
		}
	}
	if h.swapCD > 0 {
		h.swapCD--
	}
	h.Shield.Tick()
	h.Animation.Tick()
	for _, c := range h.chars {
		c.Tick()
	}
}

// InitializeTeam will set up resonance event hooks and calculate
// all character base stats
func (h *Handler) InitializeTeam() error {
	var err error
	for _, char := range h.chars {
		err = char.UpdateBaseStats()
		if err != nil {
			return err
		}
	}
	// loop again to initialize
	for i, char := range h.chars {
		err = char.Init()
		if err != nil {
			return err
		}
		char.GetEquip().Weapon.Init()
		for _, v := range char.GetEquip().Sets {
			v.Init()
		}

		// TODO
		// set each char's starting hp ratio
		// switch {
		// case char.StartHP > 0 && char.StartHPRatio > 0:
		// 	char.SetHPByRatio(float64(char.StartHPRatio) / 100.0)
		// 	char.ModifyHPByAmount(float64(char.StartHP))
		// case char.StartHP > 0:
		// 	char.SetHPByAmount(float64(char.StartHP))
		// case char.StartHPRatio > 0:
		// 	char.SetHPByRatio(float64(char.StartHPRatio) / 100.0)
		// default:
		// 	char.SetHPByRatio(1)
		// }

		h.Log.NewEvent("starting hp set", glog.LogCharacterEvent, i).
			Write("starting_hp_ratio", char.CurrentHPRatio()).
			Write("starting_hp", char.CurrentHP())
	}
	return nil
}
