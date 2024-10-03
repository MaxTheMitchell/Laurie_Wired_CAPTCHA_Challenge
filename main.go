package main

import (
	"bytes"
	"embed"
	"errors"
	"image"
	"log"
	"math"
	"time"

	"image/color"
	_ "image/gif"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/exp/rand"
)

/*
 Hello this is my somewhat evangelion themed captcha (sorry I have not seen Stiens Gate)

 The goal is to click on the hexes on the big hex grid such that the red hexes match the ones in the tiny hex grid

 compile by running `go build main.go``
 (you may need to install the dependencies listed here https://ebitengine.org/en/documents/install.html)
 (also to cross compile to windows run `GOOS=windows GOARCH=amd64 ^C build main.go`)

 then just run that exe

 Thank you so much for checking out my captcha, had to throw this together before I left on a trip to Europe!

 LOVE THE CHANNEL SO MUCH!!!!! <3
*/

//go:embed assets/imgs/*
var assets embed.FS

var (
	whiteImage = ebiten.NewImage(3, 3)
)

func init() {
	whiteImage.Fill(color.White)
	rand.Seed(uint64(time.Now().UnixNano()))
}

type Game struct {
	HexManager HexManager
	KeyManager HexManager
	Alerts     Alerts
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

func Green() Color {
	return Color{R: 0.208, G: 1, B: 0.694, A: 1}
}

func Orange() Color {
	return Color{R: 0.93, G: 0.5, B: 0.06, A: 1}
}

type Alerts struct {
	Frame int
	Img   *ebiten.Image
}

func (alerts *Alerts) Update() {
	alerts.Frame++
	if alerts.Frame > 40 {
		alerts.Frame = 0
	}
}

func (alerts Alerts) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()
	if alerts.Frame > 20 {
		for y := 0; y < screenHeight; y += alerts.Img.Bounds().Dy() {
			ops := &ebiten.DrawImageOptions{}
			ops.GeoM.Translate(float64(screenWidth)-float64(alerts.Img.Bounds().Dx()), float64(y))
			screen.DrawImage(alerts.Img, ops)
		}

	}
}

type HexManager struct {
	IsKey bool
	Hexes []Hex
}

func (hexManager HexManager) NeedsInit() bool {
	return len(hexManager.Hexes) == 0
}

func (hexManager *HexManager) InitHexes(centerX, centerY, height int) {
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
				X:      float64(centerX),
				Y:      float64(centerY),
				Radius: (float64(height) / 10) * 0.8,
				IsKey:  hexManager.IsKey,
			}

			hex.X += hex.Radius * math.Sqrt(3) * float64(x)
			hex.Y += hex.Radius*2*float64(y-2) - (float64(y) * hex.Radius * 0.5)
			if y%2 == 1 {
				hex.X -= ((hex.Radius * math.Sqrt(3)) / 2) * (float64(x) / math.Abs(float64(x)))
			}

			hexManager.Hexes = append(hexManager.Hexes, hex)
		}
	}

	if hexManager.IsKey {
		hexManager.selectRandom()
	}
}

func (hexManager HexManager) selectRandom() {
	const randomMax = 6

	for i := 0; i < randomMax; i++ {
		hexManager.Hexes[rand.Intn(len(hexManager.Hexes))].Selected = true
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
	X, Y, Radius             float64
	Selected, Hovered, IsKey bool
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
	if hex.IsKey {
		hex.draw(screen, hex.fillColor())
	} else {
		hex.draw(screen, Red())
		hex.overlayHex().draw(screen, hex.fillColor())
	}
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
	if hex.IsKey {
		return Green()
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
		return errors.New("YOU RAN AWAY, IDOIT SHINJI!")
	}

	if !g.HexManager.NeedsInit() && g.passed() {
		return errors.New("You passed the captcha! Congratulations! Congratulations! Congratulations! Congratulations! Congratulations! Congratulations! Congratulations! Congratulations! Congratulations! Congratulations! \n Thank you all :)")
	}

	g.HexManager.Update()
	g.Alerts.Update()

	return nil
}

func (g *Game) passed() bool {
	for i, key := range g.KeyManager.Hexes {
		if key.Selected != g.HexManager.Hexes[i].Selected {
			return false
		}
	}
	return true
}

func (g *Game) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()

	if g.HexManager.NeedsInit() {
		g.HexManager.InitHexes(screenWidth/2, screenHeight/2+100, screenHeight)
		g.KeyManager.InitHexes(screenWidth/5, screenHeight/4, screenHeight/2)
	}

	g.HexManager.Draw(screen)
	g.KeyManager.Draw(screen)
	g.Alerts.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	alertBytes, err := assets.ReadFile("assets/imgs/alert.gif")
	if err != nil {
		panic(err)
	}

	alertImg, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(alertBytes))
	if err != nil {
		panic(err)
	}

	ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)

	ebiten.SetWindowTitle("Eva CAPTCHA")
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&Game{
		HexManager: HexManager{},
		KeyManager: HexManager{IsKey: true},
		Alerts:     Alerts{Img: alertImg}}); err != nil {
		log.Fatal(err)
	}
}
