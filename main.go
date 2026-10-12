package main

import (
	"embed"
	"fmt"
	_ "image/png"
	"log/slog"
	"os"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/garinyaroslav/duel/internal"
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
	internal.InitLogger()

	typeOfPlayer, _ := prompt.New().Ask("Welcome to Duel! Would you like to be server or client?").
		Choose([]string{"server", "client", "quit"}, choose.WithHelp(true))

	slog.Debug("User choice", "typeOfPlayer", typeOfPlayer)

	switch typeOfPlayer {
	case "server":
		fmt.Println("Starting server...")
	case "client":
		fmt.Println("Starting client...")
	case "quit":
		fmt.Println("Exiting game...")
		os.Exit(0)
		return
	default:
		fmt.Println("Invalid choice. Exiting game...")
		os.Exit(0)
		return
	}

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
