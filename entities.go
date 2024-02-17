package tmgang

import (
	"github.com/gdamore/tcell/v2"
)

// Entity is the base of any game.
//
// You will draw entities on your TUI to display your game. Any resize
// from the starting screen size will change the ratioX and ratioY
// variables.
type Entity interface {
	Draw(s tcell.Screen, camera Coordinates, ratioX, ratioY float32)
}

// InteractiveEntity is a special entitity that processes events.
//
// These entities will attempt to process keyboard events if
// they are currently displayed on the screen. As a game developer
// you create entities that implement this interface and give them
// to a state.
type InteractiveEntity interface {
	ProcessKeyEvent(e *tcell.EventKey, camera Coordinates)
}

// Overlays are entities that are drawn at the end, regardless of camera position.
type Overlay interface {
	Draw(s tcell.Screen)
}

// Coordinates of a position
type Coordinates struct {
	X int
	Y int
}

// RectArea defines top left and dimensions of any rectangle area based entitiy.
type RectArea struct {
	Coordinates
	Width  int
	Height int
}

// StrictRectangle is the basic drawing entitiy
type StrictRectangle struct {
	RectArea
	BackgroundStyle, BorderStyle tcell.Style
}

func (r *StrictRectangle) Draw(s tcell.Screen, camera Coordinates, _, _ float32) {
	for row := r.Y; row <= r.Y+r.Height; row++ {
		for col := r.X; col <= r.X+r.Width; col++ {
			s.SetContent(col, row, ' ', nil, r.BackgroundStyle)
		}
	}

	for col := r.X; col <= r.X+r.Width; col++ {
		s.SetContent(col, r.Y, tcell.RuneHLine, nil, r.BorderStyle)
		s.SetContent(col, r.Y+r.Height, tcell.RuneHLine, nil, r.BorderStyle)
	}
	for row := r.Y; row <= r.Y+r.Height; row++ {
		s.SetContent(r.X, row, tcell.RuneVLine, nil, r.BorderStyle)
		s.SetContent(r.X+r.Width, row, tcell.RuneVLine, nil, r.BorderStyle)
	}

	if r.Height != 0 && r.Width != 0 {
		s.SetContent(r.X, r.Y, tcell.RuneULCorner, nil, r.BorderStyle)
		s.SetContent(r.X+r.Width, r.Y, tcell.RuneURCorner, nil, r.BorderStyle)
		s.SetContent(r.X, r.Y+r.Height, tcell.RuneLLCorner, nil, r.BorderStyle)
		s.SetContent(r.X+r.Width, r.Y+r.Height, tcell.RuneLRCorner, nil, r.BorderStyle)
	}
}

// FlexRectangle is a basic drawing entitiy that adjusts with the resize
type FlexRectangle struct {
	RectArea

	OriginalRect                 RectArea
	BackgroundStyle, BorderStyle tcell.Style

	ratioX, ratioY float32
}

func (r *FlexRectangle) Draw(s tcell.Screen, camera Coordinates, ratioX, ratioY float32) {
	if r.ratioX != ratioX {
		r.ratioX = ratioX
		r.X = int(float32(r.OriginalRect.X) * ratioX)
		r.Width = int(float32(r.OriginalRect.Width) * ratioX)
	}
	if r.ratioY != ratioY {
		r.ratioY = ratioY
		r.Y = int(float32(r.OriginalRect.Y) * ratioX)
		r.Height = int(float32(r.OriginalRect.Height) * ratioX)
	}

	for row := r.Y; row <= r.Y+r.Height; row++ {
		for col := r.X; col <= r.X+r.Width; col++ {
			s.SetContent(col, row, ' ', nil, r.BackgroundStyle)
		}
	}

	for col := r.X; col <= r.X+r.Width; col++ {
		s.SetContent(col, r.Y, tcell.RuneHLine, nil, r.BorderStyle)
		s.SetContent(col, r.Y+r.Height, tcell.RuneHLine, nil, r.BorderStyle)
	}
	for row := r.Y; row <= r.Y+r.Height; row++ {
		s.SetContent(r.X, row, tcell.RuneVLine, nil, r.BorderStyle)
		s.SetContent(r.X+r.Width, row, tcell.RuneVLine, nil, r.BorderStyle)
	}

	if r.Height != 0 && r.Width != 0 {
		s.SetContent(r.X, r.Y, tcell.RuneULCorner, nil, r.BorderStyle)
		s.SetContent(r.X+r.Width, r.Y, tcell.RuneURCorner, nil, r.BorderStyle)
		s.SetContent(r.X, r.Y+r.Height, tcell.RuneLLCorner, nil, r.BorderStyle)
		s.SetContent(r.X+r.Width, r.Y+r.Height, tcell.RuneLRCorner, nil, r.BorderStyle)
	}
}

