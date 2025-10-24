package target

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/geometry"
)

type SingleTarget struct {
	Target keys.Target
}

func (s *SingleTarget) PointInShape(p geometry.Point) bool            { return true }
func (s *SingleTarget) IntersectCircle(in geometry.Circle) bool       { return false }
func (s *SingleTarget) IntersectRectangle(in geometry.Rectangle) bool { return false }
func (s *SingleTarget) Pos() geometry.Point                           { return geometry.Point{X: 0, Y: 0} }
func (s *SingleTarget) String() string                                { return fmt.Sprintf("single target: %v", s.Target) }
