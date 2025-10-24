package player

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (h *Handler) SetAirborne(src info.AirborneSource) error {
	if src < info.AirborneGrounded || src >= info.TerminateAirborne {
		// do nothing
		return fmt.Errorf("invalid airborne source: %v", src)
	}
	h.airborne = src
	return nil
}

func (h *Handler) Airborne() info.AirborneSource {
	return h.airborne
}
