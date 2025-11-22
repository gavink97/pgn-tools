package chessbase

import "golang.org/x/exp/mmap"

type CBRead struct {
	Read mmap.ReaderAt
	File []byte
}

type CBReadParams struct {
	Read mmap.ReaderAt
	File []byte
}

func NewCBReader(params CBReadParams) *CBRead {
	return &CBRead{
		Read: params.Read,
		File: params.File,
	}
}

type ChessBaseGameInfo struct {
	GameLength       int
	ATypicalStart    bool
	IsEncoded        bool
	IsSpecialEncoded bool
	Is960            bool
}

type ChessBaseGameInfoParams struct {
	GameLength       int
	ATypicalStart    bool
	IsEncoded        bool
	IsSpecialEncoded bool
	Is960            bool
}

func NewChessBaseGameInfo(params ChessBaseGameInfoParams) *ChessBaseGameInfo {
	return &ChessBaseGameInfo{
		GameLength:       params.GameLength,
		ATypicalStart:    params.ATypicalStart,
		IsEncoded:        params.IsEncoded,
		IsSpecialEncoded: params.IsSpecialEncoded,
		Is960:            params.Is960,
	}
}

type ChessBaseRecord struct {
	CBHRecord []byte
	CBP       []byte
	CBT       []byte
	CBG       []byte
}

type ChessBaseRecordParams struct {
	CBHRecord []byte
	CBP       []byte
	CBT       []byte
	CBG       []byte
}

func NewChessBaseRecord(params ChessBaseRecordParams) *ChessBaseRecord {
	return &ChessBaseRecord{
		CBHRecord: params.CBHRecord,
		CBP:       params.CBP,
		CBT:       params.CBT,
		CBG:       params.CBG,
	}
}

type Coord struct {
	X int // File
	Y int // Rank
}

type PieceInfo struct {
	PieceType int
	PieceNo   int
}

type Chessboard struct {
	Position  [8][8]PieceInfo
	PieceList [13][8]Coord
}

type ChessboardParams struct {
	Position  [8][8]PieceInfo
	PieceList [13][8]Coord
}

func InitialChessboard() *Chessboard {
	// we'll make PieceNo None = -1 because 0-7 represents a rank
	// position indicates the number of a type of piece with position in the [][] being the file and rank
	return &Chessboard{
		Position: [8][8]PieceInfo{
			{{W_ROOK, 0}, {W_PAWN, 0}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {B_PAWN, 0},
				{B_ROOK, 0}}, // A file 1-8
			{{W_KNIGHT, 0}, {W_PAWN, 1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {B_PAWN, 1},
				{B_KNIGHT, 0}}, // B file 1-8
			{{W_BISHOP, 0}, {W_PAWN, 2}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {B_PAWN, 2},
				{B_BISHOP, 0}}, // C file 1-8
			{{W_QUEEN, 0}, {W_PAWN, 3}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {B_PAWN, 3},
				{B_QUEEN, 0}}, // D file 1-8
			{{W_KING, -1}, {W_PAWN, 4}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {B_PAWN, 4},
				{B_KING, -1}}, // E file 1-8
			{{W_BISHOP, 1}, {W_PAWN, 5}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {B_PAWN, 5},
				{B_BISHOP, 1}}, // F file 1-8
			{{W_KNIGHT, 1}, {W_PAWN, 6}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {B_PAWN, 6},
				{B_KNIGHT, 1}}, // G file 1-8
			{{W_ROOK, 1}, {W_PAWN, 7}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {B_PAWN, 7},
				{B_ROOK, 1}}, // H file 1-8
		},
		// piece list gives the coordinates of the piece
		// -1 for empty squares
		PieceList: [13][8]Coord{
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // Empty Byte
			{{3, 0}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}},   // white queens
			{{1, 0}, {6, 0}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}},     // white knights
			{{2, 0}, {5, 0}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}},     // white bishops
			{{0, 0}, {7, 0}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}},     // white rooks
			{{3, 7}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}},   // black queens
			{{1, 7}, {6, 7}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}},     // black knights
			{{2, 7}, {5, 7}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}},     // black bishops
			{{0, 7}, {7, 7}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}},     // black rooks
			{{4, 0}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}},   // white king
			{{4, 7}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}},   // black king
			{{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {7, 1}},                 // white pawns
			{{0, 6}, {1, 6}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {7, 6}},                 // black pawns
		},
	}
}

func EmptyChessboard() *Chessboard {
	return &Chessboard{
		Position: [8][8]PieceInfo{
			{{EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1},
				{EMPTY, -1}}, // A file 1-8
			{{EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1},
				{EMPTY, -1}}, // B file 1-8
			{{EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1},
				{EMPTY, -1}}, // C file 1-8
			{{EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1},
				{EMPTY, -1}}, // D file 1-8
			{{EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1},
				{EMPTY, -1}}, // E file 1-8
			{{EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1},
				{EMPTY, -1}}, // F file 1-8
			{{EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1},
				{EMPTY, -1}}, // G file 1-8
			{{EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1}, {EMPTY, -1},
				{EMPTY, -1}}, // H file 1-8
		},
		PieceList: [13][8]Coord{
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // Empty Byte
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // white queens
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // white knights
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // white bishops
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // white rooks
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // black queens
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // black knights
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // black bishops
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // black rooks
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // white king
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // black king
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // white pawns
			{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // black pawns
		},
	}
}

func NewChessboard(params ChessboardParams) *Chessboard {
	return &Chessboard{
		Position:  params.Position,
		PieceList: params.PieceList,
	}
}

func (cb *Chessboard) Clone() *Chessboard {
	return &Chessboard{
		Position:  cb.Position,
		PieceList: cb.PieceList,
	}
}

type State struct {
	Chessboard  *Chessboard
	MoveNo      int
	IsWhiteTurn bool
}

type StateParams struct {
	Chessboard  *Chessboard
	MoveNo      int
	IsWhiteTurn bool
}

func NewState(params StateParams) *State {
	return &State{
		Chessboard:  params.Chessboard,
		MoveNo:      params.MoveNo,
		IsWhiteTurn: params.IsWhiteTurn,
	}
}
