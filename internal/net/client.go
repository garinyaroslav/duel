package net

import (
	"errors"
	"net"
	"time"

	"github.com/garinyaroslav/duel/internal/entity"
	"github.com/garinyaroslav/duel/internal/protocol"
	"github.com/garinyaroslav/duel/pkg"
)

const smoothing = 0.35

type Client struct {
	World *entity.World

	conn       *net.UDPConn
	snapshotCh chan protocol.Snapshot
	tick       uint32

	hasSnap  bool
	lastSnap uint32
	targets  map[entity.PlayerID]pkg.Vec2
}

func NewClient() *Client {
	return &Client{
		World:      entity.NewWorld(),
		snapshotCh: make(chan protocol.Snapshot, 8),
		targets:    make(map[entity.PlayerID]pkg.Vec2),
	}
}

func (c *Client) Start(serverAddr string) error {
	raddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	c.conn = conn
	go c.readLoop()
	return nil
}

func (c *Client) readLoop() {
	buf := make([]byte, 4096)

	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			time.Sleep(10 * time.Millisecond)
			continue
		}

		snap, ok := protocol.DecodeSnapshot(buf[:n])
		if !ok {
			continue
		}

		select {
		case c.snapshotCh <- snap:
		default:
		}
	}
}

func (c *Client) Step(myInput pkg.Input, dt float64) {
	myInput.Tick = c.tick
	c.tick++

	_, _ = c.conn.Write(protocol.EncodeInput(myInput))

	for drained := false; !drained; {
		select {
		case snap := <-c.snapshotCh:
			c.apply(snap)
		default:
			drained = true
		}
	}

	c.smooth()
}

func (c *Client) apply(s protocol.Snapshot) {
	if c.hasSnap && s.Tick <= c.lastSnap {
		return
	}
	c.hasSnap = true
	c.lastSnap = s.Tick

	newPhase := entity.Phase(s.Phase)
	teleport := c.World.Phase == entity.PhaseRoundOver && newPhase != entity.PhaseRoundOver
	c.World.Phase = newPhase
	c.World.Timer = float64(s.Timer)

	seen := make(map[entity.PlayerID]bool, len(s.Players))

	for _, ps := range s.Players {
		id := entity.PlayerID(ps.ID)
		seen[id] = true

		target := pkg.Vec2{X: float64(ps.X), Y: float64(ps.Y)}
		p, ok := c.World.Players[id]
		if !ok {
			p = c.World.AddPlayer(id)
		}
		if !ok || teleport {
			p.Position = target
		}
		p.Angle = float64(ps.Angle)
		p.HP = int(ps.HP)
		c.targets[id] = target
	}

	for id := range c.World.Players {
		if !seen[id] {
			delete(c.World.Players, id)
			delete(c.targets, id)
		}
	}

	c.World.Bullets = c.World.Bullets[:0]
	for _, b := range s.Bullets {
		c.World.Bullets = append(c.World.Bullets, &entity.Projectile{
			Owner:    entity.PlayerID(b.Owner),
			Position: pkg.Vec2{X: float64(b.X), Y: float64(b.Y)},
		})
	}
}

func (c *Client) smooth() {
	for id, p := range c.World.Players {
		if t, ok := c.targets[id]; ok {
			p.Position = pkg.LerpVec(p.Position, t, smoothing)
		}
	}
}
