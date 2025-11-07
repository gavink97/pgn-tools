package main

import (
	"fmt"
	"os"

	"github.com/gavink97/pgn-tools/internal/global"
	"github.com/gavink97/pgn-tools/internal/parser"
	"github.com/gavink97/pgn-tools/internal/run"
)

func main() {
	args := os.Args[1:]
	global.SetDefaultLogger()
	parser.ParseArgs(args)

	program := args[0]

	switch program {
	case "convert":
	case "merge":
	case "query":
		run.Query(args)
	default:
		global.Logger.Error(fmt.Sprintf("invalid program: %s", program))
		os.Exit(1)
	}

}
