package resource

//go:generate mockgen -destination=mock_resource/mocktype.go . Type

type Type interface {
	Name() string
	Load(path string) (Resource, error)
}
