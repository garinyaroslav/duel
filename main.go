package main

import (
	"embed"
	"fmt"
	_ "image/png"
	"log"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/input"
	"github.com/garinyaroslav/duel/internal"
	"github.com/hajimehoshi/ebiten/v2"
)

const DEFAULT_PORT = "7777"

//go:embed assets/*
var assetsFs embed.FS

func main() {
	internal.InitLogger()

	role, err := prompt.New().Ask("Welcome to Duel! Would you like to be server or client?").
		Choose([]string{"server", "client", "quit"}, choose.WithHelp(true))
	if err != nil {
		log.Fatal(err)
	}
	slog.Debug("User choice", "role", role)

	var addr, mode string

	switch role {
	case "server":
		mode = "host"
		addr = ":" + DEFAULT_PORT

	case "client":
		mode = "client"

		serverIP, err := prompt.New().Ask("Enter server IP (port "+DEFAULT_PORT+" by default):").
			Input("127.0.0.1", input.WithHelp(true))
		if err != nil {
			log.Fatal(err)
		}
		addr = normalizeAddr(serverIP)

	default:
		fmt.Println("Exiting game...")
		os.Exit(0)
	}

	game, err := internal.NewGame(&assetsFs, mode, addr)
	if err != nil {
		log.Fatal(err)
	}

	if mode == "host" {
		fmt.Println("Server started. Tell your opponent one of these addresses:")
		for _, ip := range localIPs() {
			fmt.Printf("  %s\n", ip)
		}
		fmt.Printf("Waiting for opponent on port %s... (Ctrl+C to cancel)\n", DEFAULT_PORT)

		game.WaitForOpponent()
		fmt.Println("Opponent connected, starting!")
	}

	ebiten.SetWindowSize(internal.SCREEN_W, internal.SCREEN_H)
	ebiten.SetWindowTitle("Duel — " + role)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(internal.TPS)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func normalizeAddr(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		s = "127.0.0.1"
	}
	if _, _, err := net.SplitHostPort(s); err == nil {
		return s
	}
	return net.JoinHostPort(s, DEFAULT_PORT)
}

func localIPs() []string {
	var res []string

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return res
	}
	for _, a := range addrs {
		ipn, ok := a.(*net.IPNet)
		if ok && !ipn.IP.IsLoopback() && ipn.IP.To4() != nil {
			res = append(res, ipn.IP.String())
		}
	}
	return res
}
