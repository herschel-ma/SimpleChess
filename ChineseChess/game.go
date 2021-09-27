package ChineseChess

import (
	"bytes"
	"image"
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/herschel-ma/SimpleChess/res"
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	imageFile, _, err := image.Decode(bytes.NewReader(res.ImgChessBoard))
	if err != nil {
		log.Fatal("load ChessBoard fail, err:", err)
	}
	img := ebiten.NewImageFromImage(imageFile)
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
