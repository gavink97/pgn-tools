package writer

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gavink97/pgn-tools/internal/global"
	"github.com/gavink97/pgn-tools/internal/types"
)

func WritePGN(filename string, game any) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = f.Close()
		if err != nil {
			global.Logger.Error(fmt.Sprintf("an unexpected error occured closing file: %s", filename))
			global.Logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	var content string
	switch g := game.(type) {
	case *types.Game:
		content = formatPGN(g)
	case []*types.Game:
		var sb strings.Builder
		for _, game := range g {
			sb.WriteString(formatPGN(game))
		}
		content = sb.String()
	default:
		global.Logger.Error("unsupported type, expected *types.Game or []*types.Game")
		global.Logger.Error("type: %s", g)
		os.Exit(1)
	}

	_, err = f.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
}

func formatPGN(game *types.Game) string {
	out := fmt.Sprintf(`[Event "%s"]
[Site "%s"]
[Date "%s"]
[Round "%s"]
[White "%s"]
[Black "%s"]
[Result "%s"]`,
		game.Event,
		game.Site,
		game.Date,
		game.Round,
		game.White,
		game.Black,
		game.Result,
	)

	if game.WhiteElo != 0 {
		out += "\n"
		out += fmt.Sprintf(`[WhiteElo "%d"]`, game.WhiteElo)
	}

	if game.BlackElo != 0 {
		out += "\n"
		out += fmt.Sprintf(`[BlackElo "%d"]`, game.BlackElo)
	}

	if game.EventDate != "" {
		out += "\n"
		out += fmt.Sprintf(`[EventDate "%s"]`, game.EventDate)
	}

	if game.ECO != "" {
		out += "\n"
		out += fmt.Sprintf(`[ECO "%s"]`, game.ECO)
	}

	if game.FEN != "" {
		out += "\n"
		out += fmt.Sprintf(`[FEN "%s"]`, game.FEN)
		out += "\n"
		out += `[SetUp "1"]`
	}

	if game.Source != "" {
		out += "\n"
		out += fmt.Sprintf(`[Source "%s"]`, game.Source)
	}

	out += fmt.Sprintf("\n%s\n\n", game.Game)
	return out
}
