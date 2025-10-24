package player

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func (h *Handler) PlayerTarget() core.Character {
	return h.GetCharacter(h.ActiveCharacter())
}

func (h *Handler) GetCharacter(index int) core.Character {
	return h.chars[index]
}

func (h *Handler) GetCharacterByKey(key keys.Char) core.Character {
	if idx, ok := h.charPos[key]; ok {
		return h.chars[idx]
	}
	return nil
}

func (h *Handler) GetCharacterByTarget(key keys.Target) core.Character {
	for _, c := range h.chars {
		if c.Key() == key {
			return c
		}
	}
	return nil
}
