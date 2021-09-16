package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "hello world")
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 380, 160
}
func main() {
	ebiten.SetWindowSize(400, 200)
	ebiten.SetWindowTitle("hello world")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
