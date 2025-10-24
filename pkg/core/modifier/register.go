package modifier

import (
	"sync"

	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	mu          sync.Mutex
	modifierMap = make(map[string]Config)
)

func Register(key string, config Config) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := modifierMap[key]; dup {
		panic("modifier.Register called twice for mod " + key)
	}
	modifierMap[key] = config
}

type Config struct {
	Stacking StackingBehavior
	MaxCount int

	Duration          int
	ThinkInterval     int
	ElementDurability info.Durability
	AffectedByHitlag  bool

	// TODO:
	// isUnique
	// isLimitedProperties

	Listeners Listeners
}

type StackingBehavior int

const (
	Refresh StackingBehavior = iota
	Unique
	Prolong
	RefreshAndAddDurability
	Multiple
	MultipleRefresh
	MultipleRefreshNoRemove
	MultipleAllRefresh
	GlobalUnique
	Overlap
	RefreshUniqueDurability
	OverlapRefreshDuration
)
