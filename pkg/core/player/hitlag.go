package player

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (h *Handler) ApplyHitlag(char int, factor, dur float64) {
	// make sure we only apply hitlag if this character is on field
	if char != h.active {
		return
	}

	h.chars[char].ApplyHitlag(factor, dur)

	// also extend infusion
	// TODO: this is a really awkward place to apply this
	// h.ExtendInfusion(char, factor, dur)

	// extend the dash cd by the hitlag extension amount
	if h.dashCDExpirationFrame > *h.F {
		ext := int(math.Ceil(dur * (1 - factor)))
		h.dashCDExpirationFrame += ext

		var evt glog.Event
		if h.dashLockout {
			evt = h.Log.NewEvent("dash cd hitlag extended", glog.LogHitlagEvent, char)
		} else {
			evt = h.Log.NewEvent("dash lockout evaluation hitlag extended", glog.LogHitlagEvent, char)
		}
		evt.Write("extension", ext).
			Write("expiry", h.dashCDExpirationFrame-*h.F).
			Write("expiry_frame", h.dashCDExpirationFrame).
			Write("lockout", h.dashLockout)
	}
}
