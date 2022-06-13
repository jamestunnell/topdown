package sprite

import (
	"fmt"
	"strings"

	"github.com/jamestunnell/topdown"
	"github.com/jamestunnell/topdown/resource"
)

const (
	LinkSepID    = "#"
	LinkSepPoint = ":"
)

type Link interface {
	FindSprite(mgr resource.Manager) (*Sprite, bool)
}

type IDLink struct {
	SpriteSheetRef, SpriteID string
}

type PointLink struct {
	SpriteSheetRef string
	Origin         topdown.Point[int]
}

func ParseLink(link string) (l Link, err error) {
	switch {
	case strings.ContainsRune(link, []rune(LinkSepID)[0]):
		pieces := strings.SplitN(link, LinkSepID, 2)

		l = &IDLink{
			SpriteSheetRef: pieces[0],
			SpriteID:       pieces[1],
		}
	case strings.ContainsRune(link, []rune(LinkSepPoint)[0]):
		pieces := strings.SplitN(link, LinkSepPoint, 2)

		pt, parseErr := topdown.ParseIntPoint(pieces[1])
		if parseErr != nil {
			err = parseErr
		} else {
			l = &PointLink{
				SpriteSheetRef: pieces[0],
				Origin:         pt,
			}
		}
	default:
		err = fmt.Errorf("missing type separator")
	}

	return l, err
}

func (l *IDLink) FindSprite(mgr resource.Manager) (*Sprite, bool) {
	sheet, err := resource.GetAs[*Sheet](mgr, l.SpriteSheetRef)
	if err != nil {
		return nil, false
	}

	return sheet.FindSpriteByID(l.SpriteID)
}

func (l *IDLink) String() string {
	return strings.Join([]string{l.SpriteSheetRef, l.SpriteID}, LinkSepID)
}

func (l *PointLink) FindSprite(mgr resource.Manager) (*Sprite, bool) {
	sheet, err := resource.GetAs[*Sheet](mgr, l.SpriteSheetRef)
	if err != nil {
		return nil, false
	}

	return sheet.FindSpriteByOrigin(l.Origin)
}

func (l *PointLink) String() string {
	return strings.Join([]string{l.SpriteSheetRef, l.Origin.String()}, LinkSepPoint)
}
