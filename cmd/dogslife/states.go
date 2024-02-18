package main

import (
	tg "github.com/luisferreira32/tmgang"
)

const (
	startingMenuStateKey tg.StateKey = "startingMenu"
	walkingDogStateKey   tg.StateKey = "walkingDog"
	gameOverStateKey     tg.StateKey = "gameOver"
)

type startingMenuState struct {
	gameLogo *gameLogoDrawable
}

func (m *startingMenuState) GetDrawables() []tg.Entity {
	return []tg.Entity{m.gameLogo}
}

func (m *startingMenuState) GetInteractables() []tg.InteractiveEntity {
	return []tg.InteractiveEntity{}
}

func (m *startingMenuState) GetTimers() []tg.TimeEntity {
	return []tg.TimeEntity{}
}

func (*startingMenuState) GetCamera() tg.Coordinates {
	// basically means we're aligned with top/left of visible screen - no adjustments
	return tg.Coordinates{X: 0, Y: 0}
}

func (*startingMenuState) NextState() tg.StateKey {
	return startingMenuStateKey
}

type walkingDogState struct {
}

func (m *walkingDogState) GetDrawables() []tg.Entity {
	return []tg.Entity{}
}

func (m *walkingDogState) GetInteractables() []tg.InteractiveEntity {
	return []tg.InteractiveEntity{}
}

func (m *walkingDogState) GetTimers() []tg.TimeEntity {
	return []tg.TimeEntity{}
}

func (*walkingDogState) GetCamera() tg.Coordinates {
	// basically means we're aligned with top/left of visible screen - no adjustments
	return tg.Coordinates{X: 0, Y: 0}
}

func (*walkingDogState) NextState() tg.StateKey {
	return walkingDogStateKey
}

type gameOverState struct {
}

func (m *gameOverState) GetDrawables() []tg.Entity {
	return []tg.Entity{}
}

func (m *gameOverState) GetInteractables() []tg.InteractiveEntity {
	return []tg.InteractiveEntity{}
}

func (m *gameOverState) GetTimers() []tg.TimeEntity {
	return []tg.TimeEntity{}
}

func (*gameOverState) GetCamera() tg.Coordinates {
	// basically means we're aligned with top/left of visible screen - no adjustments
	return tg.Coordinates{X: 0, Y: 0}
}

func (*gameOverState) NextState() tg.StateKey {
	return gameOverStateKey
}
