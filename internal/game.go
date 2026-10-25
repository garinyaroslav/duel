package internal

import (
	"embed"
	"fmt"
	"image/color"
	"math"

	"github.com/garinyaroslav/duel/internal/entity"
	"github.com/garinyaroslav/duel/internal/net"
	"github.com/garinyaroslav/duel/pkg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	SCREEN_W, SCREEN_H = 1920, 1080
	TPS                = 60
	FIXED_DT           = 1.0 / TPS
)

type Game struct {
	background color.RGBA
	world      *entity.World
	step       func(pkg.Input, float64)
	localID    entity.PlayerID
	big        *ebiten.Image
	mapImg     *ebiten.Image

	playerSprite     *ebiten.Image
	projectileSprite *ebiten.Image
	host             *net.Host
}

func NewGame(assetFs *embed.FS, mode, addr string) (*Game, error) {
	g := &Game{
		background:       color.RGBA{0, 181, 226, 255},
		playerSprite:     pkg.LoadImage("assets/player.png", assetFs),
		projectileSprite: pkg.LoadImage("assets/projectile.png", assetFs),
		big:              ebiten.NewImage(256, 16),
		mapImg:           renderMap(),
	}

	switch mode {
	case "host":
		h := net.NewHost()
		if err := h.Start(addr); err != nil {
			return nil, err
		}
		g.host = h
		g.world, g.step, g.localID = h.World, h.Step, entity.HostID
	case "client":
		c := net.NewClient()
		if err := c.Start(addr); err != nil {
			return nil, err
		}
		g.world, g.step, g.localID = c.World, c.Step, entity.ClientID
	default:
		return nil, fmt.Errorf("unknown mode %q", mode)
	}

	return g, nil
}

func (g *Game) Update() error {
	g.step(pkg.ReadLocalInput(), FIXED_DT)
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_W, SCREEN_H
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.mapImg, nil)

	g.drawBullets(screen)
	g.drawPlayers(screen)
	g.drawOverlay(screen)
	drawDebugInfo(screen)
}

func (g *Game) drawPlayers(screen *ebiten.Image) {
	b := g.playerSprite.Bounds()
	w, h := float64(b.Dx()), float64(b.Dy())

	for id, p := range g.world.Players {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Rotate(p.Angle)
		op.GeoM.Translate(p.Position.X, p.Position.Y)

		switch {
		case !p.Alive():
			op.ColorScale.Scale(0.4, 0.4, 0.4, 0.5)
		case id != g.localID:
			op.ColorScale.Scale(1, 0.5, 0.5, 1)
		}
		screen.DrawImage(g.playerSprite, op)

		if p.Alive() {
			ebitenutil.DebugPrintAt(screen,
				fmt.Sprintf("HP %d/%d", p.HP, entity.MaxHP),
				int(p.Position.X)-20, int(p.Position.Y-entity.PlayerRadius)-24)
		}
	}
}

func (g *Game) drawOverlay(screen *ebiten.Image) {
	secs := int(math.Ceil(g.world.Timer))

	switch g.world.Phase {
	case entity.PhaseWaiting:
		g.drawBigText(screen, "Waiting for opponent...", 3, 0)
	case entity.PhaseCountdown:
		g.drawBigText(screen, fmt.Sprint(secs), 14, 0)
	case entity.PhaseRoundOver:
		g.drawBigText(screen, g.roundResult(), 8, -80)
		g.drawBigText(screen, fmt.Sprintf("Respawn in %d", secs), 4, 60)
	}
}

func (g *Game) roundResult() string {
	me := g.world.Players[g.localID]
	if me == nil {
		return ""
	}
	if me.Alive() {
		return "YOU WIN"
	}
	for id, p := range g.world.Players {
		if id != g.localID && p.Alive() {
			return "YOU DIED"
		}
	}
	return "DRAW"
}

func (g *Game) drawBigText(screen *ebiten.Image, s string, scale, dy float64) {
	g.big.Clear()
	ebitenutil.DebugPrint(g.big, s)

	w, h := float64(len(s)*6), 16.0
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(SCREEN_W/2, SCREEN_H/2+dy)
	screen.DrawImage(g.big, op)
}

func (g *Game) drawBullets(screen *ebiten.Image) {
	b := g.projectileSprite.Bounds()
	w, h := float64(b.Dx()), float64(b.Dy())

	for _, p := range g.world.Bullets {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.Position.X-w/2, p.Position.Y-h/2)
		screen.DrawImage(g.projectileSprite, op)
	}
}

func drawDebugInfo(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %v, Y: %v, FPS: %.0f", x, y, ebiten.ActualFPS()))
}

func (g *Game) WaitForOpponent() {
	if g.host != nil {
		g.host.WaitForClient()
	}
}
