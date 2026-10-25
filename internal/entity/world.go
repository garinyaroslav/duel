package entity

import (
	"math"

	"github.com/garinyaroslav/duel/pkg"
)

const (
	WorldW = 1920.0
	WorldH = 1080.0

	CountdownTime = 3.0
	RespawnTime   = 3.0
)

type Phase uint8

const (
	PhaseWaiting Phase = iota
	PhaseCountdown
	PhasePlaying
	PhaseRoundOver
)

type World struct {
	Tick    uint64
	Phase   Phase
	Timer   float64
	Players map[PlayerID]*Player
	Bullets []*Projectile
}

func NewWorld() *World {
	return &World{
		Players: make(map[PlayerID]*Player),
		Bullets: make([]*Projectile, 0, 64),
	}
}

func spawnPoint(id PlayerID) (pkg.Vec2, float64) {
	if id == ClientID {
		return pkg.Vec2{X: WorldW - 300, Y: WorldH / 2}, math.Pi
	}
	return pkg.Vec2{X: 300, Y: WorldH / 2}, 0
}

func (w *World) AddPlayer(id PlayerID) *Player {
	p := &Player{ID: id}
	p.Reset(spawnPoint(id))
	w.Players[id] = p
	return p
}

func (w *World) RemovePlayer(id PlayerID) {
	delete(w.Players, id)
	w.Phase = PhaseWaiting
	w.resetRound()
}

func (w *World) StartMatch() {
	w.resetRound()
	w.Phase = PhaseCountdown
	w.Timer = CountdownTime
}

func (w *World) resetRound() {
	for id, p := range w.Players {
		p.Reset(spawnPoint(id))
	}
	w.Bullets = w.Bullets[:0]
}

func (w *World) Step(inputs map[PlayerID]pkg.Input, dt float64) {
	switch w.Phase {
	case PhaseCountdown:
		w.Timer -= dt
		if w.Timer <= 0 {
			w.Phase = PhasePlaying
		}

	case PhasePlaying:
		w.updatePlayers(inputs, dt, true)
		w.moveBullets(dt)
		w.resolveHits()

		if w.anyDead() {
			w.Phase = PhaseRoundOver
			w.Timer = RespawnTime
		}

	case PhaseRoundOver:
		w.updatePlayers(inputs, dt, false)
		w.moveBullets(dt)

		w.Timer -= dt
		if w.Timer <= 0 {
			w.resetRound()
			w.Phase = PhasePlaying
		}
	}

	w.Tick++
}

func (w *World) updatePlayers(inputs map[PlayerID]pkg.Input, dt float64, canShoot bool) {
	for id, in := range inputs {
		p, ok := w.Players[id]
		if !ok || !p.Alive() {
			continue
		}
		if !canShoot {
			in.Shoot = false
		}
		if b := p.Update(in, dt); b != nil {
			w.Bullets = append(w.Bullets, b)
		}
	}
}

func (w *World) moveBullets(dt float64) {
	alive := w.Bullets[:0]
	for _, b := range w.Bullets {
		b.Position = b.Position.Add(b.Velocity.Mul(dt))
		b.Life -= dt

		inside := b.Position.X >= 0 && b.Position.X <= WorldW &&
			b.Position.Y >= 0 && b.Position.Y <= WorldH
		if b.Life > 0 && inside && !hitsAnyWall(b.Position, BulletRadius) {
			alive = append(alive, b)
		}
	}
	w.Bullets = alive
}

func (w *World) resolveHits() {
	alive := w.Bullets[:0]
	for _, b := range w.Bullets {
		hit := false
		for id, p := range w.Players {
			if id == b.Owner || !p.Alive() {
				continue
			}
			if p.Position.Sub(b.Position).Len() < PlayerRadius+BulletRadius {
				p.HP = max(p.HP-1, 0)
				hit = true
				break
			}
		}
		if !hit {
			alive = append(alive, b)
		}
	}
	w.Bullets = alive
}

func (w *World) anyDead() bool {
	for _, p := range w.Players {
		if !p.Alive() {
			return true
		}
	}
	return false
}
