package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Modifier struct {
	Key    string
	Source keys.Target
	Props  attributes.PropMap
	State  any
}
