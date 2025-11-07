package writer

import (
	"fmt"
	"log"
	"os"

	"github.com/gavink97/pgn-tools/internal/global"
	"github.com/gavink97/pgn-tools/internal/types"
)

func WritePGN(filename string, game *types.Game) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = f.Close()
		if err != nil {
			global.Logger.Error("an unexpected error occured closing file: %s\n%v", filename, err)
			os.Exit(1)
		}
	}()

	content := formatPGN(game)

	_, err = f.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
}

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
