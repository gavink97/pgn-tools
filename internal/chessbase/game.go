package chessbase

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gavink97/pgn-tools/internal/global"
	"github.com/gavink97/pgn-tools/internal/types"
)

func (cb *ChessBaseRecord) ExtractGame() (*types.Game, error) {
	if !isGame(cb.CBHRecord) || isMarkedDeleted(cb.CBHRecord) {
		return nil, errors.New("record is not a game or has been marked deleted")
	}

	whiteOffset, err := getWhiteOffset(cb.CBHRecord)
	if err != nil {
		return nil, err
	}

	white, err := getPlayer(cb.CBP, whiteOffset)
	if err != nil {
		return nil, err
	}

	blackOffset, err := getBlackOffset(cb.CBHRecord)
	if err != nil {
		return nil, err
	}

	black, err := getPlayer(cb.CBP, blackOffset)
	if err != nil {
		return nil, err
	}

	date, err := getDate(cb.CBHRecord)
	if err != nil {
		return nil, err
	}

	result, err := getResult(cb.CBHRecord)
	if err != nil {
		return nil, err
	}

	tournamentOffset, err := getTournamentOffset(cb.CBHRecord)
	if err != nil {
		return nil, err
	}

	tournamentData, err := getEventSiteRounds(cb.CBT, tournamentOffset)
	if err != nil {
		return nil, err
	}

	event := tournamentData[0]
	site := tournamentData[1]

	round, err := getRoundSubround(cb.CBHRecord)
	if err != nil {
		return nil, err
	}

	ratingData, err := getRatings(cb.CBHRecord)
	if err != nil {
		return nil, err
	}

	whiteElo := ratingData[0]
	blackElo := ratingData[1]

	gameOffset, err := getGameOffset(cb.CBHRecord)
	if err != nil {
		return nil, err
	}

	gameInfo, err := getGameInfo(cb.CBG, gameOffset)
	if err != nil {
		return nil, err
	}

	// print debug game information
	if !gameInfo.verifiedCBGame(cb.CBHRecord) {
		return nil, errors.New("invalid Chessbase Game")
	}

	var fen string
	var chessboard *Chessboard
	var decodeOffset int

	if gameInfo.ATypicalStart {
		fen, chessboard, err = decodeStartPosition(cb.CBG, gameOffset)
		if err != nil {
			global.Logger.Warn("unable to decode starting position")
			return nil, err
		}

		decodeOffset = 32
	} else {
		fen = ""
		chessboard = InitialChessboard()
		decodeOffset = 4
	}

	game, err := Decode(cb.CBG[gameOffset+decodeOffset:gameOffset+gameInfo.GameLength], chessboard, fen)
	if err != nil {
		global.Logger.Warn("unable to decode chess game due to error.")
		return nil, err
	}

	game += result

	return types.NewGame(types.GameParams{
		Event:    event,
		Site:     site,
		Date:     date,
		Round:    round,
		White:    white,
		Black:    black,
		Result:   result,
		BlackElo: blackElo,
		WhiteElo: whiteElo,
		FEN:      fen,
		Game:     game,
	}), nil
}

func getGameInfo(cbgFile []byte, gameNo int) (*ChessBaseGameInfo, error) {
	if len(cbgFile) < gameNo+4 {
		return nil, fmt.Errorf("cbg too short: expected atleast %d bytes, got %d", gameNo+4, len(cbgFile))
	}

	data := binary.BigEndian.Uint32(cbgFile[gameNo : gameNo+4])

	aTypicalStart := ((data & MASK_START_WITH_INITIAL) >> 30) != 0
	isEncoded := ((data & MASK_IS_ENCODED) >> 31) != 0
	isSpecialEncoded := ((data & MASK_SPECIAL_ENCODING) >> 26) != 0
	is960 := (data & MASK_IS_960) != 0
	gameLen := int(data & MASK_GAME_LEN)

	gameInfo := NewChessBaseGameInfo(ChessBaseGameInfoParams{
		GameLength:       gameLen,
		ATypicalStart:    aTypicalStart,
		IsEncoded:        isEncoded,
		IsSpecialEncoded: isSpecialEncoded,
		Is960:            is960,
	})

	return gameInfo, nil
}

