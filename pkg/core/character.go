package core

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Character interface {
	Target
	CharacterBase
	CharacterHP
	CharacterAction
	CharacterStam
	CharacterEnergy
	CharacterStats
	CharacterCooldown
	CharacterAnimation
	CharacterHitlag

	Init() error

	GetIndex() int
	SetIndex(index int)

	Condition(fields []string) (any, error)

	SetDirectionToClosestEnemy() // looks for closest enemy
}

type CharacterBase interface {
	UpdateBaseStats() error

	GetBase() info.CharacterBase
	GetWeapon() info.WeaponProfile
	GetTalents() info.TalentProfile
	GetSets() info.Sets

	GetEquip() EquipInfo

	Data() *model.AvatarData

	TalentLvlAttack() int
	TalentLvlSkill() int
	TalentLvlBurst() int

	SetWeapon(w Weapon)
	SetArtifactSet(key keys.Set, s Set)
}

type CharacterHP interface {
	MaxHP() float64
	CurrentHPRatio() float64
	CurrentHP() float64
	CurrentHPDebt() float64
	CurrentHPDebtRatio() float64

	SetHPByAmount(float64)
	SetHPByRatio(float64)
	ModifyHPByAmount(float64)
	ModifyHPByRatio(float64)

	ModifyHPDebtByAmount(float64)
	ModifyHPDebtByRatio(float64)

	// Heal(*info.HealInfo) (float64, float64) // return actual hp healed and amount of hp debt cleared
	// Drain(*info.DrainInfo) float64

	// ReceiveHeal(*info.HealInfo, float64) float64
}

type CharacterAction interface {
	Attack(p map[string]int) (action.Info, error)
	Aimed(p map[string]int) (action.Info, error)
	ChargeAttack(p map[string]int) (action.Info, error)
	HighPlungeAttack(p map[string]int) (action.Info, error)
	LowPlungeAttack(p map[string]int) (action.Info, error)
	Skill(p map[string]int) (action.Info, error)
	Burst(p map[string]int) (action.Info, error)

	ActionReady(a info.Action, p map[string]int) (bool, info.Failure)
	NextQueueItemIsValid(targetChar keys.Char, a info.Action, p map[string]int) error

	ResetNormalCounter()
	AdvanceNormalIndex()
	NormalCounter() int
	NextNormalCounter() int

	DashLockout() bool
	RemainingDashCD() int

	SetDashCD(cd int, lockout bool)
}

type CharacterStam interface {
	ActionStam(a info.Action, p map[string]int) float64
	AbilStamCost(a info.Action, p map[string]int) float64

	Dash(p map[string]int) (action.Info, error)
	Walk(p map[string]int) (action.Info, error)
	Jump(p map[string]int) (action.Info, error)

	ApplyDashCD()
	QueueDashStaminaConsumption(p map[string]int)

	DashLength() int
	DashToJumpLength() int
	JumpLength() int
}

type CharacterEnergy interface {
	Energy() float64
	EnergyMax() float64
	AddEnergy(src string, amt float64)
	SetParticleDelay(delay int)
}

type CharacterStats interface {
	Stats() attributes.Stats
	Stat(attr attributes.Prop) float64
	Snapshot(attackInfo *info.Attack) info.Snapshot

	StatusIsActive(status string) bool // TODO: rework?
	StatusDuration(status string) int  // TODO: rework?

	Tag(tag string) int
	SetTag(tag string, val int)

	AddModifier(info info.Modifier)
	RemoveModifier(key string) bool
}

type CharacterCooldown interface {
	SetCD(a info.Action, dur int)
	Cooldown(a info.Action) int
	SetNumCharges(a info.Action, num int)
	Charges(a info.Action) int
	ResetActionCooldown(a info.Action)
	ReduceActionCooldown(a info.Action, v int)
}

type CharacterAnimation interface {
	AnimationStartDelay(k info.AnimationDelayKey) int
}

type CharacterHitlag interface {
	ApplyHitlag(factor, dur float64)
	FramePausedOnHitlag() bool
}
