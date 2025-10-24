package core

import "github.com/genshinsim/gcsim/pkg/core/info"

// ActionEvaluator provides method for getting next action
type ActionEvaluator interface {
	NextAction() (*info.ActionEval, error) // NextAction should reuturn the next action, or nil if no actions left
	Continue()
	Exit() error
	Err() error
	Start()
}
