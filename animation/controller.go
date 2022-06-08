package animation

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Controller struct {
	frameTag      string
	frameImages   Images
	frameDur      time.Duration
	currentIndex  int
	currentOffset time.Duration
}

func NewController() *Controller {
	return &Controller{
		frameTag:      "",
		frameImages:   Images{},
		frameDur:      0,
		currentIndex:  0,
		currentOffset: 0,
	}
}

func (c *Controller) StartAnimation(tag string, images Images, frameDur time.Duration) bool {
	n := len(images)
	if n == 0 {
		return false
	}

	c.frameTag = tag
	c.frameImages = images
	c.frameDur = frameDur
	c.currentIndex = 0
	c.currentOffset = 0

	return true
}

func (c *Controller) Update(delta time.Duration) {
	c.currentOffset += delta

	for c.currentOffset >= c.frameDur {
		c.currentOffset -= c.frameDur

		if c.currentIndex == (len(c.frameImages) - 1) {
			c.currentIndex = 0
		} else {
			c.currentIndex++
		}
	}
}

func (c *Controller) CurrentFrameTag() string {
	return c.frameTag
}

func (c *Controller) CurrentFrameIndex() int {
	return c.currentIndex
}

func (c *Controller) CurrentFrameImage() *ebiten.Image {
	return c.frameImages[c.currentIndex]
}
