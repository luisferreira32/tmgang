package tmgang

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	textPaddingX = 3
	textPaddingY = 2
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

// TimeEntity is a special entitity that processes the passage of time.
//
// These entities will attempt to process the passage of time for each
// frame, an example for such usage would be for entities that have a
// speed associated with them.
type TimeEntity interface {
	ProcessFrameDuration(d time.Duration, camera Coordinates)
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

// StrictRectangle is the basic drawing entitiy.
type StrictRectangle struct {
	RectArea
	Style tcell.Style
}

func (r *StrictRectangle) Draw(s tcell.Screen, camera Coordinates, _, _ float32) {
	for row := r.Y; row <= r.Y+r.Height; row++ {
		for col := r.X; col <= r.X+r.Width; col++ {
			s.SetContent(col, row, ' ', nil, r.Style)
		}
	}

	for col := r.X; col <= r.X+r.Width; col++ {
		s.SetContent(col, r.Y, tcell.RuneHLine, nil, r.Style)
		s.SetContent(col, r.Y+r.Height, tcell.RuneHLine, nil, r.Style)
	}
	for row := r.Y; row <= r.Y+r.Height; row++ {
		s.SetContent(r.X, row, tcell.RuneVLine, nil, r.Style)
		s.SetContent(r.X+r.Width, row, tcell.RuneVLine, nil, r.Style)
	}

	if r.Height != 0 && r.Width != 0 {
		s.SetContent(r.X, r.Y, tcell.RuneULCorner, nil, r.Style)
		s.SetContent(r.X+r.Width, r.Y, tcell.RuneURCorner, nil, r.Style)
		s.SetContent(r.X, r.Y+r.Height, tcell.RuneLLCorner, nil, r.Style)
		s.SetContent(r.X+r.Width, r.Y+r.Height, tcell.RuneLRCorner, nil, r.Style)
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

// BasicText writes bounded unscrollable text
type BasicText struct {
	RectArea
	TextStyle tcell.Style
	Content   string
}

func (t *BasicText) Draw(s tcell.Screen, camera Coordinates, _, _ float32) {
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
	Style tcell.Style

	MenuItems       []string
	CurrentSelected int
}

func (bm *BasicMenu) Draw(s tcell.Screen, camera Coordinates, _, _ float32) {
	for row := bm.Y; row <= bm.Y+bm.Height; row++ {
		for col := bm.X; col <= bm.X+bm.Width; col++ {
			switch {
			case col == bm.X+1 && bm.CurrentSelected == row-bm.Y-1:
				s.SetContent(col, row, '*', nil, bm.Style)
			case col >= bm.X+3 && row >= bm.Y+1 && len(bm.MenuItems) > row-bm.Y-1 && len(bm.MenuItems[row-bm.Y-1]) > col-bm.X-3:
				s.SetContent(col, row, rune(bm.MenuItems[row-bm.Y-1][col-bm.X-3]), nil, bm.Style)
			default:
				s.SetContent(col, row, ' ', nil, bm.Style)
			}
		}
	}

	for col := bm.X; col <= bm.X+bm.Width; col++ {
		s.SetContent(col, bm.Y, tcell.RuneHLine, nil, bm.Style)
		s.SetContent(col, bm.Y+bm.Height, tcell.RuneHLine, nil, bm.Style)
	}
	for row := bm.Y; row <= bm.Y+bm.Height; row++ {
		s.SetContent(bm.X, row, tcell.RuneVLine, nil, bm.Style)
		s.SetContent(bm.X+bm.Width, row, tcell.RuneVLine, nil, bm.Style)
	}

	if bm.Height != 0 && bm.Width != 0 {
		s.SetContent(bm.X, bm.Y, tcell.RuneULCorner, nil, bm.Style)
		s.SetContent(bm.X+bm.Width, bm.Y, tcell.RuneURCorner, nil, bm.Style)
		s.SetContent(bm.X, bm.Y+bm.Height, tcell.RuneLLCorner, nil, bm.Style)
		s.SetContent(bm.X+bm.Width, bm.Y+bm.Height, tcell.RuneLRCorner, nil, bm.Style)
	}
}

func (bm *BasicMenu) ProcessKeyEvent(e *tcell.EventKey, camera Coordinates) {
	if e.Key() == tcell.KeyDown && bm.CurrentSelected < len(bm.MenuItems)-1 {
		bm.CurrentSelected++
	}
	if e.Key() == tcell.KeyUp && bm.CurrentSelected > 0 {
		bm.CurrentSelected--
	}
}

// FlexChatBox writes bounded "scrollable" text when the Space key is pressed.
type FlexChatBox struct {
	RectArea

	OriginalRect   RectArea
	ratioX, ratioY float32
	Style          tcell.Style

	currentWord int
	nextWord    int
	Content     []string
}

func (r *FlexChatBox) Draw(s tcell.Screen, camera Coordinates, ratioX, ratioY float32) {
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
			s.SetContent(col, row, ' ', nil, r.Style)
		}
	}

	r.nextWord = 0
	idx, row, col := 0, r.Y+textPaddingY, r.X+textPaddingX
	for {
		if r.currentWord+r.nextWord >= len(r.Content) { // reached the end of the content to draw
			break
		}
		if row+textPaddingY > r.Y+r.Height { // cannot draw more on this flexbox to respect the padding
			break
		}

		if idx == 0 && col+len(r.Content[r.currentWord+r.nextWord])+textPaddingX > r.X+r.Width { // this word would be split at the end, wrap it
			col = r.X + textPaddingX
			row++
			continue
		}

		if idx >= len(r.Content[r.currentWord+r.nextWord]) { // we finished a word, space and go next.
			s.SetContent(col, row, ' ', nil, r.Style)
			col++
			idx = 0
			r.nextWord++
			continue
		}

		s.SetContent(col, row, rune(r.Content[r.currentWord+r.nextWord][idx]), nil, r.Style)
		col++
		idx++
	}

	if len(r.Content) > r.currentWord+r.nextWord { // we have not finished the text
		s.SetContent(r.X+r.Width-1, r.Y+r.Height-1, '.', nil, r.Style)
		s.SetContent(r.X+r.Width-2, r.Y+r.Height-1, '.', nil, r.Style)
		s.SetContent(r.X+r.Width-3, r.Y+r.Height-1, '.', nil, r.Style)
	}

	for col := r.X; col <= r.X+r.Width; col++ {
		s.SetContent(col, r.Y, tcell.RuneHLine, nil, r.Style)
		s.SetContent(col, r.Y+r.Height, tcell.RuneHLine, nil, r.Style)
	}
	for row := r.Y; row <= r.Y+r.Height; row++ {
		s.SetContent(r.X, row, tcell.RuneVLine, nil, r.Style)
		s.SetContent(r.X+r.Width, row, tcell.RuneVLine, nil, r.Style)
	}

	if r.Height != 0 && r.Width != 0 {
		s.SetContent(r.X, r.Y, tcell.RuneULCorner, nil, r.Style)
		s.SetContent(r.X+r.Width, r.Y, tcell.RuneURCorner, nil, r.Style)
		s.SetContent(r.X, r.Y+r.Height, tcell.RuneLLCorner, nil, r.Style)
		s.SetContent(r.X+r.Width, r.Y+r.Height, tcell.RuneLRCorner, nil, r.Style)
	}
}

func (r *FlexChatBox) ProcessKeyEvent(e *tcell.EventKey, camera Coordinates) {
	if e.Rune() == ' ' && len(r.Content) > r.currentWord+r.nextWord {
		r.currentWord += r.nextWord
	}
}
