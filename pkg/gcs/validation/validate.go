package validation

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func ValidateCharParamKeys(c keys.Char, a info.Action, keys []string) error {
	f, ok := charValidParamKeys[c]
	if !ok {
		// all is ok if no validation function registered
		return nil
	}
	return f(a, keys)
}
