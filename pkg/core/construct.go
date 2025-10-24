package core

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

type Construct interface {
	OnDestruct()
	Key() int
	Type() info.GeoConstructType
	Expiry() int
	IsLimited() bool
	Count() int
	Direction() geometry.Point
	Pos() geometry.Point
}
