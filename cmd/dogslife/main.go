package main

import (
	"context"
	"fmt"

	tg "github.com/luisferreira32/tmgang"
)

func main() {
	eng, err := tg.NewEngine()
	if err != nil {
		panic(err)
	}

	w, h := eng.ScreenSize()

	eng.Configure(&tg.EngineOpts{
		Fps:          10,
		InitialState: walkingDogStateKey,
		StateMachine: map[tg.StateKey]tg.State{
			startingMenuStateKey: &startingMenuState{
				gameLogo: newGameLogoDrawable(w, h),
			},
			walkingDogStateKey: &walkingDogState{},
			gameOverStateKey:   &gameOverState{},
		},
	})

	if err := eng.Run(context.Background()); err != nil {
		fmt.Printf("run failed with: %v\n", err)
	}
}
