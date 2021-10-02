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
)

type Game struct {
	images         map[int]*ebiten.Image // 图片资源
	audios         map[int]*audio.Player // 音效
	audioContext   *audio.Context        // 音效器
	singlePosition *Position             // 棋局单例
	bFilpped       bool                  // 是否翻转棋盘
	mvLast         int                   // 上一步棋
	sqSelected     int                   // 需中的格子的值
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBoard(screen)
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
	// fmt.Printf("p.ucpcSquares: %#v", game.singlePosition.ucpcSquares)
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

// drawPiece 绘制棋子
func (g *Game) drawPiece(x, y int, screen, img *ebiten.Image) {
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}

// drawBoard 绘制棋盘，并且加载棋子的位置
func (g *Game) drawBoard(screen *ebiten.Image) {
	// 棋盘
	if v, ok := g.images[ImgChessBoard]; ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		screen.DrawImage(v, op)
	}
	// 棋子
	for x := Left; x <= Right; x++ {
		for y := Top; y <= Bottom; y++ {
			xPos, yPos := 0, 0
			if g.bFilpped {
				xPos = BoardEdge + (xFlip(x)-Left)*SquareSize
				yPos = BoardEdge + (yFlip(y)-Top)*SquareSize
			} else {
				xPos = BoardEdge + (x-Left)*SquareSize
				yPos = BoardEdge + (y-Top)*SquareSize
			}
			// sq -> [0, 255]
			sq := squareXY(x, y)
			pc := g.singlePosition.ucpcSquares[sq]
			if pc != 0 {
				g.drawPiece(xPos, yPos, screen, g.images[pc])
			}
			if sq == g.sqSelected || sq == src(g.mvLast) || sq == dst(g.mvLast) {
				g.drawPiece(xPos, yPos, screen, g.images[ImgSelect])
			}
		}
	}
}
