package chessbase

import (
	"bytes"
	"fmt"
	"os"

	"github.com/gavink97/pgn-tools/internal/global"
)

func getEventSiteRounds(cbtFile []byte, tournamentNo int) ([]string, error) {
	var offset int

	switch cbtFile[0x18] {
	case 4:
		offset = 32 + (tournamentNo * 99)
	case 0:
		offset = 28 + (tournamentNo * 99)
	default:
		global.Logger.Error("Unknown CBT file version")
		os.Exit(1)
	}

	if len(cbtFile) < offset+99 {
		return nil, fmt.Errorf("cbt too short: expected atleast %d bytes, got %d", offset+99, len(cbtFile))
	}

	record := cbtFile[offset : offset+99]

	title := string(bytes.TrimRight(record[9:9+40], "\x00\xfe"))
	site := string(bytes.TrimRight(record[49:49+30], "\x00\xfe"))

	return []string{title, site}, nil
}