func decodePieceLocation(stream string) (*Chessboard, error) {
	sIdx := 0
	bIdx := 0

	board := EmptyChessboard()
	pieceCount := [13]int{}

	for sIdx < len(stream) && bIdx < 64 {
		if string(stream[sIdx]) == "0" {
			sIdx++
			bIdx++
			continue
		}

		if len(stream)-sIdx < 5 {
			return nil, fmt.Errorf("error decoding position: %s", stream)
		}

		piece := stream[sIdx : sIdx+5]
		abs := ABS_TO_XY[bIdx]
		i, j := abs.X, abs.Y

		switch piece {
		case "10001":
			board.Position[i][j] = PieceInfo{PieceType: W_KING, PieceNo: 0} // can only have 1
			board.PieceList[W_KING][0] = Coord{X: i, Y: j}
		case "10010":
			idx := pieceCount[W_QUEEN]
			board.Position[i][j] = PieceInfo{PieceType: W_QUEEN, PieceNo: idx}
			board.PieceList[W_QUEEN][idx] = Coord{X: i, Y: j}
			pieceCount[W_QUEEN]++
		case "10011":
			idx := pieceCount[W_KNIGHT]
			board.Position[i][j] = PieceInfo{PieceType: W_KNIGHT, PieceNo: idx}
			board.PieceList[W_KNIGHT][idx] = Coord{X: i, Y: j}
			pieceCount[W_KNIGHT]++
		case "10100":
			idx := pieceCount[W_BISHOP]
			board.Position[i][j] = PieceInfo{PieceType: W_BISHOP, PieceNo: idx}
			board.PieceList[W_BISHOP][idx] = Coord{X: i, Y: j}
			pieceCount[W_BISHOP]++
		case "10101":
			idx := pieceCount[W_ROOK]
			board.Position[i][j] = PieceInfo{PieceType: W_ROOK, PieceNo: idx}
			board.PieceList[W_ROOK][idx] = Coord{X: i, Y: j}
			pieceCount[W_ROOK]++
		case "10110":
			idx := pieceCount[W_PAWN]
			board.Position[i][j] = PieceInfo{PieceType: W_PAWN, PieceNo: idx}
			board.PieceList[W_PAWN][idx] = Coord{X: i, Y: j}
			pieceCount[W_PAWN]++
		case "11001":
			board.Position[i][j] = PieceInfo{PieceType: B_KING, PieceNo: 0} // can only have 1
			board.PieceList[B_KING][0] = Coord{X: i, Y: j}
		case "11010":
			idx := pieceCount[B_QUEEN]
			board.Position[i][j] = PieceInfo{PieceType: B_QUEEN, PieceNo: idx}
			board.PieceList[B_QUEEN][idx] = Coord{X: i, Y: j}
			pieceCount[B_QUEEN]++
		case "11011":
			idx := pieceCount[B_KNIGHT]
			board.Position[i][j] = PieceInfo{PieceType: B_KNIGHT, PieceNo: idx}
			board.PieceList[B_KNIGHT][idx] = Coord{X: i, Y: j}
			pieceCount[B_KNIGHT]++
		case "11100":
			idx := pieceCount[B_BISHOP]
			board.Position[i][j] = PieceInfo{PieceType: B_BISHOP, PieceNo: idx}
			board.PieceList[B_BISHOP][idx] = Coord{X: i, Y: j}
			pieceCount[B_BISHOP]++
		case "11101":
			idx := pieceCount[B_ROOK]
			board.Position[i][j] = PieceInfo{PieceType: B_ROOK, PieceNo: idx}
			board.PieceList[B_ROOK][idx] = Coord{X: i, Y: j}
			pieceCount[B_ROOK]++
		case "11110":
			idx := pieceCount[B_PAWN]
			board.Position[i][j] = PieceInfo{PieceType: B_PAWN, PieceNo: idx}
			board.PieceList[B_PAWN][idx] = Coord{X: i, Y: j}
			pieceCount[B_PAWN]++
		default:
			return nil, fmt.Errorf("invalid piece : %s pos: %v from %s", piece, sIdx, stream)
		}

		sIdx += 5
		bIdx++
	}

	return board, nil

}

