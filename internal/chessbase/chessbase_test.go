package chessbase

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gavink97/pgn-tools/internal/global"
	"github.com/gavink97/pgn-tools/internal/types"
)

func TestMain(m *testing.M) {
	global.SetDefaultLogger()
	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestConvertMain(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Unable to locate user home directory: %v", err)
	}

	base := "~/Downloads/SetupMega2020/Bases/MegaDatabase2020/Mega Database 2020.cbh"
	input := filepath.Join(homeDir, base[2:])

	if !VerifyChessbaseInput(input) {
		t.Errorf("Invalid input: %s", input)
	}

	fileName := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
	dir := filepath.Dir(input)

	cbhReader, err := ReadMMap(fmt.Sprintf("%s/%s.cbh", dir, fileName))
	if err != nil {
		t.Errorf("Unable to read MMAP: %s.cbh", fileName)
	}

	cbpReader, err := ReadMMap(fmt.Sprintf("%s/%s.cbp", dir, fileName))
	if err != nil {
		t.Errorf("Unable to read MMAP: %s.cbp", fileName)
	}

	cbtReader, err := ReadMMap(fmt.Sprintf("%s/%s.cbt", dir, fileName))
	if err != nil {
		t.Errorf("Unable to read MMAP: %s.cbt", fileName)
	}

	cbgReader, err := ReadMMap(fmt.Sprintf("%s/%s.cbg", dir, fileName))
	if err != nil {
		t.Errorf("Unable to read MMAP: %s.cbg", fileName)
	}

	defer func() {
		err := cbhReader.Read.Close()
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("An error occured closing %s.cbh", fileName))
			global.Logger.Warn(err.Error())
		}

		err = cbpReader.Read.Close()
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("An error occured closing %s.cbp", fileName))
			global.Logger.Warn(err.Error())
		}

		err = cbtReader.Read.Close()
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("An error occured closing %s.cbt", fileName))
			global.Logger.Warn(err.Error())
		}

		err = cbgReader.Read.Close()
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("An error occured closing %s.cbg", fileName))
			global.Logger.Warn(err.Error())
		}
	}()

	cbh := cbhReader.File
	cbp := cbpReader.File
	cbt := cbtReader.File
	cbg := cbgReader.File

	headerByte := cbh[0:46]
	headerId := headerByte[0:6]

	switch string(headerId) {
	case "\x00\x00\x2c\x00\x2e\x01":
		global.Logger.Info("created by CB9+")
	case "\x00\x00\x24\x00\x2e\x01":
		global.Logger.Info("created by Chess Program X/CB Light")
	default:
		global.Logger.Info(fmt.Sprintf("unknown header id: %s", string(headerId)))
	}

	i := 7510212
	cbhRecord := cbh[46*i : 46*(i+1)]

	record := NewChessBaseRecord(ChessBaseRecordParams{
		CBHRecord: cbhRecord,
		CBP:       cbp,
		CBT:       cbt,
		CBG:       cbg,
	})

	result, err := record.ExtractGame()
	if err != nil {
		t.Errorf("An error occured extracting game data from record")
	}

	res := fmt.Sprintf(
		`Event: %s
Site: %s
Date: %s
Round: %s
White: %s
Black: %s
Result: %s
BlackElo: %d
ECO: %s
EventDate: %s
WhiteElo: %d
Source: %s
FEN: %s
Game: %s
`,
		result.Event,
		result.Site,
		result.Date,
		result.Round,
		result.White,
		result.Black,
		result.Result,
		result.BlackElo,
		result.ECO,
		result.EventDate,
		result.WhiteElo,
		result.Source,
		result.FEN,
		result.Game,
	)

	fmt.Print(res)

	expected := types.NewGame(types.GameParams{
		Event:     "Moscow Aeroflot op-A 17th",
		Site:      "Moscow",
		Date:      "2018.02.21",
		Round:     "2",
		White:     "Andreikin, Dmitry",
		WhiteElo:  2712,
		Black:     "Vavulin, Maksim",
		BlackElo:  2575,
		Result:    "1/2-1/2",
		EventDate: "2018.02.21",
		ECO:       "",
		Source:    "",
		FEN:       "",
		//Game:      "1. d4 Nf6 2. c4 e6 3. Nc3 Bb4 4. Qc2 d5 5. a3 Bxc3+ 6. Qxc3 dxc4 7. Qxc4 b6 8. Bg5 Ba6 9. Qc3 Qd5 10. Bxf6 gxf6 11. f3 Nd7 12. Rc1 c5 13. e4 Qb7 14. dxc5 bxc5 15. Bxa6 Qxa6 16. Ne2 Rg8 17. Kf2 Rb8 18. Rc2 Ne5 19. Rd1 c4 20. Qd4 Kf8 21. Nf4 Qb6 22. Kf1 Ke7 23. Qxb6 1/2-1/2",
		Game: "1. d4 Nf6 2. c4 e6 3. Nc3 Bb4 4. Qc2 d5 5. a3 Bxc3 6. Qxc3 dxc4 7. Qxc4 b6 8. Bg5 Ba6 9. Qc3 Qd5 10. Bxf6 gxf6 11. f3 Nd7 12. Rc1 c5 13. e4 Qb7 14. dxc5 bxc5 15. Bxa6 Qxa6 16. Ne2 Rg8 17. Kf2 Rb8 18. Rc2 Ne5 19. Rd1 c4 20. Qd4 Kf8 21. Nf4 Qb6 22. Kf1 Ke7 23. Qxb6 1/2-1/2",
	})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Incorrect Result: \nresult: %v \nexpected: %v", result, expected)
	}
}

func TestIsWhiteTurnFEN(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	result, err := IsWhiteTurnFEN(fen)
	if err != nil {
		t.Errorf("unexpected character in FEN")
	}

	if !result {
		t.Errorf("Incorrect Result: \nresult: %v \nexpected: %v", result, true)
	}
}
