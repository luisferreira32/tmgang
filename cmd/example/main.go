package main

import (
	"fmt"

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
}

func (m *mainMenuState) GetDrawables() []tg.Entity {
	return []tg.Entity{m.background, m.menuEntitiy, m.flexRect}
}

func (m *mainMenuState) GetInteractables() []tg.InteractiveEntity {
	return []tg.InteractiveEntity{m.menuEntitiy}
}

func (*mainMenuState) GetCamera() tg.Coordinates {
	// basically means we're aligned with top/left of visible screen - no adjustments
	return tg.Coordinates{X: 0, Y: 0}
}

func (*mainMenuState) NextState() tg.StateKey {
	return mainMenuStateKey
}

func createMainMenu(w, h int) *mainMenuState {
	flexRectArea := tg.RectArea{
		Coordinates: tg.Coordinates{X: 10 + w/3, Y: 10},
		Width:       w / 3,
		Height:      h / 8,
	}
	flexRect := &tg.FlexRectangle{
		RectArea:     flexRectArea,
		OriginalRect: flexRectArea,
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
			BackgroundStyle: tcell.StyleDefault.Background(0),
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
			TextStyle:       tcell.StyleDefault.Background(0),
			BorderStyle:     tcell.StyleDefault.Background(0),
			CurrentSelected: 0,
			MenuItems: []string{
				"foobar",
				"barfoobar",
				"foo",
			},
		},
		flexRect: flexRect,
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
