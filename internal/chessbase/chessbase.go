package chessbase

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gavink97/pgn-tools/internal/global"
	"golang.org/x/exp/mmap"
)

func VerifyChessbaseInput(file string) bool {
	if file == "" {
		global.Logger.Error("Enter input filepath")
		return false
	}

	global.Logger.Debug(fmt.Sprintf("input file: %s", file))

	filetype := filepath.Ext(file)
	filename := strings.TrimSuffix(file, filetype)

	if strings.EqualFold(filetype, "cbh") || strings.EqualFold(filetype, "cbg") {
		global.Logger.Error(fmt.Sprintf("Invalid Filetype: %s", file))
		return false
	}

	_, err := os.Stat(file)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Invalid Filepath: %s", file))
		return false
	}

	for _, e := range []string{"cbh", "cbg", "cbp", "cbt", "cbe"} {
		alt := strings.Join([]string{filename, e}, ".")

		_, err := os.Stat(alt)
		if err != nil {
			global.Logger.Error(fmt.Sprintf("Invalid Filepath: %s", alt))
			return false
		}
	}

	return true
}

func ReadMMap(file string) (CBRead, error) {
	read, err := mmap.Open(file)
	if err != nil {
		global.Logger.Debug(err.Error())
		return CBRead{}, err
	}

	cb := make([]byte, read.Len())

	_, err = read.ReadAt(cb, 0)
	if err != nil {
		global.Logger.Debug(err.Error())
		return CBRead{}, err
	}

	cbRead := NewCBReader(CBReadParams{
		Read: *read,
		File: cb,
	})

	global.Logger.Debug(fmt.Sprintf("%s, %d", filepath.Base(file), read.Len()))

	return *cbRead, nil
}

func (cbi ChessBaseGameInfo) verifiedCBGame(cbhRecord []byte) bool {
	// ignore these games and split them off for debugging later
	if cbi.IsSpecialEncoded {
		return false
	}

	if isGame(cbhRecord) && !isMarkedDeleted(cbhRecord) && !cbi.IsEncoded && !cbi.Is960 {
		return true
	}
	return false
}

func IsWhiteTurnFEN(fen string) (bool, error) {
	splits := strings.Split(fen, " ")

	switch splits[1] {
	case "w":
		return true, nil
	case "b":
		return false, nil
	default:
		return false, fmt.Errorf("unexpected character in FEN: %s", splits[1])
	}
}

func GetMoveNoFEN(fen string) (int, error) {
	splits := strings.Split(fen, " ")

	i, err := strconv.Atoi(splits[len(splits)-1])
	if err != nil {
		return -1, err
	}

	return i, nil
}

//nolint:unused // Good for debugging
func printPosition(chessboard *Chessboard) {
	position := chessboard.Position

	for i := 7; i >= 0; i-- {
		square := ""
		for j := range 8 {
			t := position[j][i].PieceType
			pieceNo := position[j][i].PieceNo

			var c string
			if pieceNo == -1 {
				c = "_"
			} else {
				c = fmt.Sprintf("%d", pieceNo)
			}

			switch t {
			case W_QUEEN:
				square += fmt.Sprintf(" (q,%s)", c)
			case W_KING:
				square += fmt.Sprintf(" (k,%s)", c)
			case W_ROOK:
				square += fmt.Sprintf(" (r,%s)", c)
			case W_BISHOP:
				square += fmt.Sprintf(" (b,%s)", c)
			case W_KNIGHT:
				square += fmt.Sprintf(" (n,%s)", c)
			case W_PAWN:
				square += fmt.Sprintf(" (p,%s)", c)
			case B_QUEEN:
				square += fmt.Sprintf(" (Q,%s)", c)
			case B_KING:
				square += fmt.Sprintf(" (K,%s)", c)
			case B_ROOK:
				square += fmt.Sprintf(" (R,%s)", c)
			case B_BISHOP:
				square += fmt.Sprintf(" (B,%s)", c)
			case B_KNIGHT:
				square += fmt.Sprintf(" (N,%s)", c)
			case B_PAWN:
				square += fmt.Sprintf(" (P,%s)", c)
			default:
				square += " ....."
			}
		}
		fmt.Println(square)
	}
}
