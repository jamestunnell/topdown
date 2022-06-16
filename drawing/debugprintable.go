package drawing

import "github.com/jamestunnell/topdown/debug"

type DebugPrintable interface {
	DebugData() *debug.Dataset
}
