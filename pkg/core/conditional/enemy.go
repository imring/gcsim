package conditional

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/reactable"
)

func evalDebuff(c core.Core, fields []string) (bool, error) {
	// .debuff.res.t1.name
	if err := fieldsCheck(fields, 4, "debuff"); err != nil {
		return false, err
	}
	typ := fields[1]
	trg := fields[2]
	mod := fields[3]

	e, err := parseTarget(c, trg)
	if err != nil {
		return false, fmt.Errorf("bad debuff condition: %w", err)
	}

	switch typ {
	case "def":
		return e.HasModifier(mod), nil
	case "res":
		return e.HasModifier(mod), nil
	default:
		return false, fmt.Errorf("bad debuff condition: invalid type %s", typ)
	}
}

func evalElement(c core.Core, fields []string) (float64, error) {
	// .element.t1.pyro
	if err := fieldsCheck(fields, 3, "element"); err != nil {
		return 0, err
	}
	trg := fields[1]
	ele := fields[2]

	e, err := parseTarget(c, trg)
	if err != nil {
		return 0, fmt.Errorf("bad element condition: %w", err)
	}

	if ele == "burning" {
		result := info.Durability(0)
		if e.IsBurning() {
			result = e.GetAuraDurability(attributes.ElementBurning)
		}
		return float64(result), nil
	}

	elekey := attributes.StringToEle(ele)
	if elekey == attributes.ElementNone {
		return 0, fmt.Errorf("bad element condition: invalid element %s", ele)
	}
	result := info.Durability(0)
	for i := attributes.ElementNone; i < attributes.ElementCOUNT; i++ {
		dur := e.GetAuraDurability(i)
		if dur > reactable.ZeroDur && dur > result {
			result = dur
		}
	}
	return float64(result), nil
}

func parseTarget(c core.Core, trg string) (core.Enemy, error) {
	trg = strings.TrimPrefix(trg, "t")
	tid, err := strconv.ParseInt(trg, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid target %v", trg)
	}

	t := c.GetEnemy(int(tid))
	if t == nil {
		return nil, fmt.Errorf("invalid target %v", tid)
	}
	return t, nil
}
