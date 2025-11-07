package parser

import (
	"os"
	"strconv"
	"strings"

	"github.com/gavink97/pgn-tools/internal/types"
)

func NewPGNGame(str string) *types.Game {
	game := &types.Game{}

	var metadata []string
	var moves []string

	lines := strings.SplitSeq(str, "\n")

	for line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") {
			metadata = append(metadata, line)
		} else {
			moves = append(moves, line)
		}
	}

	for _, data := range metadata {
		if strings.HasPrefix(data, "[") && strings.HasSuffix(data, "]") {
			content := data[1 : len(data)-1]
			parts := strings.SplitN(content, " ", 2)

			if len(parts) == 2 {
				key := parts[0]
				value := strings.TrimSuffix(parts[1], "\"")
				value = strings.TrimPrefix(value, "\"")

				switch key {
				case "Event":
					game.Event = value
				case "Site":
					game.Site = value
				case "Date":
					game.Date = value
				case "Round":
					game.Round = value
				case "White":
					game.White = value
				case "Black":
					game.Black = value
				case "Result":
					game.Result = value
				case "BlackElo":
					elo, err := strconv.Atoi(value)
					if err != nil {
						elo = -1
					}
					game.BlackElo = elo
				case "ECO":
					game.ECO = value
				case "EventDate":
					game.EventDate = value
				case "WhiteElo":
					elo, err := strconv.Atoi(value)
					if err != nil {
						elo = -1
					}
					game.WhiteElo = elo
				case "Source":
					game.Source = value
				}
			}
		}
	}

	game.Game = strings.Join(moves, " ")

	return game
}

func ParsePGN(fileName string) ([]*types.Game, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	content := string(file)
	chunks := strings.Split(content, `[Event "`)

	var games []*types.Game
	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		gameText := `[Event "` + chunk
		game := NewPGNGame(gameText)
		games = append(games, game)
	}

	return games, nil
}
