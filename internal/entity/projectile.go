package entity

import "github.com/garinyaroslav/duel/pkg"

type Projectile struct {
	Position pkg.Vec2
	Velocity pkg.Vec2
	Active   bool
}
