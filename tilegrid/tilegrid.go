package tilegrid

import (
	"fmt"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/drawing"
	"github.com/jamestunnell/topdown/mathutil"
	"github.com/jamestunnell/topdown/resource"
	"github.com/jamestunnell/topdown/sliceutil"
	"github.com/jamestunnell/topdown/sprite"
)

type TileGrid struct {
	MinPosition topdown.Vector     `json:"minPosition"`
	TileSize    topdown.Size[int]  `json:"tileSize"`
	TileDefs    map[string]TileDef `json:"tileDefs"`
	TileRows    []string           `json:"tileRows"`

	worldArea topdown.Rectangle[float64]
	center    topdown.Point[float64]
	rows      []*Row
	nRows     int
	nCols     int
}

type TileDef struct {
	SpriteSetRef string             `json:"spriteSetRef"`
	StartPoint   topdown.Point[int] `json:"startPoint"`
}

const (
	RefIDSeparator = " "
)

func New(tileSize topdown.Size[int]) *TileGrid {
	tg := &TileGrid{
		MinPosition: topdown.Vec(0, 0),
		TileSize:    tileSize,
		TileDefs:    map[string]TileDef{},
		TileRows:    []string{},
		rows:        []*Row{},
	}

	return tg
}

func (tg *TileGrid) Initialize(mgr resource.Manager) error {
	spriteSetRefs := mapset.NewSet[string]()

	for _, tileDef := range tg.TileDefs {
		spriteSetRefs.Add(tileDef.SpriteSetRef)
	}

	spriteSets := map[string]*sprite.SpriteSet{}

	for _, ref := range spriteSetRefs.ToSlice() {
		sprites, err := resource.GetAs[*sprite.SpriteSet](mgr, ref)
		if err != nil {
			return fmt.Errorf("failed to get sprites: %w", err)
		}

		spriteSets[ref] = sprites
	}

	if err := tg.MakeRows(spriteSets); err != nil {
		return fmt.Errorf("failed to make rows: %w", err)
	}

	dX := float64(tg.nCols * tg.TileSize.Width)
	dY := float64(tg.nRows * tg.TileSize.Height)
	tg.worldArea = topdown.Rectangle[float64]{
		Min: topdown.Pt(tg.MinPosition.X, tg.MinPosition.Y),
		Max: topdown.Pt(tg.MinPosition.X+dX, tg.MinPosition.Y+dY),
	}

	tg.center = tg.worldArea.Min.Add(tg.worldArea.Size().Center())

	return nil
}

func (tg *TileGrid) Center() topdown.Point[float64] {
	return tg.worldArea.Size().Center()
}

func (tg *TileGrid) MakeRows(spriteSets map[string]*sprite.SpriteSet) error {
	images := map[string]*ebiten.Image{}
	xScales := map[string]float64{}
	yScales := map[string]float64{}

	for defID, tileDef := range tg.TileDefs {
		sprites, found := spriteSets[tileDef.SpriteSetRef]
		if !found {
			return fmt.Errorf("missing sprites %s", tileDef.SpriteSetRef)
		}

		img, _, found := sprites.FindSprite(tileDef.StartPoint)
		if !found {
			return fmt.Errorf("no sprite found in sprite set %s at %v",
				tileDef.SpriteSetRef, tileDef.StartPoint)
		}

		rect := img.Bounds()
		dx := rect.Dx()
		dy := rect.Dy()

		if dx != tg.TileSize.Width || dy != tg.TileSize.Height {
			xScales[defID] = float64(tg.TileSize.Width) / float64(dx)
			yScales[defID] = float64(tg.TileSize.Height) / float64(dy)
		} else {
			xScales[defID] = 1
			yScales[defID] = 1
		}

		images[defID] = img
	}

	rows := make([]*Row, len(tg.TileRows))

	columnCounts := mapset.NewSet[int]()

	for i, tileRow := range tg.TileRows {
		defIDs := strings.Split(tileRow, RefIDSeparator)
		nCols := len(defIDs)
		row := NewRow(nCols)

		for j, defID := range defIDs {
			row.Images[j] = images[defID]
			row.XScales[j] = xScales[defID]
			row.YScales[j] = yScales[defID]
		}

		columnCounts.Add(nCols)

		rows[i] = row
	}

	if len(rows) == 0 {
		return fmt.Errorf("grid has no rows")
	}

	switch columnCounts.Cardinality() {
	case 0:
		return fmt.Errorf("grid has no columns")
	case 1:
		// do nothing
	default:
		columnCountStrings := sliceutil.Map(columnCounts.ToSlice(), func(i int) string {
			return strconv.Itoa(i)
		})
		columnCountsStr := strings.Join(columnCountStrings, ", ")

		return fmt.Errorf("grid has inconsistent column counts: %s", columnCountsStr)
	}

	tg.nRows = len(rows)
	tg.nCols = columnCounts.ToSlice()[0]
	tg.rows = rows

	return nil
}

func (tg *TileGrid) WorldLayer() int {
	return drawing.LayerBackground
}

func (tg *TileGrid) VisibleColumns(visible topdown.Rectangle[float64]) (int, int) {
	first := int((visible.Min.X - tg.MinPosition.X) / float64(tg.TileSize.Width))
	last := int((visible.Max.X - tg.MinPosition.X) / float64(tg.TileSize.Width))

	return mathutil.Clamp(first, 0, tg.nCols-1), mathutil.Clamp(last, 0, tg.nCols-1)
}

func (tg *TileGrid) VisibleRows(visible topdown.Rectangle[float64]) (int, int) {
	first := int((visible.Min.Y - tg.MinPosition.Y) / float64(tg.TileSize.Height))
	last := int((visible.Max.Y - tg.MinPosition.Y) / float64(tg.TileSize.Height))

	return mathutil.Clamp(first, 0, tg.nRows-1), mathutil.Clamp(last, 0, tg.nRows-1)
}

func (tg *TileGrid) WorldDraw(worldSurface *ebiten.Image, visible topdown.Rectangle[float64]) {
	// skip drawing if there is no visible portion of the tile grid
	if tg.worldArea.Intersect(visible).Empty() {
		return
	}

	// only draw the visible tiles
	firstColumn, lastColumn := tg.VisibleColumns(visible)
	firstRow, lastRow := tg.VisibleRows(visible)

	for row := firstRow; row <= lastRow; row++ {
		for col := firstColumn; col <= lastColumn; col++ {
			tileImg := tg.rows[row].Images[col]
			opts := &ebiten.DrawImageOptions{}
			sx := tg.rows[row].XScales[col]
			sy := tg.rows[row].YScales[col]
			tx := float64(col * tg.TileSize.Width)
			ty := float64(row * tg.TileSize.Height)

			if sx != 1 || sy != 1 {
				opts.GeoM.Scale(sx, sy)
			}

			opts.GeoM.Translate(tx, ty)

			worldSurface.DrawImage(tileImg, opts)
		}
	}
}
