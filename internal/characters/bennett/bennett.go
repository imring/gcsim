package bennett

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Bennett, NewChar)
}

type char struct {
	*character.Character
	core core.Core
}

func NewChar(s core.Core, p info.CharacterProfile) (core.Character, error) {
	var err error

	c := &char{core: s}
	c.Character, err = character.New(character.Opt{
		Core:    s,
		Profile: p,
		Data:    base,

		EnergyMax:    60,
		NormalHitNum: normalHitNum,
		SkillCon:     3,
		BurstCon:     5,
	})
	if err != nil {
		return nil, err
	}

	c.SetParticleDelay(80) // special default for bennett
	return c, nil
}

func (c *char) Init() error {
	if c.GetBase().Cons >= 2 {
		// c.c2()
	}
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 7
	}
	return c.Character.AnimationStartDelay(k)
}
