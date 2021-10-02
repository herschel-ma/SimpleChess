package ChineseChess

import (
	"bytes"
	"fmt"
	"image"
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/herschel-ma/SimpleChess/res"
)

type Game struct {
	images         map[int]*ebiten.Image // 图片资源
	audios         map[int]*audio.Player // 音效
	audioContext   *audio.Context        // 音效器
	singlePosition *Position             // 棋局单例
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
	return BoardWidth, BoardHeight
}

// NewGame 创建象棋程序
func NewGame() bool {
	game := &Game{
		images:         make(map[int]*ebiten.Image),
		audios:         make(map[int]*audio.Player),
		singlePosition: NewPosition(),
	}
	if game == nil || game.singlePosition == nil {
		return false
	}
	// 初始化 audioContext
	game.audioContext = audio.NewContext(sampleRate)
	// 加载资源
	if ok := game.LoadResource(); !ok {
		return false
	}
	// 设置窗口大小
	ebiten.SetWindowSize(BoardWidth, BoardHeight)
	// 设置窗口标题
	ebiten.SetWindowTitle("中国象棋")
	// 后端初始化棋盘
	game.singlePosition.startup()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// LoadResource 加载资源
func (g *Game) LoadResource() bool {
	for k, v := range resMap {
		if k >= MusicSelect {
			// 加载音效
			stream, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(v))
			if err != nil {
				fmt.Println(err)
				return false
			}
			player, err := audio.NewPlayer(g.audioContext, stream)
			if err != nil {
				fmt.Println(err)
				return false
			}
			g.audios[k] = player
		} else {
			// 加载图片
			imgFile, _, err := image.Decode(bytes.NewReader(v))
			if err != nil {
				fmt.Println(err)
				return false
			}
			ebitenImage := ebiten.NewImageFromImage(imgFile)
			g.images[k] = ebitenImage
		}
	}
	return true
}

// playAudio 播放音效
func (g *Game) playAudio(key int) bool {
	if player, ok := g.audios[key]; ok {
		player.Rewind()
		player.Play()
		return true
	} else {
		return false
	}
}
