package run

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gavink97/pgn-tools/internal/chessbase"
	"github.com/gavink97/pgn-tools/internal/global"
	"github.com/gavink97/pgn-tools/internal/parser"
	"github.com/gavink97/pgn-tools/internal/writer"
)

// pgn-tools convert [INPUT_PATH] [OUTPUT_PATH]

// should make a struct to hold chessbase byte arrays
// should use a struct to do the conversions
func Convert(args []string) {
	start := time.Now()
	defer func() {
		global.Logger.Info(fmt.Sprintf("conversion took: %v\n", time.Since(start)))
	}()

	input := args[1]

	// check for output for assign from --output
	output := args[2]

	if !chessbase.VerifyChessbaseInput(input) {
		global.Logger.Error(fmt.Sprintf("Invalid input: %s", input))
		os.Exit(1)
	}

	if !parser.VerifyPGNOutput(output) {
		global.Logger.Error(fmt.Sprintf("Invalid output: %s", output))
		os.Exit(1)
	}

	fileName := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))

	dir := filepath.Dir(input)

	cbhReader, err := chessbase.ReadMMap(fmt.Sprintf("%s/%s.cbh", dir, fileName))
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Unable to read MMAP: %s.cbh", fileName))
		os.Exit(1)
	}

	cbpReader, err := chessbase.ReadMMap(fmt.Sprintf("%s/%s.cbp", dir, fileName))
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Unable to read MMAP: %s.cbp", fileName))
		os.Exit(1)
	}

	cbtReader, err := chessbase.ReadMMap(fmt.Sprintf("%s/%s.cbt", dir, fileName))
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Unable to read MMAP: %s.cbt", fileName))
		os.Exit(1)
	}

	cbgReader, err := chessbase.ReadMMap(fmt.Sprintf("%s/%s.cbg", dir, fileName))
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Unable to read MMAP: %s.cbg", fileName))
		os.Exit(1)
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

	nrRecords := cbhReader.Read.Len() / 46

	games := 0
	errorMsgs := 0

	global.Logger.Info(fmt.Sprintf("converting %d chess games", nrRecords))

	for i := range nrRecords {
		cbhRecord := cbh[46*i : 46*(i+1)]

		record := chessbase.NewChessBaseRecord(chessbase.ChessBaseRecordParams{
			CBHRecord: cbhRecord,
			CBP:       cbp,
			CBT:       cbt,
			CBG:       cbg,
		})

		// write debugging errors to file
		game, err := record.ExtractGame()
		if err != nil {
			global.Logger.Warn(err.Error())
			errorMsgs++
			continue
		}

		writer.WritePGN(output, game)
		games++
	}

	if errorMsgs > 0 {
		global.Logger.Info(fmt.Sprintf("%d errors occured", errorMsgs))
	}

	global.Logger.Info(fmt.Sprintf("Extracted %d games from %s", games, input))
}
