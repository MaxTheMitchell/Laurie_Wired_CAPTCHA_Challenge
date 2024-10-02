package main

import (
	"bytes"
	"embed"
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed shinji.png
var assets embed.FS

type Game struct {
	Shinji *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(500, 400)
	screen.DrawImage(g.Shinji, &op)
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 800
}

func main() {
	shinjiBytes, err := assets.ReadFile("shinji.png")
	if err != nil {
		panic(err)
	}

	shinji, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(shinjiBytes))
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{Shinji: shinji}); err != nil {
		log.Fatal(err)
	}
}
