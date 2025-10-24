package core

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type Reactable interface {
	Tick()

	React(a *info.AttackEvent)
	AttachOrRefill(a *info.AttackEvent) bool
	SetAuraDurability(mod attributes.ElementType, dur info.Durability, decay info.Durability)

	ActiveAuraString() []string
	AuraCount() int
	GetAuraDurability(mod attributes.ElementType) info.Durability
	GetDurability() []info.Durability
	GetAuraDecayRate(mod attributes.ElementType) info.Durability
	AuraContains(e ...attributes.ElementType) bool

	ReactableBloom
	ReactableBurning
	ReactableCatalyze
	ReactableCrystallize
	ReactableEC
	ReactableFreeze
	ReactableMelt
	ReactableOverload
	ReactableSuperconduct
	ReactableSwirl
	ReactableVaporize
}

type ReactableBloom interface {
	TryBloom(a *info.AttackEvent) bool
}

type ReactableBurning interface {
	TryBurning(a *info.AttackEvent) bool
	IsBurning() bool
}

type ReactableCatalyze interface {
	TryAggravate(a *info.AttackEvent) bool
	TrySpread(a *info.AttackEvent) bool
	TryQuicken(a *info.AttackEvent) bool
}

type ReactableCrystallize interface {
	TryCrystallizeElectro(a *info.AttackEvent) bool
	TryCrystallizeHydro(a *info.AttackEvent) bool
	TryCrystallizeCryo(a *info.AttackEvent) bool
	TryCrystallizePyro(a *info.AttackEvent) bool
	TryCrystallizeFrozen(a *info.AttackEvent) bool
}

type ReactableEC interface {
	TryAddEC(a *info.AttackEvent) bool
}

type ReactableFreeze interface {
	TryFreeze(a *info.AttackEvent) bool
	PoiseDMGCheck(a *info.AttackEvent) bool
	ShatterCheck(a *info.AttackEvent) bool
	SetFreezeResist(resist float64)
}

type ReactableMelt interface {
	TryMelt(a *info.AttackEvent) bool
}

type ReactableOverload interface {
	TryOverload(a *info.AttackEvent) bool
}

type ReactableSuperconduct interface {
	TrySuperconduct(a *info.AttackEvent) bool
	TryFrozenSuperconduct(a *info.AttackEvent) bool
}

type ReactableSwirl interface {
	TrySwirlElectro(a *info.AttackEvent) bool
	TrySwirlHydro(a *info.AttackEvent) bool
	TrySwirlCryo(a *info.AttackEvent) bool
	TrySwirlPyro(a *info.AttackEvent) bool
	TrySwirlFrozen(a *info.AttackEvent) bool
}

type ReactableVaporize interface {
	TryVaporize(a *info.AttackEvent) bool
}
