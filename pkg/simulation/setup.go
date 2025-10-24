package simulation

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/target/enemy"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

func (c *Core) addChar(p geometry.Point, r float64, profile info.CharacterProfile) (int, error) {
	var err error

	// initialize character
	f, ok := core.NewCharFuncMap[profile.Base.Key]
	if !ok {
		return -1, fmt.Errorf("invalid character: %v", profile.Base.Key.String())
	}
	char, err := f(c, profile)
	if err != nil {
		return -1, err
	}
	index := c.player.AddChar(char)
	char.SetShape(*geometry.NewCircle(p, r, geometry.DefaultDirection(), 360)) // TODO: ?

	// initialize weapon
	wf, ok := core.NewWeaponFuncMap[profile.Weapon.Key]
	if !ok {
		return -1, fmt.Errorf("unrecognized weapon %v for character %v", profile.Weapon.Key, profile.Base.Key.String())
	}
	weap, err := wf(c, char, profile.Weapon)
	if err != nil {
		return -1, err
	}
	char.SetWeapon(weap)

	// set bonus
	total := 0
	for key, count := range profile.Sets {
		total += count
		af, ok := core.NewSetFuncMap[key]
		if ok {
			s, err := af(c, char, count, profile.SetParams[key])
			if err != nil {
				return -1, err
			}
			char.SetArtifactSet(key, s)
		} else {
			return -1, fmt.Errorf("character %v has unrecognized artifact: %v", profile.Base.Key.String(), key)
		}
	}
	// TODO: this should be handled by parser
	if total > 5 {
		return -1, fmt.Errorf("total set count cannot exceed 5, got %v", total)
	}

	return index, nil
}

func (c *Core) setupCharacters(p geometry.Point, r float64, chars []info.CharacterProfile, initial keys.Char) error {
	if len(chars) > 4 {
		return errors.New("cannot have more than 4 characters per team")
	}
	dup := make(map[keys.Char]bool)

	active := -1
	for i := range chars {
		// if using random stats, ignore all stats except main
		if chars[i].RandomSubstats != nil {
			stats, err := generateRandSubs(chars[i].RandomSubstats, c.rand)
			if err != nil {
				return err
			}
			chars[i].Stats = stats
			clear(chars[i].StatsByLabel)
		}

		i, err := c.addChar(p, r, chars[i])
		if err != nil {
			return err
		}

		if chars[i].Base.Key == initial {
			c.player.SetActive(i)
			active = i
		}

		if _, ok := dup[chars[i].Base.Key]; ok {
			return fmt.Errorf("duplicated character %v", chars[i].Base.Key)
		}
		dup[chars[i].Base.Key] = true
	}

	if active == -1 {
		return errors.New("no active character set")
	}
	return nil
}

func (c *Core) setupTargets(targets []info.EnemyProfile) error {
	// add targets
	for i := range targets {
		v := &targets[i]
		if v.Pos.R == 0 {
			return fmt.Errorf("target cannot have 0 radius (index %v): %v", i, v)
		}
		e := enemy.New(c, *v)
		c.combat.AddEnemy(e)
		// s.stats.ElementUptime[i+1] = make(map[core.EleType]int)
	}

	c.player.SetPlayerDirectionToClosestEnemy()

	return nil
}
