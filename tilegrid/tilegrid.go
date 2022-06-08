package tilegrid

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/drawing"
	"github.com/jamestunnell/topdown/imageset"
	"github.com/jamestunnell/topdown/mathutil"
	"github.com/jamestunnell/topdown/resource"
	"github.com/jamestunnell/topdown/sliceutil"
)

type TileGrid struct {
	MinPosition topdown.Vector     `json:"minPosition"`
	TileSize    topdown.Size       `json:"tileSize"`
	TileDefs    map[string]TileDef `json:"tileDefs"`
	TileRows    []string           `json:"tileRows"`

	imageSets  map[string]*imageset.ImageSet
	imgRows    [][]*ebiten.Image
	xScaleRows [][]float64
	yScaleRows [][]float64
	nRows      int
	nCols      int
}

type TileDef struct {
	ImageSetRef string        `json:"imageSetRef"`
	StartPixel  topdown.Pixel `json:"startPixel"`
}

const (
	RefIDSeparator = " "
	TypeName       = "tilegrid"
)

func New(tileSize topdown.Size) *TileGrid {
	tg := &TileGrid{
		MinPosition: topdown.NewVector(0, 0),
		TileSize:    tileSize,
		TileDefs:    map[string]TileDef{},
		TileRows:    []string{},
		imageSets:   map[string]*imageset.ImageSet{},
		imgRows:     [][]*ebiten.Image{},
	}

	return tg
}

func (tg *TileGrid) Initialize(mgr resource.Manager) error {
	imageSets := map[string]*imageset.ImageSet{}

	for _, tileDef := range tg.TileDefs {
		if _, found := imageSets[tileDef.ImageSetRef]; found {
			continue
		}

		is, err := resource.GetAs[*imageset.ImageSet](mgr, tileDef.ImageSetRef)
		if err != nil {
			return fmt.Errorf("failed to get '%s' from dependencies: %w", tileDef.ImageSetRef, err)
		}

		imageSets[tileDef.ImageSetRef] = is
	}

	tg.imageSets = imageSets

	if err := tg.MakeImageRows(); err != nil {
		return fmt.Errorf("failed to make image rows: %w", err)
	}

	return nil
}

func (tg *TileGrid) WorldArea() image.Rectangle {
	x := int(tg.MinPosition.X)
	y := int(tg.MinPosition.Y)
	dX := tg.nCols * int(tg.TileSize.Width)
	dY := tg.nRows * int(tg.TileSize.Height)

	return image.Rect(x, y, x+dX, y+dY)
}

func (tg *TileGrid) Center() image.Point {
	if len(tg.imgRows) == 0 || len(tg.imgRows[0]) == 0 {
		return image.Point{}
	}

	cX := len(tg.imgRows[0]) * int(tg.TileSize.Width) / 2
	cY := len(tg.imgRows) * int(tg.TileSize.Height) / 2

	return image.Pt(cX, cY)
}

func (tg *TileGrid) MakeImageRows() error {
	images := map[string]*ebiten.Image{}
	xScales := map[string]float64{}
	yScales := map[string]float64{}

	for refID, tileDef := range tg.TileDefs {
		is, found := tg.imageSets[tileDef.ImageSetRef]
		if !found {
			return fmt.Errorf("missing imageset %s", tileDef.ImageSetRef)
		}

		img, _, subImageFound := is.SubImage(tileDef.StartPixel)
		if !subImageFound {
			return fmt.Errorf("no sub-image found in imageset %s at %s",
				tileDef.ImageSetRef, tileDef.StartPixel.String())
		}

		rect := img.Bounds()
		dx := float64(rect.Dx())
		dy := float64(rect.Dy())

		if dx != tg.TileSize.Width || dy != tg.TileSize.Height {
			xScales[refID] = tg.TileSize.Width / dx
			yScales[refID] = tg.TileSize.Height / dy
		} else {
			xScales[refID] = 1
			yScales[refID] = 1
		}

		images[refID] = img
	}

	imgRows := [][]*ebiten.Image{}
	xScaleRows := [][]float64{}
	yScaleRows := [][]float64{}

	columnCounts := mapset.NewSet[int]()

	for _, tileRow := range tg.TileRows {
		refIDs := strings.Split(tileRow, RefIDSeparator)
		imgRow := make([]*ebiten.Image, len(refIDs))
		xScaleRow := make([]float64, len(refIDs))
		yScaleRow := make([]float64, len(refIDs))

		for i, refID := range refIDs {
			imgRow[i] = images[refID]
			xScaleRow[i] = xScales[refID]
			yScaleRow[i] = yScales[refID]
		}

		columnCounts.Add(len(refIDs))

		imgRows = append(imgRows, imgRow)
		xScaleRows = append(xScaleRows, xScaleRow)
		yScaleRows = append(yScaleRows, yScaleRow)
	}

	if len(imgRows) == 0 {
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

	tg.nRows = len(imgRows)
	tg.nCols = columnCounts.ToSlice()[0]
	tg.imgRows = imgRows
	tg.xScaleRows = xScaleRows
	tg.yScaleRows = yScaleRows

	return nil
}

func (tg *TileGrid) WorldLayer() int {
	return drawing.LayerBackground
}

func (tg *TileGrid) visibleColumns(visible image.Rectangle) (int, int) {
	first := mathutil.Clamp(
		int((float64(visible.Min.X)-tg.MinPosition.X)/tg.TileSize.Width), 0, tg.nCols-1)
	last := mathutil.Clamp(
		int((float64(visible.Max.X)-tg.MinPosition.X)/tg.TileSize.Width), 0, tg.nCols-1)

	return first, last
}

func (tg *TileGrid) visibleRows(visible image.Rectangle) (int, int) {
	first := mathutil.Clamp(
		int((float64(visible.Min.Y)-tg.MinPosition.Y)/tg.TileSize.Height), 0, tg.nRows-1)
	last := mathutil.Clamp(
		int((float64(visible.Max.Y)-tg.MinPosition.Y)/tg.TileSize.Height), 0, tg.nRows-1)

	return first, last
}

func (tg *TileGrid) WorldDraw(worldSurface *ebiten.Image, visible image.Rectangle) {
	// skip drawing if there is no visible portion of the tile grid
	if tg.WorldArea().Intersect(visible).Empty() {
		return
	}

	// only draw the visible tiles
	firstColumn, lastColumn := tg.visibleColumns(visible)
	firstRow, lastRow := tg.visibleRows(visible)

	for row := firstRow; row <= lastRow; row++ {
		for col := firstColumn; col <= lastColumn; col++ {
			tileImg := tg.imgRows[row][col]
			opts := &ebiten.DrawImageOptions{}
			sx := tg.xScaleRows[row][col]
			sy := tg.yScaleRows[row][col]
			tx := float64(col) * tg.TileSize.Width
			ty := float64(row) * tg.TileSize.Height

			if sx != 1 || sy != 1 {
				opts.GeoM.Scale(sx, sy)
			}

			opts.GeoM.Translate(tx, ty)

			worldSurface.DrawImage(tileImg, opts)
		}
	}
}