// Text writes bounded unscrollable text
type Text struct {
	RectArea
	TextStyle tcell.Style
	Content   string
}

func (t *Text) Draw(s tcell.Screen, camera Coordinates, _, _ float32) {
	row := t.Y
	col := t.X
	for _, r := range []rune(t.Content) {
		s.SetContent(col, row, r, nil, t.TextStyle)
		col++
		if col >= t.X+t.Width {
			row++
			col = t.X
		}
		if row > t.Y+t.Height {
			break
		}
	}
}

// BasicMenu will, given some options, allow selection with a * pointer.
type BasicMenu struct {
	RectArea
	TextStyle, BorderStyle tcell.Style

	MenuItems       []string
	CurrentSelected int
}

func (bm *BasicMenu) Draw(s tcell.Screen, camera Coordinates, _, _ float32) {
	for row := bm.Y; row <= bm.Y+bm.Height; row++ {
		for col := bm.X; col <= bm.X+bm.Width; col++ {
			switch {
			case col == bm.X+1 && bm.CurrentSelected == row-bm.Y-1:
				s.SetContent(col, row, '*', nil, bm.TextStyle)
			case col >= bm.X+3 && row >= bm.Y+1 && len(bm.MenuItems) > row-bm.Y-1 && len(bm.MenuItems[row-bm.Y-1]) > col-bm.X-3:
				s.SetContent(col, row, rune(bm.MenuItems[row-bm.Y-1][col-bm.X-3]), nil, bm.TextStyle)
			default:
				s.SetContent(col, row, ' ', nil, bm.TextStyle)
			}
		}
	}

	for col := bm.X; col <= bm.X+bm.Width; col++ {
		s.SetContent(col, bm.Y, tcell.RuneHLine, nil, bm.BorderStyle)
		s.SetContent(col, bm.Y+bm.Height, tcell.RuneHLine, nil, bm.BorderStyle)
	}
	for row := bm.Y; row <= bm.Y+bm.Height; row++ {
		s.SetContent(bm.X, row, tcell.RuneVLine, nil, bm.BorderStyle)
		s.SetContent(bm.X+bm.Width, row, tcell.RuneVLine, nil, bm.BorderStyle)
	}

	if bm.Height != 0 && bm.Width != 0 {
		s.SetContent(bm.X, bm.Y, tcell.RuneULCorner, nil, bm.BorderStyle)
		s.SetContent(bm.X+bm.Width, bm.Y, tcell.RuneURCorner, nil, bm.BorderStyle)
		s.SetContent(bm.X, bm.Y+bm.Height, tcell.RuneLLCorner, nil, bm.BorderStyle)
		s.SetContent(bm.X+bm.Width, bm.Y+bm.Height, tcell.RuneLRCorner, nil, bm.BorderStyle)
	}
}

// bm processing event
func (bm *BasicMenu) ProcessKeyEvent(e *tcell.EventKey, camera Coordinates) {
	if e.Key() == tcell.KeyDown && bm.CurrentSelected < len(bm.MenuItems)-1 {
		bm.CurrentSelected++
	}
	if e.Key() == tcell.KeyUp && bm.CurrentSelected > 0 {
		bm.CurrentSelected--
	}
}
