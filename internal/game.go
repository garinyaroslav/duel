package internal

import (
	"embed"
	"fmt"
	"image/color"

	"github.com/garinyaroslav/duel/pkg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenW, screenH = 1920, 1080
	playerX, playerY = 0, 0
)

type Game struct {
	backgroundColor  color.RGBA
	playerSprite     *ebiten.Image
	playerX, playerY float64
	// playerVelocity   float64
}

func NewGame(AssetFs embed.FS) *Game {
	return &Game{
		backgroundColor: color.RGBA{0, 181, 226, 255},
		playerSprite:    pkg.LoadImage("assets/player.png", AssetFs),
	}
}

// func (g *Game) reset() {
// 	g.playerX = screenW/2 - playerW/2
// 	g.playerY = screenH/2 - playerH/2
// }

func (g *Game) Update() error {
	switch true {
	case ebiten.IsKeyPressed(ebiten.KeyW):
		g.playerY -= 10
	case ebiten.IsKeyPressed(ebiten.KeyA):
		g.playerX -= 10
	case ebiten.IsKeyPressed(ebiten.KeyS):
		g.playerY += 10
	case ebiten.IsKeyPressed(ebiten.KeyD):
		g.playerX += 10
	}

	// if g.playerY < 0 || g.playerY >= screenH-playerH {
	// 	g.reset()
	// }
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.backgroundColor)

	X, Y := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %v, Y: %v, FPS: %v", X, Y, ebiten.ActualFPS()))

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.playerX, g.playerY)
	screen.DrawImage(g.playerSprite, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenW, screenH
}
