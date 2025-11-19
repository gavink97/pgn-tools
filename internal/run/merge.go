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

func Merge(args []string) {
	start := time.Now()
	defer func() {
		global.Logger.Info(fmt.Sprintf("merge took: %v\n", time.Since(start)))
	}()

	output := global.Output
	var inputs []string

	for i, arg := range args[1:] {
		if strings.EqualFold(arg, "--output") || strings.EqualFold(arg, "-o") {
			out := args[i+2]

			if parser.VerifyPGNOutput(out) {
				output = out
			} else {
				global.Logger.Error(fmt.Sprintf("Invalid output: %s", out))
				os.Exit(1)
			}

			continue
		}

		if strings.EqualFold(arg, output) {
			continue
		}

		stat, err := os.Stat(arg)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("An error occured getting stats on: %s", arg))
			global.Logger.Warn(err.Error())
			global.Logger.Info(fmt.Sprintf("Skipping file: %s", arg))
			continue
		}

		if stat.IsDir() {
			err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					global.Logger.Warn(fmt.Sprintf("An error occured accessing path: %s", path))
					global.Logger.Warn(err.Error())
					return nil
				}

				if info.IsDir() {
					return nil
				}

				if !parser.VerifyPGNInput(path) {
					global.Logger.Info(fmt.Sprintf("Skipping file: %s", arg))
					return nil
				}

				inputs = append(inputs, path)
				return nil
			})

			if err != nil {
				global.Logger.Warn(fmt.Sprintf("Error walking directory: %s", arg))
				global.Logger.Warn(err.Error())
			}

			continue
		}

		if !parser.VerifyPGNInput(arg) {
			global.Logger.Info(fmt.Sprintf("Skipping file: %s", arg))
			continue
		}

		inputs = append(inputs, arg)
	}

	for _, input := range inputs {
		games, err := parser.ParsePGN(input)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("An error occured merging %s", input))
			global.Logger.Warn(err.Error())
			continue
		}

		writer.WritePGN(output, games)
	}
}
