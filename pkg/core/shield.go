package core

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type Shield interface {
	ShieldOwner() int
	ShieldTarget() int // -1 to apply to all characters, character index otherwise
	Key() int
	Type() info.ShieldType
	ShieldStrength(ele attributes.ElementType, bonus float64) float64
	OnDamage(dmg float64, ele attributes.ElementType, bonus float64) (float64, bool) // return dmg taken and shield stays
	OnExpire()
	OnOverwrite()
	Expiry() int
	SetExpiry(expiry int)
	CurrentHP() float64
	Element() attributes.ElementType
	Desc() string
}
