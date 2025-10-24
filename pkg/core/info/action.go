package info

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Action int

const (
	InvalidAction Action = iota
	ActionSkill
	ActionBurst
	ActionAttack
	ActionCharge
	ActionHighPlunge
	ActionLowPlunge
	ActionAim
	ActionDash
	ActionJump
	// following action have to implementations
	ActionSwap
	ActionWalk
	ActionWait  // character should stand around and wait
	ActionDelay // delay before executing next action
	EndActionType
	// these are only used for frames purposes and that's why it's after end
	ActionSkillHoldFramesOnly
)

var astr = []string{
	"invalid",
	"skill",
	"burst",
	"attack",
	"charge",
	"high_plunge",
	"low_plunge",
	"aim",
	"dash",
	"jump",
	"swap",
	"walk",
	"wait",
	"delay",
}

func (a Action) String() string {
	return astr[a]
}

func (a Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(astr[a])
}

func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range astr {
		if v == s {
			*a = Action(i)
			return nil
		}
	}
	return errors.New("unrecognized action")
}

func StringToAction(s string) Action {
	for i, v := range astr {
		if v == s {
			return Action(i)
		}
	}
	return InvalidAction
}

type AnimationState int

const (
	AnimationStateIdle AnimationState = iota
	AnimationStateNormalAttack
	AnimationStateChargeAttack
	AnimationStatePlungeAttack
	AnimationStateSkill
	AnimationStateBurst
	AnimationStateAim
	AnimationStateDash
	AnimationStateJump
	AnimationStateWalk
	AnimationStateSwap
)

var statestr = []string{
	"idle",
	"normal",
	"charge",
	"plunge",
	"skill",
	"burst",
	"aim",
	"dash",
	"jump",
	"walk",
	"swap",
}

func (a AnimationState) String() string {
	return statestr[a]
}

// ActionEval represents a sim action
type ActionEval struct {
	Char   keys.Char
	Action Action
	Param  map[string]int
}

type LastAction struct {
	Type  Action
	Param map[string]int
	Char  int
}
