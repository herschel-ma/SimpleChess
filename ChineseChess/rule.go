package ChineseChess

// Position 局面结构
type Position struct {
	sdPlayer    int      // 轮到谁走, 红方为0, 黑方为1
	ucpcSquares [256]int // 棋盘上的棋子
}

// NewPosition 初始化局面结构体
func NewPosition() *Position {
	p := &Position{}
	if p == nil {
		return nil
	}
	return p
}

// startup 初始化棋盘
func (p *Position) startup() {
	p.sdPlayer = 0
	for sq := 0; sq < 256; sq++ {
		p.ucpcSquares[sq] = cucpcStartUp[sq]
	}
}

// changeSide 交换走子方
func (p *Position) changeSide() {
	p.sdPlayer = 1 - p.sdPlayer
}

// addPiece 往棋盘上放置一枚棋子
func (p *Position) addPiece(sq, pc int) {
	p.ucpcSquares[sq] = pc
}

// delPiece 从棋盘上拿走一枚棋子
func (p *Position) delPiece(sq int) {
	p.ucpcSquares[sq] = 0
}

// movePiece 搬一步棋
func (p *Position) movePiece(mv int) {
	sqSrc := src(mv)
	sqDst := dst(mv)
	p.delPiece(sqDst)
	pc := p.ucpcSquares[sqSrc]
	p.delPiece(sqSrc)
	p.addPiece(sqDst, pc)
}

// makeMove 走一步棋
func (p *Position) makeMove(mv int) {
	p.movePiece(mv)
	p.changeSide()
}
