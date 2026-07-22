package internal

import (
	"embed"
	"fmt"
	"image/color"

	"github.com/garinyaroslav/duel/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenW, screenH = 1920, 1080
	playerX, playerY = 0, 0
	NumberOfPlayers  = 2
)

type Game struct {
	background *color.RGBA
	players    [NumberOfPlayers]*entity.Player
}

func NewGame(assetFs *embed.FS) *Game {
	return &Game{
		background: &color.RGBA{0, 181, 226, 255},
		players:    [NumberOfPlayers]*entity.Player{entity.NewPlayer(100, 100, 0, assetFs), entity.NewPlayer(1820, 100, 1, assetFs)},
	}
}

func (g *Game) Update() error {
	updatePlayers(g)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.background)

	drawPlayers(g, screen)

	drawDebugInfo(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenW, screenH
}

func updatePlayers(g *Game) {
	for _, p := range g.players {
		p.Update()
	}
}

func drawPlayers(g *Game, screen *ebiten.Image) {
	for _, p := range g.players {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.Position.X-40, p.Position.Y-40)
		screen.DrawImage(p.Sprite, op)

		for _, proj := range p.Projectiles {
			if !proj.Active {
				return
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(proj.Position.X-10, proj.Position.Y-10)
			screen.DrawImage(p.ProjectileSprite, op)
		}
	}
}

func drawDebugInfo(screen *ebiten.Image) {
	X, Y := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %v, Y: %v, FPS: %v", X, Y, ebiten.ActualFPS()))
}
