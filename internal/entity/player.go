package entity

import (
	"time"

	"github.com/garinyaroslav/duel/pkg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	Position     pkg.Vec2
	Velocity     pkg.Vec2
	MaxSpeed     float64
	Acceleration float64
	Deceleration float64
	Sprite       *ebiten.Image

	lastTime time.Time
}

func NewPlayer(x, y float64, sprite *ebiten.Image) *Player {
	return &Player{
		Position:     pkg.Vec2{x, y},
		Velocity:     pkg.Vec2{},
		MaxSpeed:     600,
		Acceleration: 100,
		Deceleration: 10,
		Sprite:       sprite,
		lastTime:     time.Now(),
	}
}

func (p *Player) Update() {
	currentTime := time.Now()
	dt := currentTime.Sub(p.lastTime).Seconds()
	p.lastTime = currentTime

	if dt > 0.1 {
		dt = 0.1
	}

	input := pkg.Vec2{0, 0}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		input.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		input.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		input.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		input.X += 1
	}

	if input.Len() > 0 {
		input = input.Normalized()
	}

	desired := input.Mul(p.MaxSpeed)
	accel := p.Acceleration
	if input.Len() == 0 {
		accel = p.Deceleration
	}

	p.Velocity = pkg.LerpVec(p.Velocity, desired, accel*dt)
	p.Position = p.Position.Add(p.Velocity.Mul(dt))
}
