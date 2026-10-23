package pkg

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	Tick  uint32
	MoveX float32
	MoveY float32
	AimX  float32
	AimY  float32
	Shoot bool
}

func ReadLocalInput() Input {
	var in Input

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		in.MoveY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		in.MoveY += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		in.MoveX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		in.MoveX += 1
	}

	cx, cy := ebiten.CursorPosition()
	in.AimX = float32(cx)
	in.AimY = float32(cy)
	in.Shoot = inpututil.IsKeyJustPressed(ebiten.KeySpace)

	return in
}
