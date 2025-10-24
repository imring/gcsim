package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Handler struct {
	targets  map[keys.Target]core.Target
	keycount int
}

func New() *Handler {
	return &Handler{
		targets:  make(map[keys.Target]core.Target),
		keycount: 0,
	}
}

func (h *Handler) Add(t core.Target) keys.Target {
	result := keys.Target(h.keycount)
	h.targets[result] = t
	h.keycount++
	return result
}

func (h *Handler) Targets() map[keys.Target]core.Target {
	return h.targets
}