// only supports standard chess no 960 X-FEN
func posToFEN(position [8][8]PieceInfo, ep_file int, isBlackTurn bool, whiteLong bool, whiteShort bool, blackLong bool, blackShort bool, nextMoveNo int) (string, error) {
	fen := ""

	for i := 7; i >= 0; i-- {
		square := 0
		for j := range 8 {
			pieceType := position[j][i].PieceType
			if pieceType == 0 {
				square++
			} else {
				if square > 0 {
					fen += strconv.Itoa(square)
					square = 0
				}

				switch pieceType {
				case W_KING:
					fen += "K"
				case W_QUEEN:
					fen += "Q"
				case W_ROOK:
					fen += "R"
				case W_BISHOP:
					fen += "B"
				case W_KNIGHT:
					fen += "N"
				case W_PAWN:
					fen += "P"
				case B_KING:
					fen += "k"
				case B_QUEEN:
					fen += "q"
				case B_ROOK:
					fen += "r"
				case B_BISHOP:
					fen += "b"
				case B_KNIGHT:
					fen += "n"
				case B_PAWN:
					fen += "p"
				default:
					return "", fmt.Errorf("unknown piece: %v", pieceType)
				}
			}
		}
		if square > 0 {
			fen += strconv.Itoa(square)
		}
		fen += "/"
	}

	fen = fen[:len(fen)-1]

	if isBlackTurn {
		fen += " b"
	} else {
		fen += " w"
	}

	if whiteShort || whiteLong || blackShort || blackLong {
		fen += " "
		if whiteShort {
			fen += "K"
		}

		if whiteLong {
			fen += "Q"
		}

		if blackShort {
			fen += "k"
		}

		if blackLong {
			fen += "q"
		}
	} else {
		fen += " -"
	}

	if ep_file > 0 {
		fen += " "
		switch ep_file {
		case 1:
			fen += "a"
		case 2:
			fen += "b"
		case 3:
			fen += "c"
		case 4:
			fen += "d"
		case 5:
			fen += "e"
		case 6:
			fen += "f"
		case 7:
			fen += "g"
		case 8:
			fen += "h"
		default:
			return "", fmt.Errorf("unknown ep file encoding: %v", ep_file)
		}

		if isBlackTurn {
			fen += "3 "
		} else {
			fen += "6 "
		}
	} else {
		fen += " - "
	}

	fen += "0 "
	fen += fmt.Sprintf("%d", nextMoveNo)
	return fen, nil
}

func decodeStartPosition(cbgFile []byte, gameNo int) (string, *Chessboard, error) {
	ep_file := int(cbgFile[gameNo+4+1] & byte(MASK_EP_FILE))
	isBlackTurn := int((cbgFile[gameNo+4+1]&byte(MASK_TURN))>>4) == 1
	whiteCastleLong := int(cbgFile[gameNo+4+2]&byte(MASK_WHITE_CASTLE_LONG)) == 1
	whiteCastleShort := int((cbgFile[gameNo+4+2]&byte(MASK_WHITE_CASTLE_SHORT))>>1) == 1
	blackCastleLong := int((cbgFile[gameNo+4+2]&byte(MASK_BLACK_CASTLE_LONG))>>2) == 1
	blackCastleShort := int((cbgFile[gameNo+4+2]&byte(MASK_BLACK_CASTLE_SHORT))>>3) == 1
	nextMoveNo := int(cbgFile[gameNo+4+3])
	bitstream := cbgFile[gameNo+8 : gameNo+8+24]

	setupBits := ""

	for i := range len(bitstream) {
		setupBits += fmt.Sprintf("%08b", bitstream[i])
	}

	chessboard, err := decodePieceLocation(setupBits)
	if err != nil {
		return "", nil, err
	}

	fen, err := posToFEN(chessboard.Position, ep_file, isBlackTurn, whiteCastleLong, whiteCastleShort, blackCastleLong, blackCastleShort, nextMoveNo)
	if err != nil {
		return "", nil, err
	}

	return fen, chessboard, nil
}

func decreasePieceNR(cb *Chessboard, targetPieceInfo PieceInfo) {
	pieceType := targetPieceInfo.PieceType
	pieceNo := targetPieceInfo.PieceNo

	for nr := pieceNo; nr < 7; nr++ {
		cb.PieceList[pieceType][nr] = cb.PieceList[pieceType][nr+1]
	}

	cb.PieceList[pieceType][7] = Coord{-1, -1}

	for x := range 8 {
		for y := range 8 {
			p, t := cb.Position[x][y].PieceType, cb.Position[x][y].PieceNo
			if p == pieceType && t > pieceNo {
				cb.Position[x][y] = PieceInfo{PieceType: p, PieceNo: t - 1}
			}
		}
	}
}

