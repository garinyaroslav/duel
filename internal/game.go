package internal

import (
	"embed"
	"fmt"
	"image/color"

	"github.com/garinyaroslav/duel/internal/entity"
	"github.com/garinyaroslav/duel/pkg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenW, screenH = 1920, 1080
	playerX, playerY = 0, 0
)

type Game struct {
	background color.RGBA
	player     entity.Player
}

func NewGame(AssetFs embed.FS) *Game {
	return &Game{
		background: color.RGBA{0, 181, 226, 255},
		player:     *entity.NewPlayer(0, 0, pkg.LoadImage("assets/player.png", AssetFs)),
	}
}

func (g *Game) Update() error {
	g.player.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.background)

	X, Y := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %v, Y: %v, FPS: %v", X, Y, ebiten.ActualFPS()))

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.player.Position.X, g.player.Position.Y)
	screen.DrawImage(g.player.Sprite, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenW, screenH
}
