package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/gcs"
)

func NewCore(seed int64, debug bool, cfg *gcs.ActionList) (*core.Core, error) {
	return core.New(core.CoreOpt{
		Seed:         seed,
		Debug:        debug,
		Delays:       cfg.Settings.Delays,
		DefHalt:      cfg.Settings.DefHalt,
		DamageMode:   cfg.Settings.DamageMode,
		EnableHitlag: cfg.Settings.EnableHitlag,
	})
}
