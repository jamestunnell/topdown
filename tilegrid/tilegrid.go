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
	Origin    topdown.Point[float64] `json:"origin"`
	TileSize  topdown.Size[int]      `json:"tileSize"`
	TileLinks map[string]string      `json:"tileLinks"`
	TileRows  []string               `json:"tileRows"`

	worldArea topdown.Rectangle[float64]
	center    topdown.Point[float64]
	rows      []*Row
	nRows     int
	nCols     int
}

type Tile struct {
	Image          *ebiten.Image
	XScale, YScale float64
}

type Row struct {
	Tiles []*Tile
}

const (
	RefIDSeparator = " "
)

func New(tileSize topdown.Size[int]) *TileGrid {
	tg := &TileGrid{
		Origin:    topdown.Pt[float64](0, 0),
		TileSize:  tileSize,
		TileLinks: map[string]string{},
		TileRows:  []string{},
		rows:      []*Row{},
	}

	return tg
}

func (tg *TileGrid) Initialize(mgr resource.Manager) error {
	tiles := map[string]*Tile{}

	for tileID, spriteLink := range tg.TileLinks {
		l, err := sprite.ParseLink(spriteLink)
		if err != nil {
			return fmt.Errorf("failed to parse sprite link '%s': %w", spriteLink, err)
		}

		sprite, found := l.FindSprite(mgr)
		if !found {
			return fmt.Errorf("failed to find sprite with link '%s'", spriteLink)
		}

		rect := sprite.Image.Bounds()
		dx := rect.Dx()
		dy := rect.Dy()

		tile := &Tile{
			Image:  sprite.Image,
			XScale: 1.0,
			YScale: 1.0,
		}

		if dx != tg.TileSize.Width || dy != tg.TileSize.Height {
			tile.XScale = float64(tg.TileSize.Width) / float64(dx)
			tile.YScale = float64(tg.TileSize.Height) / float64(dy)
		}

		tiles[tileID] = tile
	}

	if err := tg.MakeRows(tiles); err != nil {
		return fmt.Errorf("failed to make rows: %w", err)
	}

	dX := float64(tg.nCols * tg.TileSize.Width)
	dY := float64(tg.nRows * tg.TileSize.Height)
	tg.worldArea = topdown.Rectangle[float64]{
		Min: topdown.Pt(tg.Origin.X, tg.Origin.Y),
		Max: topdown.Pt(tg.Origin.X+dX, tg.Origin.Y+dY),
	}

	tg.center = tg.worldArea.Min.Add(tg.worldArea.Size().Center())

	return nil
}

func (tg *TileGrid) Center() topdown.Point[float64] {
	return tg.worldArea.Size().Center()
}

func (tg *TileGrid) MakeRows(tiles map[string]*Tile) error {
	rows := make([]*Row, len(tg.TileRows))

	columnCounts := mapset.NewSet[int]()

	for i, tileRow := range tg.TileRows {
		tileIDs := strings.Split(tileRow, RefIDSeparator)
		nCols := len(tileIDs)
		rowTiles := make([]*Tile, nCols)

		for j, tileID := range tileIDs {
			tile, found := tiles[tileID]
			if !found {
				return fmt.Errorf("tile '%s' not defined", tileID)
			}

			rowTiles[j] = tile
		}

		columnCounts.Add(nCols)

		rows[i] = &Row{Tiles: rowTiles}
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
	first := int((visible.Min.X - tg.Origin.X) / float64(tg.TileSize.Width))
	last := int((visible.Max.X - tg.Origin.X) / float64(tg.TileSize.Width))

	return mathutil.Clamp(first, 0, tg.nCols-1), mathutil.Clamp(last, 0, tg.nCols-1)
}

func (tg *TileGrid) VisibleRows(visible topdown.Rectangle[float64]) (int, int) {
	first := int((visible.Min.Y - tg.Origin.Y) / float64(tg.TileSize.Height))
	last := int((visible.Max.Y - tg.Origin.Y) / float64(tg.TileSize.Height))

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
			tileImg := tg.rows[row].Tiles[col].Image
			opts := &ebiten.DrawImageOptions{}
			sx := tg.rows[row].Tiles[col].XScale
			sy := tg.rows[row].Tiles[col].YScale
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
