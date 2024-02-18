package main

import (
	"github.com/gdamore/tcell/v2"
	tg "github.com/luisferreira32/tmgang"
)

func newGameLogoDrawable(w, h int) *gameLogoDrawable {
	return &gameLogoDrawable{}
}

type gameLogoDrawable struct {
}

func (g *gameLogoDrawable) Draw(s tcell.Screen, camera tg.Coordinates, _, _ float32) {}
