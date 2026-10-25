package internal

import (
	"image/color"

	"github.com/garinyaroslav/duel/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	floorColor  = color.RGBA{38, 50, 56, 255}
	gridColor   = color.RGBA{46, 60, 68, 255}
	wallColor   = color.RGBA{96, 108, 122, 255}
	wallEdge    = color.RGBA{150, 164, 180, 255}
	wallShadow  = color.RGBA{0, 0, 0, 90}
	borderColor = color.RGBA{150, 164, 180, 255}
)

func renderMap() *ebiten.Image {
	img := ebiten.NewImage(SCREEN_W, SCREEN_H)
	img.Fill(floorColor)

	const cell = 60
	for x := 0; x <= SCREEN_W; x += cell {
		vector.StrokeLine(img, float32(x), 0, float32(x), SCREEN_H, 1, gridColor, false)
	}
	for y := 0; y <= SCREEN_H; y += cell {
		vector.StrokeLine(img, 0, float32(y), SCREEN_W, float32(y), 1, gridColor, false)
	}

	for _, r := range entity.Walls {
		x, y, w, h := float32(r.X), float32(r.Y), float32(r.W), float32(r.H)

		vector.DrawFilledRect(img, x+8, y+8, w, h, wallShadow, false) // тень
		vector.DrawFilledRect(img, x, y, w, h, wallColor, false)
		vector.StrokeRect(img, x, y, w, h, 4, wallEdge, false)
	}

	vector.StrokeRect(img, 2, 2, SCREEN_W-4, SCREEN_H-4, 4, borderColor, false)

	return img
}