func (cb *Chessboard) doMove(pieceInfo PieceInfo, cbEnc map[byte]Coord, token byte, pawnFlip bool) (string, error) {
	pieceType := pieceInfo.PieceType
	pieceNo := pieceInfo.PieceNo

	i := cb.PieceList[pieceType][pieceNo].X
	j := cb.PieceList[pieceType][pieceNo].Y

	if i < 0 || i > 7 || j < 0 || j > 7 {
		return "", fmt.Errorf("invalid piece coordinates: (%d,%d) for piece type %d, number %d", i, j, pieceType, pieceNo)
	}

	cb.Position[i][j] = PieceInfo{PieceType: EMPTY, PieceNo: -1}

	addX := cbEnc[token].X
	addY := cbEnc[token].Y

	if pawnFlip {
		addX = -addX
		addY = -addY
	}

	i1 := (i + addX) % 8
	j1 := (j + addY) % 8

	tPieceType := cb.Position[i1][j1].PieceType

	var move string

	if tPieceType != EMPTY {
		move = strings.Join([]string{PIECE_MAP[pieceType], "x", SQN[i1][j1]}, "")
		if pieceType == W_PAWN || pieceType == B_PAWN {
			move = strings.Join([]string{SQN[i][j][:1], "x", SQN[i1][j1]}, "")
		}
	} else {
		move = strings.Join([]string{PIECE_MAP[pieceType], SQN[i1][j1]}, "")
	}

	if tPieceType != EMPTY && tPieceType != W_KING && tPieceType != B_KING && tPieceType != W_PAWN && tPieceType != B_PAWN {
		decreasePieceNR(cb, cb.Position[i1][j1])
	}

	cb.Position[i1][j1] = pieceInfo
	cb.PieceList[pieceType][pieceNo] = Coord{X: i1, Y: j1}

	if pieceType == W_KING && token == 0x76 { // castle short white
		move = "O-O"
		cb.Position[7][0] = PieceInfo{PieceType: EMPTY, PieceNo: -1}
		for idx := range len(cb.PieceList[W_ROOK]) {
			if cb.PieceList[W_ROOK][idx].X == 7 && cb.PieceList[W_ROOK][idx].Y == 0 {
				cb.PieceList[W_ROOK][idx] = Coord{5, 0}
				cb.Position[5][0] = PieceInfo{PieceType: W_ROOK, PieceNo: idx}
				break
			}
		}
	}

	if pieceType == B_KING && token == 0x76 { // castle short black
		move = "O-O"
		cb.Position[7][7] = PieceInfo{PieceType: EMPTY, PieceNo: -1}
		for idx := range len(cb.PieceList[B_ROOK]) {
			if cb.PieceList[B_ROOK][idx].X == 7 && cb.PieceList[B_ROOK][idx].Y == 7 {
				cb.PieceList[B_ROOK][idx] = Coord{5, 7}
				cb.Position[5][7] = PieceInfo{PieceType: B_ROOK, PieceNo: idx}
				break
			}
		}
	}

	if pieceType == W_KING && token == 0xB5 { // castle long white
		move = "O-O-O"
		cb.Position[0][0] = PieceInfo{PieceType: EMPTY, PieceNo: -1}
		for idx := range len(cb.PieceList[W_ROOK]) {
			if cb.PieceList[W_ROOK][idx].X == 0 && cb.PieceList[W_ROOK][idx].Y == 0 {
				cb.PieceList[W_ROOK][idx] = Coord{3, 0}
				cb.Position[3][0] = PieceInfo{PieceType: W_ROOK, PieceNo: idx}
				break
			}
		}
	}

	if pieceType == B_KING && token == 0xB5 { // castle long black
		move = "O-O-O"
		cb.Position[0][7] = PieceInfo{PieceType: EMPTY, PieceNo: -1}
		for idx := range len(cb.PieceList[B_ROOK]) {
			if cb.PieceList[B_ROOK][idx].X == 0 && cb.PieceList[B_ROOK][idx].Y == 7 {
				cb.PieceList[B_ROOK][idx] = Coord{3, 7}
				cb.Position[3][7] = PieceInfo{PieceType: B_ROOK, PieceNo: idx}
				break
			}
		}
	}

	// checks, checkmate
	// check for disambiguations in the position (rooks, knights)
	// disambiguations for multiple queens, bishops, rooks, knights

	return move, nil
}

