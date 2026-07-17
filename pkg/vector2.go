package pkg

import "math"

type Vec2 struct {
	X, Y float64
}

func (v Vec2) Add(other Vec2) Vec2     { return Vec2{v.X + other.X, v.Y + other.Y} }
func (v Vec2) Sub(other Vec2) Vec2     { return Vec2{v.X - other.X, v.Y - other.Y} }
func (v Vec2) Mul(scalar float64) Vec2 { return Vec2{v.X * scalar, v.Y * scalar} }
func (v Vec2) Len() float64            { return math.Sqrt(v.X*v.X + v.Y*v.Y) }

func (v Vec2) Normalized() Vec2 {
	len := v.Len()
	if len == 0 {
		return Vec2{0, 0}
	}
	return v.Mul(1 / len)
}

func LerpVec(a, b Vec2, t float64) Vec2 {
	if t > 1 {
		t = 1
	}
	return a.Add(b.Sub(a).Mul(t))
}
