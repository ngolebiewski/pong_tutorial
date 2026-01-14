package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png" // Essential for decoding PNGs in WASM/Go
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed images/gopher.png
	gopher_png []byte
)

// goBall is our global reference to the Gopher's image and dimensions.
// It is accessible by other files in 'package main'.
var goBall *Gopher

type Gopher struct {
	width  int
	height int
	img    *ebiten.Image
}

func init() {
	// 1. Decode the embedded bytes into a standard image.Image
	img, _, err := image.Decode(bytes.NewReader(gopher_png))
	if err != nil {
		log.Fatal("Error decoding gopher.png: ", err)
	}

	// 2. Convert standard image to ebiten.Image for GPU rendering
	ebitenImg := ebiten.NewImageFromImage(img)

	// 3. Capture dimensions so the physics engine knows how big the "ball" is
	w, h := ebitenImg.Bounds().Dx(), ebitenImg.Bounds().Dy()

	// 4. Initialize the global pointer
	goBall = &Gopher{
		width:  w,
		height: h,
		img:    ebitenImg,
	}
}
