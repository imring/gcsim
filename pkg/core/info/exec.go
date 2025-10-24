package info

import "errors"

// ErrActionNotReady is returned if the requested action is not ready; this could be
// due to any of the following:
//   - Insufficient energy (burst only)
//   - Ability on cooldown
//   - Player currently in animation
var (
	// exec-specfic errors
	ErrActionNotReady        = errors.New("action is not ready yet; cannot be executed")
	ErrPlayerNotReady        = errors.New("player still in animation; cannot execute action")
	ErrInvalidAirborneAction = errors.New("player must use low_plunge or high_plunge while airborne")
	ErrActionNoOp            = errors.New("action is a noop")
	// shared character-specific errors
	ErrInvalidChargeAction = errors.New("need to use attack right before charge")
)
