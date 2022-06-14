package movecollide

import (
	"fmt"

	"github.com/jamestunnell/topdown"
	"github.com/rs/zerolog/log"
	"github.com/zergon321/cirno"
)

//go:generate mockgen -destination=mock_movecollide/mocksystem.go . System

type System interface {
	Add(id string, resource interface{})
	// Remove(id string)
	// Clear()

	Raycast(r *Ray) (*RayHit, bool)
	MoveCollide(deltaSec float64)
}

type Ray struct {
	Origin, Direction cirno.Vector
	Distance          float64
}

type RayHit struct {
	ID       string
	Position cirno.Vector
	Object   any
}

type system struct {
	space        *cirno.Space
	movables     map[string]Movable
	collidables  map[string]Collidable
	triggerables map[string]Triggerable
}

const (
	CollisionSpaceCapacity    = 1024
	CollisionSpaceSubdivision = 10

	TriggerShapeID  = 1
	ColliderShapeID = 2
)

func NewSystem(worldWidth, worldHeight float64) (System, error) {
	spaceMin := cirno.Zero()
	spaceMax := cirno.NewVector(worldWidth, worldHeight)

	log.Debug().
		Float64("w", worldWidth).
		Float64("h", worldHeight).
		Msg("creating collision space")

	space, err := cirno.NewSpace(
		CollisionSpaceSubdivision, CollisionSpaceCapacity,
		2*worldWidth, 2*worldHeight, spaceMin, spaceMax, true)
	if err != nil {
		return nil, fmt.Errorf("failed to make collision space: %w", err)
	}

	// add boundaries to the collision space
	nw := topdown.Pt[float64](0, 0)
	ne := topdown.Pt(worldWidth, 0)
	se := topdown.Pt(worldWidth, worldHeight)
	sw := topdown.Pt(0, worldHeight)

	var addLineErr error

	addLine := func(a, b topdown.Point[float64], id string) bool {
		v1 := cirno.NewVector(a.X, a.Y)
		v2 := cirno.NewVector(b.X, b.Y)
		line, err := cirno.NewLine(v1, v2)

		line.SetIdentity(ColliderShapeID)

		// collides with anything
		line.SetMask(^0)

		line.SetData(id)

		if err != nil {
			addLineErr = fmt.Errorf("failed to make space boundary: %w", err)

			return false
		}

		if err = space.Add(line); err != nil {
			addLineErr = fmt.Errorf("failed to add space boundary: %w", err)

			return false
		}

		log.Debug().Stringer("a", a).Stringer("b", b).Msg("added boundary line")

		return true
	}

	if !addLine(nw, ne, "north") || !addLine(ne, se, "east") || !addLine(se, sw, "south") || !addLine(sw, nw, "west") {
		return nil, addLineErr
	}

	s := &system{
		space:        space,
		movables:     map[string]Movable{},
		collidables:  map[string]Collidable{},
		triggerables: map[string]Triggerable{},
	}

	return s, nil
}

func (s *system) Add(id string, x interface{}) {
	if m, ok := x.(Movable); ok {
		s.movables[id] = m

		log.Debug().Str("id", id).Msg("added movable")
	}

	if c, ok := x.(Collidable); ok {
		shape := c.ColliderShape()

		shape.SetIdentity(ColliderShapeID)

		// collides with anything
		shape.SetMask(^0)

		shape.SetData(id)

		s.space.Add(shape)

		s.collidables[id] = c

		log.Debug().Str("id", id).Msg("added collidable")
	}

	if t, ok := x.(Triggerable); ok {
		shape := t.TriggerShape()

		shape.SetIdentity(TriggerShapeID)

		// don't collide with other triggers
		shape.SetMask(^TriggerShapeID)

		s.space.Add(shape)

		s.triggerables[id] = t

		log.Debug().Str("id", id).Msg("added triggerable")
	}
}

func (s *system) Raycast(r *Ray) (*RayHit, bool) {
	hitShape, hitPos, err := s.space.Raycast(r.Origin, r.Direction, r.Distance, ColliderShapeID)
	if err != nil {
		log.Warn().Err(err).Msg("raycast failed")

		return nil, false
	}

	if hitShape == nil {
		return nil, false
	}

	hitID, ok := hitShape.Data().(string)
	if !ok {
		log.Debug().Msg("ray hit, but hit shape data is not ID string")

		return nil, false
	}

	c, found := s.collidables[hitID]
	if !found {
		log.Debug().Str("collidable ID", hitID).Msg("ray hit, but collidable not found")

		return nil, false
	}

	hit := &RayHit{ID: hitID, Position: hitPos, Object: c}

	return hit, false
}

func (s *system) MoveCollide(deltaSec float64) {
	for id, m := range s.movables {
		move := m.PlanMovement(deltaSec)
		if move.Zero() {
			continue
		}

		c, found := s.collidables[id]
		if !found {
			m.Move(move)

			continue
		}

		moveDiff := cirno.NewVector(move.X, move.Y)
		shape := c.ColliderShape()

		shapes, err := s.space.WouldBeCollidedBy(shape, moveDiff, 0)
		if err != nil {
			log.Warn().Err(err).Msg("failed to figure shape move collision")

			continue
		}

		if len(shapes) > 0 {
			moveDiff = c.ResolveCollision(moveDiff, shapes)
		}

		m.Move(topdown.Vec(moveDiff.X, moveDiff.Y))

		shape.Move(moveDiff)
		s.space.AdjustShapePosition(shape)
		_, err = s.space.Update(shape)

		if err != nil {
			log.Warn().Err(err).Msg("failed to update collision space")
		}
	}
}
