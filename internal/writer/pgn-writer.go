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

// make this FEN & if the header is empty not to include it
func formatPGN(game *types.Game) string {
	return fmt.Sprintf(`[Event "%s"]
[Site "%s"]
[Date "%s"]
[Round "%s"]
[White "%s"]
[WhiteElo "%d"]
[Black "%s"]
[BlackElo "%d"]
[Result "%s"]
[EventDate "%s"]
[ECO "%s"]
[Source "%s"]
%s

`, game.Event,
		game.Site,
		game.Date,
		game.Round,
		game.White,
		game.WhiteElo,
		game.Black,
		game.BlackElo,
		game.Result,
		game.EventDate,
		game.ECO,
		game.Source,
		game.Game)
}
