package core

import (
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Weapon interface {
	SetIndex(int)
	Init() error
	Data() *model.WeaponData
}

type Set interface {
	SetIndex(int)
	Init() error
}

type EquipInfo struct {
	Weapon Weapon
	Sets   map[keys.Set]Set
}