// promotion moves or when the fourth+ kind of one piece is moved (eg fourth white queen)
func (cb *Chessboard) do2bMove(coords []Coord, promotionPiece uint16) (string, error) {
	src := coords[0]
	dst := coords[1]

	pieceInfo := cb.Position[src.X][src.Y]
	pieceType := pieceInfo.PieceType
	pieceNo := pieceInfo.PieceNo

	cb.Position[src.X][src.Y] = PieceInfo{PieceType: EMPTY, PieceNo: -1}

	targetPieceInfo := cb.Position[dst.X][dst.Y]
	tPieceType := targetPieceInfo.PieceType

	if tPieceType != EMPTY && tPieceType != W_KING && tPieceType != B_KING && tPieceType != W_PAWN && tPieceType != B_PAWN {
		decreasePieceNR(cb, targetPieceInfo)
	}

	promotionStr := ""
	promotedPieceType := 0

	var move string

	if pieceType != W_PAWN && pieceType != B_PAWN {
		// assuming two byte encodings never occur for pawn moves unless a promotion
		cb.Position[dst.X][dst.Y] = pieceInfo
		cb.PieceList[pieceType][pieceNo] = Coord{dst.X, dst.Y}

		if tPieceType != EMPTY {
			move = fmt.Sprintf("%sx%s", PIECE_MAP[pieceType], SQN[dst.X][dst.Y])
		} else {
			move = fmt.Sprintf("%s%s", PIECE_MAP[pieceType], SQN[dst.X][dst.Y])
		}

		return move, nil
	} else {
		// two byte moves are used for promotions
		// two byte moves should usually never be castles

		if pieceType == W_PAWN && dst.Y == 7 {
			switch promotionPiece {
			case 0:
				promotedPieceType = W_QUEEN
				promotionStr += "Q"
			case 1:
				promotedPieceType = W_ROOK
				promotionStr += "R"
			case 2:
				promotedPieceType = W_BISHOP
				promotionStr += "B"
			case 3:
				promotedPieceType = W_KNIGHT
				promotionStr += "N"
			default:
				return "", fmt.Errorf("unknown promotion piece type: %d", promotionPiece)
			}
		}

		if pieceType == B_PAWN && dst.Y == 0 {
			switch promotionPiece {
			case 0:
				promotedPieceType = B_QUEEN
				promotionStr += "Q"
			case 1:
				promotedPieceType = B_ROOK
				promotionStr += "R"
			case 2:
				promotedPieceType = B_BISHOP
				promotionStr += "B"
			case 3:
				promotedPieceType = B_KNIGHT
				promotionStr += "N"
			default:
				return "", fmt.Errorf("unknown promotion piece type: %d", promotionPiece)
			}
		}
	}

	if promotedPieceType != EMPTY {
		free_idx := -1
		for idx := range 8 {
			if cb.PieceList[promotedPieceType][idx].X == -1 && cb.PieceList[promotedPieceType][idx].Y == -1 {
				free_idx = idx
				break
			}
		}

		if free_idx == -1 {
			return "", fmt.Errorf("no free index for promotion piece: %d", promotedPieceType)
		}

		cb.PieceList[promotedPieceType][free_idx] = Coord{X: dst.X, Y: dst.Y}
		cb.Position[dst.X][dst.Y] = PieceInfo{PieceType: promotedPieceType, PieceNo: free_idx}
	}

	if tPieceType != EMPTY {
		move = fmt.Sprintf("%sx%s=%s", SQN[src.X][src.Y][:1], SQN[dst.X][dst.Y], promotionStr)
	} else {
		move = fmt.Sprintf("%s=%s", SQN[dst.X][dst.Y], promotionStr)
	}

	return move, nil
}

