package run

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gavink97/pgn-tools/internal/global"
	"github.com/gavink97/pgn-tools/internal/parser"
	"github.com/gavink97/pgn-tools/internal/writer"
)

func Query(args []string) {
	start := time.Now()
	defer func() {
		global.Logger.Info(fmt.Sprintf("query took: %v\n", time.Since(start)))
	}()

	input := args[1]
	keys := args[2]

	query, err := parser.ParseQuery(keys)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Error parsing query: %s", err))
		os.Exit(1)
	}

	ext := filepath.Ext(input)

	dir, file := filepath.Split(input)
	baseName := strings.TrimSuffix(file, ext)

	// name output based on query or set value?
	output := fmt.Sprintf("%s%s_modified%s", dir, baseName, ext)

	global.Logger.Debug(fmt.Sprintf("Modifying pgn at: %s", output))

	games, err := parser.ParsePGN(input)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Fatal Error: %v", err))
		os.Exit(1)
	}

	matches := 0

	for _, game := range games {
		match, err := query.Match(game)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Error evaluation game: %v", err))
			continue
		}

		if match {
			writer.WritePGN(output, game)
			matches++
		}
	}

	global.Logger.Info(fmt.Sprintf("Matched %d games out of %d", matches, len(games)))
}
