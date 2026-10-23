package entity

import (
	"math"

	"github.com/garinyaroslav/duel/pkg"
)

type PlayerID uint8

const (
	HostID PlayerID = iota
	ClientID
)

const (
	PlayerRadius = 40.0
	BulletRadius = 10.0
	MaxHP        = 2

	MaxSpeed     = 600.0
	Acceleration = 10.0
	Deceleration = 8.0

	ShotCooldown = 0.15
	BulletSpeed  = 900.0
	BulletLife   = 1.5
	MuzzleOffset = 50.0
)

type Player struct {
	ID       PlayerID
	Position pkg.Vec2
	Velocity pkg.Vec2
	Angle    float64
	HP       int

	cooldown float64
}

func (p *Player) Alive() bool { return p.HP > 0 }

func (p *Player) Reset(pos pkg.Vec2, angle float64) {
	p.Position = pos
	p.Velocity = pkg.Vec2{}
	p.Angle = angle
	p.HP = MaxHP
	p.cooldown = 0
}

func (p *Player) Update(in pkg.Input, dt float64) *Projectile {
	p.move(in, dt)

	aim := pkg.Vec2{X: float64(in.AimX), Y: float64(in.AimY)}.Sub(p.Position)
	if aim.Len() > 1 {
		p.Angle = math.Atan2(aim.Y, aim.X)
	}

	p.cooldown -= dt
	if in.Shoot && p.cooldown <= 0 && aim.Len() > 1 {
		p.cooldown = ShotCooldown
		dir := aim.Normalized()
		return &Projectile{
			Owner:    p.ID,
			Position: p.Position.Add(dir.Mul(MuzzleOffset)),
			Velocity: dir.Mul(BulletSpeed),
			Life:     BulletLife,
		}
	}
	return nil
}

func (p *Player) move(in pkg.Input, dt float64) {
	dir := pkg.Vec2{X: float64(in.MoveX), Y: float64(in.MoveY)}
	accel := Deceleration
	if dir.Len() > 0 {
		dir = dir.Normalized()
		accel = Acceleration
	}

	p.Velocity = pkg.LerpVec(p.Velocity, dir.Mul(MaxSpeed), accel*dt)
	p.Position = p.Position.Add(p.Velocity.Mul(dt))

	p.Position.X = max(PlayerRadius, min(p.Position.X, WorldW-PlayerRadius))
	p.Position.Y = max(PlayerRadius, min(p.Position.Y, WorldH-PlayerRadius))
}
