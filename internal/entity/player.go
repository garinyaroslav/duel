package entity

import (
	"embed"
	"time"

	"github.com/garinyaroslav/duel/pkg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type CharacterType int

const (
	SelfPlayer CharacterType = iota
	Enemy
)

type Player struct {
	Position         pkg.Vec2
	Velocity         pkg.Vec2
	MaxSpeed         float64
	Acceleration     float64
	Deceleration     float64
	Sprite           *ebiten.Image
	ProjectileSprite *ebiten.Image
	ProjectileX      float64
	ProjectileY      float64
	Projectiles      []*Projectile
	characterType    CharacterType

	lastTime time.Time
	lastShot time.Time
}

func NewPlayer(x, y float64, characterType CharacterType, assetFs *embed.FS) *Player {
	return &Player{
		characterType:    characterType,
		Position:         pkg.Vec2{x, y},
		Velocity:         pkg.Vec2{},
		MaxSpeed:         600,
		Acceleration:     100,
		Deceleration:     10,
		Sprite:           pkg.LoadImage("assets/player.png", assetFs),
		ProjectileSprite: pkg.LoadImage("assets/projectile.png", assetFs),
		Projectiles:      make([]*Projectile, 0, 50),
		lastTime:         time.Now(),
	}
}

func (p *Player) Update() {
	currentTime := time.Now()

	dt := currentTime.Sub(p.lastTime).Seconds()
	p.lastTime = currentTime

	if dt > 0.1 {
		dt = 0.1
	}

	updateMovement(p, dt)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.shoot()
	}

	for i := range p.Projectiles {
		if !p.Projectiles[i].Active {
			continue
		}
		p.Projectiles[i].Position = p.Projectiles[i].Position.Add(p.Projectiles[i].Velocity.Mul(1 * dt))
	}
}

func (p *Player) shoot() {
	const shotCooldown = 150 * time.Millisecond

	if time.Since(p.lastShot) < shotCooldown {
		return
	}

	cursorX, cursorY := ebiten.CursorPosition()

	dir := pkg.Vec2{float64(cursorX), float64(cursorY)}.Sub(pkg.Vec2{p.Position.X, p.Position.Y})

	if dir.Len() == 0 {
		return
	}

	dir = dir.Normalized()

	bulletSpeed := 800.0

	p.Projectiles = append(p.Projectiles, &Projectile{
		Position: pkg.Vec2{
			X: p.Position.X,
			Y: p.Position.Y,
		},
		Velocity: dir.Mul(bulletSpeed),
		Active:   true,
	})
	p.lastShot = time.Now()
}

func updateMovement(p *Player, dt float64) {
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
