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
	HexManager HexManager
}

type Color struct {
	R, G, B, A float32
}

func Black() Color {
	return Color{R: 0, G: 0, B: 0, A: 1}
}

func OffRed() Color {
	return Color{R: 1, G: .2, B: 0, A: 1}
}

func Red() Color {
	return Color{R: 1, G: 0, B: 0, A: 1}
}

func Orange() Color {
	return Color{R: 0.93, G: 0.5, B: 0.06, A: 1}
}

type HexManager struct {
	Hexes []Hex
}

func (hexManager HexManager) NeedsInit() bool {
	return len(hexManager.Hexes) == 0
}

func (hexManager *HexManager) InitHexes(screenWidth, screenHeight int) {

	rows := [][]int{
		{-1, 0, 1},
		{-2, -1, 1, 2},
		{-2, -1, 0, 1, 2},
		{-2, -1, 1, 2},
		{-1, 0, 1},
	}

	for y, row := range rows {
		for _, x := range row {
			hex := Hex{
				X:      float64(screenWidth) / 2,
				Y:      float64(screenHeight) / 2,
				Radius: (float64(screenHeight) / 10) * 0.8,
			}

			hex.X += hex.Radius * math.Sqrt(3) * float64(x)
			hex.Y += hex.Radius*2*float64(y-2) - (float64(y) * hex.Radius * 0.5)
			if y%2 == 1 {
				hex.X -= ((hex.Radius * math.Sqrt(3)) / 2) * (float64(x) / math.Abs(float64(x)))
			}

			hexManager.Hexes = append(hexManager.Hexes, hex)
		}
	}
}

func (hexManager HexManager) Update() {
	x, y := ebiten.CursorPosition()
	for i, hex := range hexManager.Hexes {
		hex.Hovered = false
		if hex.IsInside(float64(x), float64(y)) {
			hex.Hovered = true
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
				hex.Selected = !hex.Selected
			}
		}
		hexManager.Hexes[i] = hex
	}
}

func (hexManger HexManager) Draw(screen *ebiten.Image) {
	for _, hex := range hexManger.Hexes {
		hex.Draw(screen)
	}
}

type Hex struct {
	X, Y, Radius      float64
	Selected, Hovered bool
}

func (hex Hex) IsInside(x, y float64) bool {
	q2x := math.Abs(x - hex.X)
	q2y := math.Abs(y - hex.Y)
	vert := (hex.Radius * math.Sqrt(3)) / 4
	hori := hex.Radius
	if q2x > hori || q2y > vert*2 {
		return false
	}
	return 2*vert*hori-vert*q2x-hori*q2y >= 0
}

func (hex Hex) vertices(color Color) []ebiten.Vertex {
	const points = 6
	verts := []ebiten.Vertex{}

	for i := 0; i < points; i++ {
		rate := (float64(i) + .5) / float64(points)
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
	hex.overlayHex().draw(screen, hex.fillColor())
}

func (hex Hex) fillColor() Color {
	if hex.Selected && hex.Hovered {
		return OffRed()
	}
	if hex.Hovered {
		return Orange()
	}
	if hex.Selected {
		return Red()
	}
	return Black()
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

	g.HexManager.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.HexManager.NeedsInit() {
		g.HexManager.InitHexes(screen.Size())
	}

	g.HexManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
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
	if err := ebiten.RunGame(&Game{HexManager: HexManager{}}); err != nil {
		log.Fatal(err)
	}
}
