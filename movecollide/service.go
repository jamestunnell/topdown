package movecollide

type RaycastingService struct {
	MoveCollide System
}

const (
	RaycastingServiceName = "raycasting"
)

func (s *RaycastingService) Name() string {
	return RaycastingServiceName
}

func (s *RaycastingService) Raycast(ray *Ray) (*RayHit, bool) {
	return s.MoveCollide.Raycast(ray)
}
