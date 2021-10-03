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

// legalMove 判断走法是否合理
func (p *Position) legalMove(mv int) bool {
	// 判断开始的位置是否有自己的棋子
	sqSrc := src(mv)
	pcSrc := p.ucpcSquares[sqSrc]
	// 获得红黑标记
	pcSelfSide := sideTag(p.sdPlayer)
	if (pcSrc & pcSelfSide) == 0 {
		return false
	}

	// 判断终点的位置是否有自己的棋子
	sqDst := dst(mv)
	pcDst := p.ucpcSquares[sqDst]
	if (pcDst & pcSelfSide) != 0 {
		return false
	}

	// 根据棋子的类型判断是否符合象棋的规则
	// 拿到棋子的编号
	tmpPiece := pcSrc - pcSelfSide
	switch tmpPiece {
	case PieceJiang:
		return inFort(sqDst) && ccJiangSpan(sqSrc, sqDst)
	case PieceShi:
		return inFort(sqDst) && ccShiSpan(sqSrc, sqDst)
	case PieceXiang:
		xiangPinParam := xiangPin(sqSrc, sqDst)
		return sameRiver(sqSrc, sqDst) && ccXiangSpan(sqSrc, sqDst) && p.ucpcSquares[xiangPinParam] == 0
	case PieceMa:
		sqPin := maPin(sqSrc, sqDst)
		return sqSrc != sqPin && p.ucpcSquares[sqPin] == 0
	case PieceJu, PiecePao:
		// 某一个方向上的步长
		nDelta := 0
		if sameX(sqSrc, sqDst) {
			if sqDst < sqSrc {
				nDelta = -1
			} else {
				nDelta = 1
			}
		} else if sameY(sqSrc, sqDst) {
			if sqDst < sqSrc {
				nDelta = -16
			} else {
				nDelta = 16
			}
		} else {
			return false
		}
		sqTemp := sqSrc + nDelta
		for sqTemp != sqDst && p.ucpcSquares[sqTemp] == 0 {
			// 没有碰到子，就按nDelta往前走
			sqTemp += nDelta
		}
		if sqTemp == sqDst {
			// 如果终点没有子，不管是炮还是车，都是合法的。或者说只要是车，就合法
			return pcDst == 0 || tmpPiece == PieceJu
		} else if pcDst != 0 && tmpPiece == PiecePao {
			sqTemp += nDelta
			for sqTemp != sqDst && p.ucpcSquares[sqTemp] == 0 {
				sqTemp += nDelta
			}
			return sqTemp == sqDst
		} else {
			return false
		}
	case PieceBing:
		if hasRiver(sqDst, p.sdPlayer) && (sqDst == sqSrc-1 || sqDst == sqSrc+1) {
			return true
		}
		return sqDst == squareForward(sqSrc, p.sdPlayer)
	default:
	}
	return false
}
