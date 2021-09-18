package chinesschess

import (
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	img, _, err := ebitenutil.NewImageFromFile("./res/ChessBoard.png")
	if err != nil {
		log.Fatal("load ChessBoard fail, err:", err)
	}
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, op)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 520, 576
}

// NewGame 创建象棋程序
func NewGame() bool {
	game := &Game{}
	ebiten.SetWindowSize(520, 576)
	ebiten.SetWindowTitle("中国象棋")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
