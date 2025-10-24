package modifier

func (h *Handler) Remove(key string) bool {
	i := 0
	var removedMods []*Instance
	for _, mod := range h.modifiers {
		if mod.key == key {
			removedMods = append(removedMods, mod)
		} else {
			h.modifiers[i] = mod
			i++
		}
	}
	h.modifiers = h.modifiers[:i]
	h.emitRemove(removedMods)
	return len(removedMods) > 0
}
