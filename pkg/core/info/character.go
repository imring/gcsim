package info

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type CharacterProfile struct {
	Base           CharacterBase               `json:"base"`
	Weapon         WeaponProfile               `json:"weapon"`
	Talents        TalentProfile               `json:"talents"`
	Stats          []float64                   `json:"stats"`
	StatsByLabel   map[string][]float64        `json:"stats_by_label"`
	RandomSubstats *RandomSubstats             `json:"random_substats"`
	Sets           Sets                        `json:"sets"`
	SetParams      map[keys.Set]map[string]int `json:"-"`
	Params         map[string]int              `json:"-"`
}

type RandomSubstats struct {
	Rarity  int `json:"rarity"`
	Sand    attributes.Prop
	Goblet  attributes.Prop
	Circlet attributes.Prop
}

func (r RandomSubstats) Validate() error {
	//TODO: support more than just 5 stars
	if r.Rarity != 5 {
		return fmt.Errorf("unsupported rarity: %v", r.Rarity)
	}
	if r.Sand == attributes.NoProp {
		return errors.New("sand main stat not specified")
	}
	if r.Goblet == attributes.NoProp {
		return errors.New("goblet main stat not specified")
	}
	if r.Circlet == attributes.NoProp {
		return errors.New("circlet main stat not specified")
	}
	// main stat have to be valid
	switch r.Sand {
	case attributes.HPPercent:
	case attributes.ATKPercent:
	case attributes.DEFPercent:
	case attributes.EM:
	case attributes.ER:
	default:
		return fmt.Errorf("%v is not a valid main stat for sand", r.Sand.String())
	}

	switch r.Goblet {
	case attributes.HPPercent:
	case attributes.ATKPercent:
	case attributes.DEFPercent:
	case attributes.EM:
	case attributes.PyroP:
	case attributes.HydroP:
	case attributes.CryoP:
	case attributes.ElectroP:
	case attributes.AnemoP:
	case attributes.GeoP:
	case attributes.DendroP:
	case attributes.PhyP:
	default:
		return fmt.Errorf("%v is not a valid main stat for sand", r.Sand.String())
	}

	switch r.Circlet {
	case attributes.HPPercent:
	case attributes.ATKPercent:
	case attributes.DEFPercent:
	case attributes.EM:
	case attributes.CR:
	case attributes.CD:
	case attributes.HealBonus:
	default:
		return fmt.Errorf("%v is not a valid main stat for sand", r.Sand.String())
	}

	return nil
}

func (c *CharacterProfile) Clone() CharacterProfile {
	r := *c
	r.Weapon.Params = make(map[string]int)
	for k, v := range c.Weapon.Params {
		r.Weapon.Params[k] = v
	}
	r.Stats = make([]float64, len(c.Stats))
	copy(r.Stats, c.Stats)
	r.Sets = make(map[keys.Set]int)
	for k, v := range c.Sets {
		r.Sets[k] = v
	}

	return r
}

type CharacterBase struct {
	Key       keys.Char              `json:"key"`
	Rarity    int                    `json:"rarity"`
	Element   attributes.ElementType `json:"element"`
	Level     int                    `json:"level"`
	MaxLevel  int                    `json:"max_level"`
	Ascension int                    `json:"ascension"`
	Cons      int                    `json:"cons"`
}

type BodyType int

const (
	BodyBoy BodyType = iota
	BodyGirl
	BodyMale
	BodyLady
	BodyLoli
)

type ZoneType int

const (
	ZoneUnknown ZoneType = iota
	ZoneMondstadt
	ZoneLiyue
	ZoneInazuma
	ZoneSumeru
	ZoneFontaine
	ZoneNatlan
	ZoneSnezhnaya
)
