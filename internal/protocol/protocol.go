package protocol

import (
	"bytes"
	"encoding/binary"

	"github.com/garinyaroslav/duel/pkg"
)

const (
	MsgInput    byte = 1
	MsgSnapshot byte = 2
)

var order = binary.LittleEndian

type PlayerState struct {
	ID    uint8
	X, Y  float32
	Angle float32
	HP    uint8
}

type BulletState struct {
	Owner uint8
	X, Y  float32
}

type Snapshot struct {
	Tick    uint32
	Phase   uint8
	Timer   float32
	Players []PlayerState
	Bullets []BulletState
}

func EncodeInput(in pkg.Input) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(MsgInput)
	_ = binary.Write(buf, order, in)
	return buf.Bytes()
}

func DecodeInput(b []byte) (pkg.Input, bool) {
	var in pkg.Input
	if len(b) < 1 || b[0] != MsgInput {
		return in, false
	}
	if err := binary.Read(bytes.NewReader(b[1:]), order, &in); err != nil {
		return in, false
	}
	return in, true
}

func EncodeSnapshot(s Snapshot) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(MsgSnapshot)
	_ = binary.Write(buf, order, s.Tick)
	_ = binary.Write(buf, order, s.Phase)
	_ = binary.Write(buf, order, s.Timer)
	_ = binary.Write(buf, order, uint8(len(s.Players)))
	_ = binary.Write(buf, order, s.Players)
	_ = binary.Write(buf, order, uint16(len(s.Bullets)))
	_ = binary.Write(buf, order, s.Bullets)
	return buf.Bytes()
}

func DecodeSnapshot(b []byte) (Snapshot, bool) {
	var s Snapshot
	if len(b) < 1 || b[0] != MsgSnapshot {
		return s, false
	}

	r := bytes.NewReader(b[1:])
	read := func(v any) bool { return binary.Read(r, order, v) == nil }

	var np uint8
	var nb uint16

	if !read(&s.Tick) || !read(&s.Phase) || !read(&s.Timer) || !read(&np) {
		return s, false
	}
	s.Players = make([]PlayerState, np)
	if !read(s.Players) || !read(&nb) {
		return s, false
	}
	s.Bullets = make([]BulletState, nb)
	if !read(s.Bullets) {
		return s, false
	}

	return s, true
}
