package entity

import "github.com/garinyaroslav/duel/pkg"

type Projectile struct {
	Owner    PlayerID
	Position pkg.Vec2
	Velocity pkg.Vec2
	Life     float64
}
