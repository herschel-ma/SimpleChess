package ChineseChess

const (
	ImgChessBoard = 1
	ImgSelect     = 2
	ImgRedShuai   = 8
	ImgRedShi     = 9
	ImgRedXiang   = 10
	ImgRedMa      = 11
	ImgRedJu      = 12
	ImgRedPao     = 13
	ImgRedBing    = 14

	ImgBlackJiang = 16
	ImgBlackShi   = 17
	ImgBlackXiang = 18
	ImgBlackMa    = 19
	ImgBlackJu    = 20
	ImgBlackPao   = 21
	ImgBlackBing  = 22
)

const (
	MusicSelect   = 100
	MusicPut      = 101
	MusicEat      = 102
	MusicJiang    = 103
	MusicGameWin  = 104
	MusicGameLose = 105
)

const sampleRate = 44100

// 窗口
const (
	SquareSize  = 56
	BoardEdge   = 8
	BoardWidth  = BoardEdge*2 + SquareSize*9
	BoardHeight = BoardEdge*2 + SquareSize*10
)

// 棋盘的范围
const (
	Top    = 3
	Bottom = 12
	Left   = 3
	Right  = 11
)

// 棋盘初始设置
var cucpcStartUp = [256]int{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 20, 19, 18, 17, 16, 17, 18, 19, 20, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 21, 0, 0, 0, 0, 0, 21, 0, 0, 0, 0, 0,
	0, 0, 0, 22, 0, 22, 0, 22, 0, 22, 0, 22, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 14, 0, 14, 0, 14, 0, 14, 0, 14, 0, 0, 0, 0,
	0, 0, 0, 0, 13, 0, 0, 0, 0, 0, 13, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 12, 11, 10, 9, 8, 9, 10, 11, 12, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

// ccInBoard 判断是否在棋盘中的数组
var ccInBoard = []int{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

// getY 获得格子的Y坐标
func getY(sq int) int {
	return sq >> 4
}

// getX 获得格子的X坐标
func getX(sq int) int {
	return sq & 15
}

// squareXY 根据x, y坐标获得的格子的值 [0 - 255]
func squareXY(x, y int) int {
	return x + (y << 4)
}

// squareFlip 翻转格子
func squareFlip(sq int) int {
	return 254 - sq
}

// xFlip X坐标水平镜像
func xFlip(x int) int {
	return 14 - x
}

// yFlip y坐标垂直镜像
func yFlip(y int) int {
	return 15 - y
}

// sideTag 获得红黑标记 (红方为8, 黑方为16)
func sideTag(sd int) int {
	return 8 + (sd << 3)
}

// oppSideTag 获得对方红黑标记
func oppSideTag(sd int) int {
	return 16 - (sd << 3)
}

// move 根据起点和终点的sq值来获得一个走法
func move(sqSrc, sqDst int) int {
	return sqSrc + (sqDst << 8)
}

// src 获得走法起点的值
func src(mv int) int {
	return mv & 255
}

// dst 获得走法终点的值
func dst(mv int) int {
	return mv >> 8
}
