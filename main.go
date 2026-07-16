package main

import (
	"embed"
	_ "image/png"
	"log"

	"github.com/garinyaroslav/duel/internal"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	windowTitle  = "Duel"
	targetTPS    = 100
)

//go:embed assets/*
var AssetsFs embed.FS

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(true)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(targetTPS)

	if err := ebiten.RunGame(internal.NewGame(AssetsFs)); err != nil {
		log.Fatal(err)
	}
}
