package entity

import "github.com/garinyaroslav/duel/pkg"

type Rect struct{ X, Y, W, H float64 }

var Walls = buildWalls()

func buildWalls() []Rect {
	quadrant := []Rect{
		{180, 140, 300, 50},
		{180, 140, 50, 220},
		{420, 330, 60, 140},
		{560, 420, 60, 60},
		{680, 200, 60, 200},
	}

	center := []Rect{
		{820, 100, 280, 50},
		{900, 260, 120, 50},
	}

	var res []Rect
	for _, r := range quadrant {
		res = append(res,
			r,
			Rect{WorldW - r.X - r.W, r.Y, r.W, r.H},
			Rect{r.X, WorldH - r.Y - r.H, r.W, r.H},
			Rect{WorldW - r.X - r.W, WorldH - r.Y - r.H, r.W, r.H},
		)
	}
	for _, r := range center {
		res = append(res, r, Rect{r.X, WorldH - r.Y - r.H, r.W, r.H})
	}

	res = append(res, Rect{900, 460, 120, 160})

	return res
}

func (r Rect) closest(p pkg.Vec2) pkg.Vec2 {
	return pkg.Vec2{
		X: max(r.X, min(p.X, r.X+r.W)),
		Y: max(r.Y, min(p.Y, r.Y+r.H)),
	}
}

func (r Rect) HitsCircle(c pkg.Vec2, radius float64) bool {
	d := c.Sub(r.closest(c))
	return d.X*d.X+d.Y*d.Y < radius*radius
}

func hitsAnyWall(c pkg.Vec2, radius float64) bool {
	for _, r := range Walls {
		if r.HitsCircle(c, radius) {
			return true
		}
	}
	return false
}

func (p *Player) collideWalls() {
	for _, r := range Walls {
		cp := r.closest(p.Position)
		d := p.Position.Sub(cp)
		dist := d.Len()

		if dist == 0 || dist >= PlayerRadius {
			continue
		}

		n := d.Mul(1 / dist)
		p.Position = cp.Add(n.Mul(PlayerRadius))

		if vn := p.Velocity.X*n.X + p.Velocity.Y*n.Y; vn < 0 {
			p.Velocity = p.Velocity.Sub(n.Mul(vn))
		}
	}
}
