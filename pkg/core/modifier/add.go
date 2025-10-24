package modifier

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (h *Handler) Add(mod info.Modifier) (bool, error) {
	config, ok := modifierMap[mod.Key]
	if !ok {
		return false, fmt.Errorf("unknown modifier: %v", mod.Key)
	}

	instance := &Instance{
		key:               mod.Key,
		source:            mod.Source,
		duration:          config.Duration,
		thinkInterval:     config.ThinkInterval,
		elementDurability: config.ElementDurability,
		hitlagAffected:    config.AffectedByHitlag,
		count:             1, // TODO: always 1?
		maxCount:          config.MaxCount,
		state:             mod.State,
		handler:           h,
	}

	for k, v := range mod.Props {
		instance.props[k] = v
	}

	if instance.duration > 0 {
		instance.decayRate = instance.elementDurability / info.Durability(instance.duration)
	}
	instance.currentThinkInterval = instance.thinkInterval

	added := false
	switch config.Stacking {
	case Refresh:
		added = h.refresh(instance)
	case Unique:
		added = h.unique(instance)
	case MultipleRefreshNoRemove:
		added = h.multipleRefreshNoRemove(instance)
	case MultipleAllRefresh:
		added = h.multipleAllRefresh(instance)
	default:
		return false, fmt.Errorf("unsupported stacking method: %v", config.Stacking)
	}

	if added {
		h.emitAdd(instance)
	}

	return true, nil
}

func (h *Handler) refresh(instance *Instance) bool {
	for _, mod := range h.modifiers {
		if mod.key == instance.key {
			mod.updateDuration(instance.duration)
			h.emitExtendDuration(mod)
			return false
		}
	}
	h.modifiers = append(h.modifiers, instance)
	return true
}

func (h *Handler) unique(instance *Instance) bool {
	for _, mod := range h.modifiers {
		if mod.key == instance.key {
			return false
		}
	}
	h.modifiers = append(h.modifiers, instance)
	return true
}

func (h *Handler) multipleRefreshNoRemove(instance *Instance) bool {
	count := 0
	var oldMod *Instance
	for _, mod := range h.modifiers {
		if mod.key == instance.key {
			count += mod.count
		}
		if oldMod == nil || oldMod.Duration() > mod.Duration() {
			oldMod = mod
		}
	}

	if count+instance.count <= instance.maxCount {
		h.modifiers = append(h.modifiers, instance)
		return true
	}

	oldMod.updateDuration(instance.duration)
	h.emitExtendDuration(oldMod)
	return false
}

func (h *Handler) multipleAllRefresh(instance *Instance) bool {
	count := 0
	for _, mod := range h.modifiers {
		if mod.key == instance.key {
			count += mod.count
			mod.updateDuration(instance.duration)
			h.emitExtendDuration(mod)
		}
	}

	added := instance.maxCount > count
	if added {
		h.modifiers = append(h.modifiers, instance)
	}
	return added
}
