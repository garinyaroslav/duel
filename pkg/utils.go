package pkg

import (
	"embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadImage(assetPath string, embedFs *embed.FS) *ebiten.Image {
	f, err := embedFs.Open(assetPath)
	if err != nil {
		log.Panic(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Panic(err)
	}

	return ebiten.NewImageFromImage(img)
}
