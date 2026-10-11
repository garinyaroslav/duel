package main

import (
	"embed"
	"fmt"
	_ "image/png"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	windowTitle  = "Duel"
	targetTPS    = 100
)

//go:embed assets/*
var assetsFs embed.FS

func main() {
	typeOfPlayer, _ := prompt.New().Ask("Welcome to Duel! Would you like to be server or client?").
		Choose([]string{"server", "client", "quit"}, choose.WithHelp(true))

	fmt.Print(typeOfPlayer)

	// ebiten.SetWindowSize(screenWidth, screenHeight)
	// ebiten.SetWindowTitle(windowTitle)
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// ebiten.SetFullscreen(true)
	// ebiten.SetVsyncEnabled(true)
	// ebiten.SetTPS(targetTPS)
	//
	// if err := ebiten.RunGame(internal.NewGame(&assetsFs)); err != nil {
	// 	log.Fatal(err)
	// }
}
