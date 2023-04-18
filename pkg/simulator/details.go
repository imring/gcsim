package simulator

import (
	"github.com/genshinsim/gcsim/pkg/gcs"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

func GenerateCharacterDetails(cfg *gcs.ActionList) ([]simulation.CharacterDetail, error) {
	cpy := cfg.Copy()

	c, err := simulation.NewCore(CryptoRandSeed(), false, cpy)
	if err != nil {
		return nil, err
	}
	//create a new simulation and run
	sim, err := simulation.New(cpy, c)
	if err != nil {
		return nil, err
	}

	return sim.CharacterDetails(), nil
}
