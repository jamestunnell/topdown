package animation

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jamestunnell/topdown/sliceutil"
)

type Controller struct {
	frameTag      string
	frames        []*Frame
	currentIndex  int
	currentOffset time.Duration
	frameOffsets  []time.Duration
	nextOffset    time.Duration
}

func NewController() *Controller {
	return &Controller{
		frameTag:      "",
		frames:        []*Frame{},
		currentIndex:  0,
		currentOffset: 0,
		frameOffsets:  []time.Duration{},
		nextOffset:    0,
	}
}

func (c *Controller) StartAnimation(frameTag string, frames []*Frame) bool {
	n := len(frames)
	if n == 0 {
		return false
	}

	offset := time.Duration(0)
	framesToOffsets := func(f *Frame) time.Duration {
		offset += f.Duration

		return offset
	}

	c.frameTag = frameTag
	c.frames = frames
	c.currentIndex = 0
	c.currentOffset = 0
	c.frameOffsets = sliceutil.Map(frames, framesToOffsets)
	c.nextOffset = c.frameOffsets[0]

	return true
}

func (c *Controller) Update(delta time.Duration) {
	c.currentOffset += delta

	if c.currentOffset >= c.nextOffset {
		if c.currentIndex == (len(c.frames) - 1) {
			c.currentOffset = c.currentOffset - c.nextOffset
			c.currentIndex = 0
			c.nextOffset = c.frameOffsets[0]
		} else {
			c.currentIndex++

			c.nextOffset = c.frameOffsets[c.currentIndex]
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
	return c.frames[c.currentIndex].Image
}
