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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	images         map[int]*ebiten.Image // 图片资源
	audios         map[int]*audio.Player // 音效
	audioContext   *audio.Context        // 音效器
	singlePosition *Position             // 棋局单例
	bFilpped       bool                  // 是否翻转棋盘
	mvLast         int                   // 上一步棋
	sqSelected     int                   // 选中的格子的值
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		// fmt.Printf("点击了:%d, %d", x, y)
		xPos := Left + (x-BoardEdge)/SquareSize
		yPos := Top + (y-BoardEdge)/SquareSize
		g.clickSquare(xPos, yPos)
	}
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
	if player, ok := g.audios[key]; !ok {
		return false
	} else {
		if err := player.Rewind(); err != nil {
			fmt.Println(err)
			return false
		}
		player.Play()
		return true
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

// drawPieceScale 绘制棋子
func (g *Game) drawPieceScale(x, y, scaleX, scaleY float64, screen, img *ebiten.Image) {
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	op.GeoM.Scale(scaleX, scaleY)
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
				g.drawPiece(xPos, yPos+5, screen, g.images[pc])
			}
			if sq == g.sqSelected {
				const scaleParam = 1.02
				g.drawPieceScale(float64(xPos)/scaleParam, float64(yPos)/scaleParam, scaleParam, scaleParam, screen, g.images[pc])
			}
			if sq == src(g.mvLast) || sq == dst(g.mvLast) {
				g.drawPiece(xPos, yPos, screen, g.images[ImgSelect])
			}
		}
	}
}

// clickSquare 点击格子后的处理
func (g *Game) clickSquare(xPos, yPos int) {
	// fmt.Println(xPos, yPos)
	// 棋子的值，比如说19
	sq := squareXY(xPos, yPos) // 0-255
	piece := 0
	if !g.bFilpped {
		piece = g.singlePosition.ucpcSquares[sq]
	} else {
		piece = g.singlePosition.ucpcSquares[squareFlip(sq)]
	}
	// 获得红黑标记
	if (piece & sideTag(g.singlePosition.sdPlayer)) != 0 {
		// 点击了自己方的棋子
		// 直接选中个这棋子
		g.sqSelected = sq
		g.playAudio(MusicSelect)
	} else if g.sqSelected != 0 {
		// 点击的不是自己方的棋子，那么直接走这个棋子
		mv := move(g.sqSelected, sq)
		if g.singlePosition.legalMove(mv) {
			if g.singlePosition.makeMove(mv) {

				// 保存上一步走法
				g.mvLast = mv
				// 把我们的选中的格子清0
				g.sqSelected = 0
				if piece == 0 {
					g.playAudio(MusicPut)
				} else {
					g.playAudio(MusicEat)
				}
			} else {
				g.playAudio(MusicJiang)
			}
		}
		// 如果不符合走法，不做处理
	}
}
