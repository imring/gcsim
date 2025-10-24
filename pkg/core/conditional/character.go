package conditional

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/shortcut"
)

func evalCharacter(c core.Core, key keys.Char, fields []string) (any, error) {
	if err := fieldsCheck(fields, 2, "character"); err != nil {
		return 0, err
	}
	char := c.GetCharacterByKey(key)
	if char == nil {
		return 0, fmt.Errorf("character %v not in team when evaluating condition", key)
	}

	// special case for ability conditions. since fields are swapped
	// .kokomi.<abil>.<cond>
	typ := fields[1]
	act := info.StringToAction(typ)
	if act != info.InvalidAction {
		if err := fieldsCheck(fields, 3, "character ability"); err != nil {
			return 0, err
		}
		return evalCharacterAbil(c, char, act, fields[2])
	}

	charCat := "character " + typ

	switch typ {
	case "id":
		return int(char.GetBase().Key), nil
	case "cons":
		return char.GetBase().Cons, nil
	case "energy":
		return char.Energy(), nil
	case "energymax":
		return char.EnergyMax(), nil
	case "hp":
		return char.CurrentHP(), nil
	case "hpmax":
		return char.MaxHP(), nil
	case "hpratio":
		return char.CurrentHPRatio(), nil
	case "normal":
		return char.NextNormalCounter(), nil
	case "onfield":
		return c.ActiveCharacter() == char.GetIndex(), nil
	case "weapon":
		return int(char.GetWeapon().Key), nil
	case "status":
		if err := fieldsCheck(fields, 3, charCat); err != nil {
			return 0, err
		}
		return char.StatusDuration(fields[2]), nil
	case "mods":
		if err := fieldsCheck(fields, 3, charCat); err != nil {
			return 0, err
		}
		return char.StatusDuration(fields[2]), nil
	// TODO: infusion
	// case "infusion":
	// 	if err := fieldsCheck(fields, 3, charCat); err != nil {
	// 		return 0, err
	// 	}
	// 	return c.Player.WeaponInfuseIsActive(char.Index, fields[2]), nil
	case "tags":
		if err := fieldsCheck(fields, 3, charCat); err != nil {
			return 0, err
		}
		return char.Tag(fields[2]), nil
	case "stats":
		if err := fieldsCheck(fields, 3, charCat); err != nil {
			return 0, err
		}
		return evalCharacterStats(char, fields[2])
	case "bol":
		return char.CurrentHPDebt(), nil
	case "bolratio":
		return char.CurrentHPDebtRatio(), nil
	case "sets":
		if err := fieldsCheck(fields, 3, charCat); err != nil {
			return 0, err
		}
		return evalCharacterSets(char, fields[2])
	default: // .kokomi.*
		return char.Condition(fields[1:])
	}
}

func evalCharacterStats(char core.Character, stat string) (float64, error) {
	key := attributes.StrToPropType(stat)
	if key == -1 {
		return 0, fmt.Errorf("invalid stat key %v in character stat condition", stat)
	}
	return char.Stat(key), nil
}

func evalCharacterSets(char core.Character, set string) (int, error) {
	setKey, ok := shortcut.SetNameToKey[set]
	if !ok {
		return 0, fmt.Errorf("invalid set key %v in character set condition", set)
	}
	setInfo, ok := char.GetSets()[setKey]
	if !ok {
		return 0, nil
	}
	return setInfo, nil
}

func evalCharacterAbil(c core.Core, char core.Character, act info.Action, typ string) (any, error) {
	switch typ {
	case "cd":
		if act == info.ActionSwap {
			return c.PlayerSwapCD(), nil
		}
		if act == info.ActionDash {
			if c.ActiveCharacter() == char.GetIndex() && c.PlayerDashLockout() {
				return c.PlayerRemainingDashCD(), nil
			}
			if c.ActiveCharacter() != char.GetIndex() && char.DashLockout() {
				return char.RemainingDashCD(), nil
			}
			return 0, nil
		}
		return char.Cooldown(act), nil
	case "charge":
		return char.Charges(act), nil
	case "ready":
		if act == info.ActionSwap {
			return c.PlayerSwapCD() == 0 || c.ActiveCharacter() == char.GetIndex(), nil
		}
		// TODO: nil map may cause problems here??
		ok, _ := char.ActionReady(act, nil)
		return ok, nil
	default:
		return 0, fmt.Errorf("bad character ability condition: invalid type %v", typ)
	}
}
