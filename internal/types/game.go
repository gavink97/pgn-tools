package types

type Game struct {
	Event     string
	Site      string
	Date      string
	Round     string
	White     string
	Black     string
	Result    string
	BlackElo  int
	ECO       string
	EventDate string
	WhiteElo  int
	Source    string
	FEN       string
	Game      string
}

type GameParams struct {
	Event     string
	Site      string
	Date      string
	Round     string
	White     string
	Black     string
	Result    string
	BlackElo  int
	ECO       string
	EventDate string
	WhiteElo  int
	Source    string
	FEN       string
	Game      string
}

func NewGame(params GameParams) *Game {
	return &Game{
		Event:     params.Event,
		Site:      params.Site,
		Date:      params.Date,
		Round:     params.Round,
		White:     params.White,
		Black:     params.Black,
		Result:    params.Result,
		BlackElo:  params.BlackElo,
		ECO:       params.ECO,
		EventDate: params.EventDate,
		WhiteElo:  params.WhiteElo,
		Source:    params.Source,
		FEN:       params.FEN,
		Game:      params.Game,
	}
}
