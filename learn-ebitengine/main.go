package main

import (
	"embed"
	"image"
	_ "image/png" // used for decoding
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

var PlayerSprite = mustLoadImage("assets/PNG/playerShip2_red.png")

type Game struct{}

func (g *Game) Update() error {
	return nil
}

// The screen argument is the image displayed every frame (after the method returns). Our job is
// to draw other images (or text) on it.
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// 1. Translate()
	// op.GeoM.Translate(150, 200)

	// 2. Rotate() function expects the angle in radians, not degrees.
	// Radians measure angle using the radius of a circle. One full circle = 2π radians (~6.28).
	//
	// Since 360° = 2π radians, you can simplify that to:
	// 180° = π radians
	// So to convert any degree to radians:
	// radians = degrees * π / 180
	// op.GeoM.Rotate(45.0 * math.Pi / 180.0)

	// 3. Scale() makes the image bigger horizontally and vertically.
	// Passing negative values flips the image horizontally or vertically
	// op.GeoM.Scale(2, 2)

	// You can combine all of them (the op.GeoM) in a single DrawImage call.
	// There’s a caveat, though: the order matters. If you use Translate first
	// and then apply the rotation or scale, you may not see what you expect
	// because the change will be applied to the new position.
	//
	// To illustrate this, let’s see how to rotate the image around its center (a
	// common use case).
	//
	// You need to move the pivot point to the image’s center.
	// You can find it simply by dividing the image’s width and height by half.
	// (Calculating stuff based on the image size is also super common.)
	width := PlayerSprite.Bounds().Dx()
	height := PlayerSprite.Bounds().Dy()

	halfW := float64(width / 2)
	halfH := float64(height / 2)

	// You need to move the image by the negative values first. This is so
	// the image’s center aligns with the origin point (0, 0). Then, apply
	// the rotation and move the image “back” by the same amount.
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(45.0 * math.Pi / 180.0)
	op.GeoM.Translate(halfW, halfH)

	screen.DrawImage(PlayerSprite, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	g := &Game{}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}
