package core

import (
	"math/rand"

	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/task"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

type Core interface {
	F() int
	Events() *event.System
	Rand() *rand.Rand
	Tasks() task.Tasker
	Log() glog.Logger

	Combat
	Constructs
	Player
	Settings
	Status
}

type Combat interface {
	QueueAttackWithSnap(snap info.Snapshot, qa info.QueueAttack)
	QueueAttackEvent(ae *info.AttackEvent, dmgDelay int)
	QueueAttack(qa info.QueueAttack)

	GetEnemy(i int) Enemy
	GetEnemies() []Enemy
	SetEnemyPos(i int, pos geometry.Point)
	KillEnemy(i int)

	AddGadget(g Gadget)
	RemoveGadget(key keys.Target)
	GetGadgets() []Gadget

	SetDefaultTarget(key keys.Target)
	ClosestEnemy(pos geometry.Point) Enemy
	ClosestGadget(pos geometry.Point) Gadget
	ClosestEnemyWithinArea(ap info.AttackPattern, filter func(t Enemy) bool) Enemy
}

type Constructs interface {
	ConstructCountByType(t info.GeoConstructType) int
	ConstructExpiry(t info.GeoConstructType) int
}

type Player interface {
	PlayerTarget() Character

	ActiveCharacter() int
	GetCharacter(index int) Character
	GetCharacterByKey(key keys.Char) Character
	GetCharacterByTarget(key keys.Target) Character

	SetPlayerPos(pos geometry.Point)
	SetPlayerDirectionToClosestEnemy()

	AddPlayerShield(shd Shield)

	PlayerSwapCD() int
	SetPlayerSwapICD(cd int)
	PlayerDashLockout() bool
	PlayerRemainingDashCD() int // c.Player.DashCDExpirationFrame - c.F
	SetPlayerDashCD(lockout bool, cd int)
	PlayerUseStam(amount float64, a info.Action)
	PlayerStam() float64

	PlayerCurrentState() info.AnimationState
	GetPlayerLastAction() info.LastAction
	PlayerAirborne() info.AirborneSource
	SetPlayerAirborne(a info.AirborneSource)
}

type Settings interface {
	IsDamageMode() bool
	IsHitlagEnabled() bool
	CanBeDefenseHalt() bool
	IgnoreBurstEnergy() bool
}

type Status interface {
	StatusDuration(status string) int
}
