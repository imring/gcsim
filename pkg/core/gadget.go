package core

import "github.com/genshinsim/gcsim/pkg/core/info"

type Gadget interface {
	Target

	Src() int
	GadgetTyp() info.GadgetTyp
	Duration() int
	SetDuration(d int)

	SetOnKill(fn func())
	SetOnExpiry(fn func())
	SetThinkInterval(interval int)
	SetOnThinkInterval(fn func())
}
