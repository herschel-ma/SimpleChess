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
// 返回被吃的子
func (p *Position) movePiece(mv int) int {
	sqSrc := src(mv)
	sqDst := dst(mv)
	pcCaptured := p.ucpcSquares[sqDst]
	p.delPiece(sqDst)
	pc := p.ucpcSquares[sqSrc]
	p.delPiece(sqSrc)
	p.addPiece(sqDst, pc)
	return pcCaptured
}

// makeMove 走一步棋
// 被将军之后不让走这步棋去送死
func (p *Position) makeMove(mv int) bool {
	pcCaptured := p.movePiece(mv)
	if p.checked() {
		p.undoMovePiece(mv, pcCaptured)
		return false
	}
	p.changeSide()
	return true
}

// undoMovePiece movePiece的逆向操作
func (p *Position) undoMovePiece(mv, pcCaptured int) {
	sqSrc := src(mv)
	sqDst := dst(mv)
	pc := p.ucpcSquares[sqDst]
	p.delPiece(sqDst)
	p.addPiece(sqSrc, pc)
	if pcCaptured != 0 {
		p.addPiece(sqDst, pcCaptured)
	}
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

// checked 判断是否被将军
func (p *Position) checked() bool {
	// 获得红黑标记
	pcSelfSide := sideTag(p.sdPlayer)
	// 获得对方的红黑标记
	pcOppSide := oppSideTag(p.sdPlayer)
	// 找到本方棋盘上的将
	for sqSrc := 0; sqSrc < 256; sqSrc++ {
		if !inBoard(sqSrc) || p.ucpcSquares[sqSrc] != pcSelfSide+PieceJiang {
			continue
		}
		// 是否被对方的兵将军
		if p.ucpcSquares[squareForward(sqSrc, p.sdPlayer)] == pcOppSide+PieceBing {
			return true
		}
		// 左右移动之后是否能碰到对方的兵
		for nDelta := -1; nDelta <= 1; nDelta += 2 {
			if p.ucpcSquares[sqSrc+nDelta] == pcOppSide+PieceBing {
				return true
			}
		}

		// 是否被对方的马将军
		for i := 0; i < 4; i++ {
			// 0  2  0  0
			// 2  -1  0  0
			// 0  0  1  0
			if p.ucpcSquares[sqSrc+ccShiDelta[i]] != 0 {
				continue
			}
			for j := 0; j < 2; j++ {
				if p.ucpcSquares[sqSrc+ccMaCheckDelta[i][j]] == pcOppSide+PieceMa {
					return true
				}
			}
		}

		// 判断是否被对方的车、炮将军(包括将对脸)
		for i := 0; i < 4; i++ {
			// 确定一个方向
			nDelta := ccJiangDelta[i]
			sqDst := sqSrc + nDelta
			for inBoard(sqDst) {
				pcDst := p.ucpcSquares[sqDst]
				if pcDst != 0 {
					if pcDst == pcOppSide+PieceJu || pcDst == pcOppSide+PieceJiang {
						return true
					}
					break
				}
				sqDst += nDelta
			}
			sqDst += nDelta
			// 是否被炮将军
			for inBoard(sqDst) {
				pcDst := p.ucpcSquares[sqDst]
				if pcDst != 0 {
					if pcDst == pcOppSide+PiecePao {
						return true
					}
					break
				}
				sqDst += nDelta
			}
		}
		return false
	}
	return false
}

//  generateMoves 生成所有的走法, Return 有多少走法
func (p *Position) generateMoves(mvs []int) int {
	numGenerateMVS := 0
	// 获得红黑标记
	pcSelfSide := sideTag(p.sdPlayer)
	// 获得对方红黑标记
	pcOppSide := oppSideTag(p.sdPlayer)
	for sqSrc := 0; sqSrc < 256; sqSrc++ {
		if !inBoard(sqSrc) {
			continue
		}
		// 直到找到一个本方棋子
		if p.ucpcSquares[sqSrc]&pcSelfSide == 0 {
			continue
		}

		// 根据棋子确定走法
		switch p.ucpcSquares[sqSrc] - pcSelfSide {
		case PieceJiang:
			for i := 0; i < 4; i++ {
				sqDst := sqSrc + ccJiangDelta[i]
				if !inFort(sqDst) {
					continue
				}
				pcDst := p.ucpcSquares[sqDst]
				// 包含了pcDst 是0的这种情况
				if pcDst&pcSelfSide == 0 {
					// 走法合理，保存到mvs
					mv := move(sqSrc, sqDst)
					mvs[numGenerateMVS] = mv
					// 将走法数加1
					numGenerateMVS++
				}
			}
		case PieceShi:
			for i := 0; i < 4; i++ {
				sqDst := sqSrc + ccShiDelta[i]
				if !inFort(sqDst) {
					continue
				}
				pcDst := p.ucpcSquares[sqDst]
				if pcDst&pcSelfSide == 0 {
					mvs[numGenerateMVS] = move(sqSrc, sqDst)
					numGenerateMVS++
				}
			}
		case PieceXiang:
			for i := 0; i < 4; i++ {
				sqDst := sqSrc + ccShiDelta[i]
				if !(inBoard(sqDst) && noRiver(sqDst, p.sdPlayer) && p.ucpcSquares[sqDst] == 0) {
					continue
				}
				sqDst += ccShiDelta[i]
				pcDst := p.ucpcSquares[sqDst]
				if pcDst&pcSelfSide == 0 {
					mvs[numGenerateMVS] = move(sqSrc, sqDst)
					numGenerateMVS++
				}
			}
		case PieceMa:
			for i := 0; i < 4; i++ {
				// 0  1  0
				// 1  0  0
				// 0  0  2
				sqDst := sqSrc + ccJiangDelta[i]
				if p.ucpcSquares[sqDst] != 0 {
					continue
				}
				for j := 0; j < 2; j++ {
					sqDst := sqSrc + ccMaDelta[i][j]
					if !inBoard(sqDst) {
						continue
					}
					pcDst := p.ucpcSquares[sqDst]
					if pcDst&pcSelfSide == 0 {
						mvs[numGenerateMVS] = move(sqSrc, sqDst)
						numGenerateMVS++
					}
				}
			}
		case PieceJu:
			for i := 0; i < 4; i++ {
				nDelta := ccJiangDelta[i]
				sqDst := sqSrc + nDelta
				for inBoard(sqDst) {
					pcDst := p.ucpcSquares[sqDst]
					if pcDst == 0 {
						mvs[numGenerateMVS] = move(sqSrc, sqDst)
						numGenerateMVS++
					} else {
						if pcDst&pcOppSide != 0 {
							mvs[numGenerateMVS] = move(sqSrc, sqDst)
							numGenerateMVS++
						}
						break
					}
					sqDst += nDelta
				}
			}
		case PiecePao:
			for i := 0; i < 4; i++ {
				nDelta := ccJiangDelta[i]
				sqDst := sqSrc + nDelta
				for inBoard(sqDst) {
					pcDst := p.ucpcSquares[sqDst]
					if pcDst == 0 {
						mvs[numGenerateMVS] = move(sqSrc, sqDst)
						numGenerateMVS++
					} else {
						break
					}
					sqDst += nDelta
				}
				sqDst += nDelta
				for inBoard(sqDst) {
					pcDst := p.ucpcSquares[sqDst]
					if pcDst != 0 {
						if pcDst&pcOppSide != 0 {
							// 翻过山之后必须是对方的棋子
							mvs[numGenerateMVS] = move(sqSrc, sqDst)
							numGenerateMVS++
						}
						break
					}
					sqDst += nDelta
				}
			}
		case PieceBing:
			sqDst := squareForward(sqSrc, p.sdPlayer)
			if inBoard(sqDst) {
				pcDst := p.ucpcSquares[sqDst]
				if pcDst&pcSelfSide == 0 {
					// 空或者不是自己方的棋子
					mvs[numGenerateMVS] = move(sqSrc, sqDst)
					numGenerateMVS++
				}
				if hasRiver(sqSrc, p.sdPlayer) {
					for nDelta := -1; nDelta < 2; nDelta += 2 {
						sqDst := sqSrc + nDelta
						if inBoard(sqDst) {
							pcDst := p.ucpcSquares[sqDst]
							if pcDst&pcSelfSide == 0 {
								mvs[numGenerateMVS] = move(sqSrc, sqDst)
								numGenerateMVS++
							}
						}
					}
				}
			}
		}
	}
	return numGenerateMVS
}

// isMate 判断是否被将死
func (p *Position) isMate() bool {
	mvs := make([]int, MaxGenMoves) // 初始化一个保存走法的数组
	nGeneMoveNum := p.generateMoves(mvs)
	// fmt.Println(nGeneMoveNum)
	// fmt.Println(mvs)
	for i := 0; i < nGeneMoveNum; i++ {
		pcCaptured := p.movePiece(mvs[i])
		if !p.checked() {
			p.undoMovePiece(mvs[i], pcCaptured)
			return false
		}
		p.undoMovePiece(mvs[i], pcCaptured)
	}
	return true
}
