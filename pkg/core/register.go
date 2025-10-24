package core

import (
	"sync"

	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

var (
	mu               sync.RWMutex
	NewCharFuncMap   = make(map[keys.Char]NewCharacterFunc)
	NewSetFuncMap    = make(map[keys.Set]NewSetFunc)
	NewWeaponFuncMap = make(map[keys.Weapon]NewWeaponFunc)
)

type NewCharacterFunc func(core Core, p info.CharacterProfile) (Character, error)
type NewSetFunc func(core Core, char Character, count int, param map[string]int) (Set, error)
type NewWeaponFunc func(core Core, char Character, p info.WeaponProfile) (Weapon, error)

func RegisterCharFunc(char keys.Char, f NewCharacterFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := NewCharFuncMap[char]; dup {
		panic("combat: RegisterChar called twice for character " + char.String())
	}
	NewCharFuncMap[char] = f
}

func RegisterSetFunc(set keys.Set, f NewSetFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := NewSetFuncMap[set]; dup {
		panic("combat: RegisterSetBonus called twice for character " + set.String())
	}
	NewSetFuncMap[set] = f
}

func RegisterWeaponFunc(weap keys.Weapon, f NewWeaponFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := NewWeaponFuncMap[weap]; dup {
		panic("combat: RegisterWeapon called twice for character " + weap.String())
	}
	NewWeaponFuncMap[weap] = f
}
