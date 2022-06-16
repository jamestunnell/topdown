package drawing

import (
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func DebugPrint(
	screen *ebiten.Image, ids []string, dps []DebugPrintable) {
	sb := strings.Builder{}

	// start with some common debug printing
	sb.WriteString("FPS: ")
	sb.WriteString(strconv.FormatFloat(ebiten.CurrentFPS(), 'f', 2, 64))
	sb.WriteRune('\n')

	mX, mY := ebiten.CursorPosition()

	sb.WriteString("mouse:\n")
	sb.WriteString("  x: ")
	sb.WriteString(strconv.Itoa(mX))
	sb.WriteRune('\n')
	sb.WriteString("  y: ")
	sb.WriteString(strconv.Itoa(mY))
	sb.WriteRune('\n')

	for i, id := range ids {
		sb.WriteString(id)
		sb.WriteString(":\n")

		ds := dps[i].DebugData()

		for _, key := range ds.SortedKeys() {
			val, found := ds.Get(key)
			if !found {
				continue
			}

			sb.WriteString("  ")
			sb.WriteString(key)
			sb.WriteString(": ")
			sb.WriteString(val)
			sb.WriteRune('\n')
		}
	}

	ebitenutil.DebugPrint(screen, sb.String())
}
