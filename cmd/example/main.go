package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	tg "github.com/luisferreira32/tmgang"
)

const (
	mainMenuStateKey tg.StateKey = "mainMenu"
)

type mainMenuState struct {
	background  *tg.StrictRectangle
	menuEntitiy *tg.BasicMenu
	flexRect    *tg.FlexRectangle
	flexChatBox *tg.FlexChatBox
}

func (m *mainMenuState) GetDrawables() []tg.Entity {
	return []tg.Entity{m.background, m.menuEntitiy, m.flexRect, m.flexChatBox}
}

func (m *mainMenuState) GetInteractables() []tg.InteractiveEntity {
	return []tg.InteractiveEntity{m.menuEntitiy, m.flexChatBox}
}

func (*mainMenuState) GetCamera() tg.Coordinates {
	// basically means we're aligned with top/left of visible screen - no adjustments
	return tg.Coordinates{X: 0, Y: 0}
}

func (*mainMenuState) NextState() tg.StateKey {
	return mainMenuStateKey
}

const (
	longMenuText = `omg i'm quite a big sentence since I want to test how the texting wraps around the box and manages to go over to the next action... and in fact the text just keeps on going since this is but an example on how to use the FlexChatBox entity.`
)

func createMainMenu(w, h int) *mainMenuState {
	flexRectArea := tg.RectArea{
		Coordinates: tg.Coordinates{X: 10 + w/3, Y: 10},
		Width:       50,
		Height:      h / 8,
	}
	flexRect := &tg.FlexRectangle{
		RectArea:     flexRectArea,
		OriginalRect: flexRectArea,
	}

	flexChatBoxArea := tg.RectArea{
		Coordinates: tg.Coordinates{X: w - 100, Y: 10},
		Width:       50,
		Height:      5,
	}

	return &mainMenuState{
		background: &tg.StrictRectangle{
			RectArea: tg.RectArea{
				Coordinates: tg.Coordinates{
					X: -1,
					Y: -1,
				},
				Width:  w + 1,
				Height: h + 1,
			},
			Style: tcell.StyleDefault.Background(0),
		},
		menuEntitiy: &tg.BasicMenu{
			RectArea: tg.RectArea{
				Coordinates: tg.Coordinates{
					X: 10,
					Y: 10,
				},
				Width:  w / 9,
				Height: h / 8,
			},
			Style:           tcell.StyleDefault.Background(0),
			CurrentSelected: 0,
			MenuItems: []string{
				"foobar",
				"barfoobar",
				"foo",
			},
		},
		flexRect: flexRect,
		flexChatBox: &tg.FlexChatBox{
			RectArea:     flexChatBoxArea,
			OriginalRect: flexChatBoxArea,
			Style:        tcell.StyleDefault.Background(0),
			Content:      strings.Split(longMenuText, " "),
		},
	}
}

func main() {
	eng, err := tg.NewEngine()
	if err != nil {
		panic(err)
	}

	w, h := eng.ScreenSize()
	mainMenu := createMainMenu(w, h)

	eng.Configure(&tg.EngineOpts{
		Fps:          10,
		InitialState: mainMenuStateKey,
		StateMachine: map[tg.StateKey]tg.State{
			mainMenuStateKey: mainMenu,
		},
	})

	if err := eng.Run(); err != nil {
		fmt.Printf("run failed with: %v\n", err)
	}
}
