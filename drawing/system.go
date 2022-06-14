package drawing

import "github.com/hajimehoshi/ebiten/v2"

type System interface {
	Add(id string, r interface{})
	Remove(id string) bool
	Clear()

	Surface() *ebiten.Image
}
