package main

import (
	"embed"
	"errors"
	"image"
	"log"
	"math"

	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/imgs/*
var assets embed.FS

var (
	whiteImage = ebiten.NewImage(3, 3)
)

func init() {
	whiteImage.Fill(color.White)
}

type Game struct {
	// Shinji *ebiten.Image
	Hex Hex
}

type Hex struct {
	X, Y, Radius      float64
	Selected, Hovered bool
}

type Color struct {
	R, G, B float32
}

func Black() Color {
	return Color{0, 0, 0}
}

func Red() Color {
	return Color{1, 0, 0}
}

func (hex Hex) vertices(color Color) []ebiten.Vertex {
	const points = 6
	verts := []ebiten.Vertex{}

	for i := 0; i < points; i++ {
		rate := float64(i) / float64(points)
		verts = append(verts, ebiten.Vertex{
			DstX:   float32(hex.Radius*math.Cos(2*math.Pi*rate) + hex.X),
			DstY:   float32(hex.Radius*math.Sin(2*math.Pi*rate) + hex.Y),
			ColorR: color.R,
			ColorG: color.G,
			ColorB: color.B,
			ColorA: 1,
		})
	}

	verts = append(verts, ebiten.Vertex{
		DstX:   float32(hex.X),
		DstY:   float32(hex.Y),
		ColorR: color.R,
		ColorG: color.G,
		ColorB: color.B,
		ColorA: 1,
	})

	return verts
}

func (hex Hex) indecies() []uint16 {
	const points = 6

	indices := []uint16{}
	for i := 0; i < points; i++ {
		indices = append(indices, uint16(i), uint16(i+1)%uint16(points), uint16(points))
	}

	return indices
}

func (hex Hex) overlayHex() Hex {
	hex.Radius -= 2
	return hex
}

func (hex Hex) Draw(screen *ebiten.Image) {
	hex.draw(screen, Red())
	if !hex.Selected {
		hex.overlayHex().draw(screen, Black())
	}
}

func (hex Hex) draw(screen *ebiten.Image, color Color) {
	op := &ebiten.DrawTrianglesOptions{}
	op.Address = ebiten.AddressUnsafe
	screen.DrawTriangles(hex.vertices(color), hex.indecies(), whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("Escape")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Hex.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 800
}

func main() {
	// shinjiBytes, err := assets.ReadFile("assets/imgs/shinji.png")
	// if err != nil {
	// 	panic(err)
	// }

	// shinji, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(shinjiBytes))
	// if err != nil {
	// 	panic(err)
	// }

	ebiten.SetWindowTitle("Eva CAPTCHA")
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&Game{Hex: Hex{
		X: 640, Y: 480, Radius: 50,
	}}); err != nil {
		log.Fatal(err)
	}
}