func Decode(gameBytes []byte, chessboard *Chessboard, fen string) (string, error) {
	processedMoves := 0
	idx := 0
	game := ""

	variations := []*State{}

	var isWhiteToMove bool
	var err error
	var move string
	var moveNo int
	var moveFound bool

	if fen != "" {
		isWhiteToMove, err = IsWhiteTurnFEN(fen)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("unable to parse turn from FEN: %s", fen))
			return "", err
		}

		moveNo, err = GetMoveNoFEN(fen)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("unable to move number from FEN: %s", fen))
			return "", err
		}
	} else {
		isWhiteToMove = true
		moveNo = 1
	}

	for idx < len(gameBytes) {
		token := byte((int(gameBytes[idx]) - processedMoves) % 256)
		moveFound = false

		/*
			if int(token) == 12 {
				// last byte in gameBytes
				return game, nil
			}
		*/

		if !bytes.Contains(SPECIAL_CODES, []byte{token}) {
			processedMoves += 1
			processedMoves %= 256
		}

		if token == 0x9F {
			// byte skip
			idx += 1
			continue
		}

		if token == 0xAA {
			// null move
			if isWhiteToMove {
				game += fmt.Sprintf("%d. ", moveNo)
			} else {
				moveNo++
			}

			game += "-- "
			isWhiteToMove = !isWhiteToMove
			idx += 1
			continue
		}

		if token == 0x29 {
			// two byte move
			tmp := make([]byte, 2)
			tmp[0] = DEOBFUSCATE_2B[gameBytes[idx+1]-byte(processedMoves)]
			tmp[1] = DEOBFUSCATE_2B[gameBytes[idx+2]-byte(processedMoves)]

			twoByteMove := binary.BigEndian.Uint16(tmp)
			if twoByteMove == 0 {
				global.Logger.Warn(fmt.Sprintf("error decoding 2b move: %d", twoByteMove))
			}

			src := twoByteMove & 0x3F
			dst := (twoByteMove >> 6) & 0x3F
			promotePiece := (twoByteMove >> 12) & 0x3F

			coords := []Coord{
				ABS_TO_XY[src],
				ABS_TO_XY[dst],
			}

			move, err = chessboard.do2bMove(coords, promotePiece)
			if err != nil {
				global.Logger.Warn("an error occured performing a 2 byte move")
				return "", err
			}

			if isWhiteToMove {
				game += fmt.Sprintf("%d. ", moveNo)
			} else {
				moveNo++
			}

			game += fmt.Sprintf("%s ", move)
			isWhiteToMove = !isWhiteToMove

			processedMoves += 1
			processedMoves %= 256
			idx += 3
			continue
		}

		// create a new branch, all moves will go to this new branch and we will keep doing this for new branches
		// then we will in order append those branches back according to state with formatting
		if token == 0xDC {
			return "", fmt.Errorf("not handling variations atm")
			// start of variation

			/*
				state := NewState(StateParams{
					Chessboard:  chessboard.Clone(),
					MoveNo:      moveNo,
					IsWhiteTurn: isWhiteToMove,
				})

				variations = append(variations, state)
			*/
		}

		if token == 0x0C {
			// every game is terminated with 0x0C, ignore last otherwise pop from stack
			// end of variation

			if idx < len(gameBytes)-1 && len(variations) > 0 {
				state := variations[len(variations)-1]
				variations = variations[:len(variations)-1]

				chessboard = state.Chessboard
				moveNo = state.MoveNo
				isWhiteToMove = state.IsWhiteTurn
			}
		}

		if isWhiteToMove {
			if _, exists := CB_KING_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_KING,
					PieceNo:   0,
				}, CB_KING_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_QUEEN_1_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_QUEEN,
					PieceNo:   0,
				}, CB_QUEEN_1_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_QUEEN_2_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_QUEEN,
					PieceNo:   1,
				}, CB_QUEEN_2_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_QUEEN_3_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_QUEEN,
					PieceNo:   2,
				}, CB_QUEEN_3_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_ROOK_1_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_ROOK,
					PieceNo:   0,
				}, CB_ROOK_1_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_ROOK_2_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_ROOK,
					PieceNo:   1,
				}, CB_ROOK_2_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_ROOK_3_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_ROOK,
					PieceNo:   2,
				}, CB_ROOK_3_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_BISHOP_1_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_BISHOP,
					PieceNo:   0,
				}, CB_BISHOP_1_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_BISHOP_2_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_BISHOP,
					PieceNo:   1,
				}, CB_BISHOP_2_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_BISHOP_3_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_BISHOP,
					PieceNo:   2,
				}, CB_BISHOP_3_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_KNIGHT_1_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_KNIGHT,
					PieceNo:   0,
				}, CB_KNIGHT_1_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_KNIGHT_2_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_KNIGHT,
					PieceNo:   1,
				}, CB_KNIGHT_2_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_KNIGHT_3_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_KNIGHT,
					PieceNo:   2,
				}, CB_KNIGHT_3_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_A_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_PAWN,
					PieceNo:   0,
				}, CB_PAWN_A_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_B_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_PAWN,
					PieceNo:   1,
				}, CB_PAWN_B_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_C_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_PAWN,
					PieceNo:   2,
				}, CB_PAWN_C_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_D_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_PAWN,
					PieceNo:   3,
				}, CB_PAWN_D_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_E_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_PAWN,
					PieceNo:   4,
				}, CB_PAWN_E_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_F_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_PAWN,
					PieceNo:   5,
				}, CB_PAWN_F_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_G_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_PAWN,
					PieceNo:   6,
				}, CB_PAWN_G_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_H_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: W_PAWN,
					PieceNo:   7,
				}, CB_PAWN_H_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			}
		} else {
			if _, exists := CB_KING_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_KING,
					PieceNo:   0,
				}, CB_KING_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_QUEEN_1_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_QUEEN,
					PieceNo:   0,
				}, CB_QUEEN_1_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_QUEEN_2_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_QUEEN,
					PieceNo:   1,
				}, CB_QUEEN_2_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_QUEEN_3_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_QUEEN,
					PieceNo:   2,
				}, CB_QUEEN_3_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_ROOK_1_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_ROOK,
					PieceNo:   0,
				}, CB_ROOK_1_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_ROOK_2_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_ROOK,
					PieceNo:   1,
				}, CB_ROOK_2_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_ROOK_3_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_ROOK,
					PieceNo:   2,
				}, CB_ROOK_3_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_BISHOP_1_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_BISHOP,
					PieceNo:   0,
				}, CB_BISHOP_1_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_BISHOP_2_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_BISHOP,
					PieceNo:   1,
				}, CB_BISHOP_2_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_BISHOP_3_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_BISHOP,
					PieceNo:   2,
				}, CB_BISHOP_3_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_KNIGHT_1_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_KNIGHT,
					PieceNo:   0,
				}, CB_KNIGHT_1_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_KNIGHT_2_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_KNIGHT,
					PieceNo:   1,
				}, CB_KNIGHT_2_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_KNIGHT_3_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_KNIGHT,
					PieceNo:   2,
				}, CB_KNIGHT_3_ENC, token, false)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_A_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_PAWN,
					PieceNo:   0,
				}, CB_PAWN_A_ENC, token, true)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_B_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_PAWN,
					PieceNo:   1,
				}, CB_PAWN_B_ENC, token, true)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_C_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_PAWN,
					PieceNo:   2,
				}, CB_PAWN_C_ENC, token, true)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_D_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_PAWN,
					PieceNo:   3,
				}, CB_PAWN_D_ENC, token, true)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_E_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_PAWN,
					PieceNo:   4,
				}, CB_PAWN_E_ENC, token, true)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_F_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_PAWN,
					PieceNo:   5,
				}, CB_PAWN_F_ENC, token, true)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_G_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_PAWN,
					PieceNo:   6,
				}, CB_PAWN_G_ENC, token, true)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			} else if _, exists := CB_PAWN_H_ENC[token]; exists {
				move, err = chessboard.doMove(PieceInfo{
					PieceType: B_PAWN,
					PieceNo:   7,
				}, CB_PAWN_H_ENC, token, true)
				if err != nil {
					global.Logger.Warn("an error occured performing a move")
					return "", err
				}
				moveFound = true
			}
		}

		if moveFound {
			if isWhiteToMove {
				game += fmt.Sprintf("%d. ", moveNo)
			} else {
				moveNo++
			}

			game += fmt.Sprintf("%s ", move)
			isWhiteToMove = !isWhiteToMove
		}

		idx += 1
	}

	return game, nil
}
