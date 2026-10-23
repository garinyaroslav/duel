package net

import (
	"errors"
	"log/slog"
	"net"
	"time"

	"github.com/garinyaroslav/duel/internal/entity"
	"github.com/garinyaroslav/duel/internal/protocol"
	"github.com/garinyaroslav/duel/pkg"
)

const clientTimeout = 3 * time.Second

type clientPacket struct {
	addr  *net.UDPAddr
	input pkg.Input
}

type Host struct {
	World *entity.World

	conn       *net.UDPConn
	inputCh    chan clientPacket
	clientAddr *net.UDPAddr
	lastInput  pkg.Input
	lastRecv   time.Time
	tick       uint32
}

func NewHost() *Host {
	return &Host{
		World:   entity.NewWorld(),
		inputCh: make(chan clientPacket, 64),
	}
}

func (h *Host) Start(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	h.conn = conn
	h.World.AddPlayer(entity.HostID)

	go h.readLoop()
	return nil
}

func (h *Host) readLoop() {
	buf := make([]byte, 1024)

	for {
		n, addr, err := h.conn.ReadFromUDP(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}

		in, ok := protocol.DecodeInput(buf[:n])
		if !ok {
			continue
		}

		select {
		case h.inputCh <- clientPacket{addr: addr, input: in}:
		default:
		}
	}
}

func (h *Host) accept(pkt clientPacket) {
	if h.clientAddr == nil {
		h.clientAddr = pkt.addr
		h.lastInput = pkg.Input{}
		h.World.AddPlayer(entity.ClientID)
		h.World.StartMatch()
		slog.Info("Client joined", "addr", pkt.addr)
	} else if pkt.addr.String() != h.clientAddr.String() {
		return
	}

	h.lastRecv = time.Now()

	if pkt.input.Tick >= h.lastInput.Tick {
		h.lastInput = pkt.input
	}
}

func (h *Host) drainInputs() {
	for {
		select {
		case pkt := <-h.inputCh:
			h.accept(pkt)
		default:
			return
		}
	}
}

func (h *Host) Step(myInput pkg.Input, dt float64) {
	h.drainInputs()

	if h.clientAddr != nil && time.Since(h.lastRecv) > clientTimeout {
		slog.Info("Client left", "addr", h.clientAddr)
		h.clientAddr = nil
		h.lastInput = pkg.Input{}
		h.World.RemovePlayer(entity.ClientID)
		delete(h.World.Players, entity.ClientID)
	}

	inputs := map[entity.PlayerID]pkg.Input{entity.HostID: myInput}
	if h.clientAddr != nil {
		inputs[entity.ClientID] = h.lastInput
	}

	h.World.Step(inputs, dt)
	h.tick++

	if h.clientAddr != nil {
		_, _ = h.conn.WriteToUDP(protocol.EncodeSnapshot(h.snapshot()), h.clientAddr)
	}
}

func (h *Host) snapshot() protocol.Snapshot {
	s := protocol.Snapshot{
		Tick:  h.tick,
		Phase: uint8(h.World.Phase),
		Timer: float32(h.World.Timer),
	}

	for id, p := range h.World.Players {
		s.Players = append(s.Players, protocol.PlayerState{
			ID:    uint8(id),
			X:     float32(p.Position.X),
			Y:     float32(p.Position.Y),
			Angle: float32(p.Angle),
			HP:    uint8(p.HP),
		})
	}

	for _, b := range h.World.Bullets {
		s.Bullets = append(s.Bullets, protocol.BulletState{
			Owner: uint8(b.Owner),
			X:     float32(b.Position.X),
			Y:     float32(b.Position.Y),
		})
	}
	return s
}

func (h *Host) WaitForClient() {
	for h.clientAddr == nil {
		h.accept(<-h.inputCh)
	}
}
