package chessbase

import (
	"bytes"
	"fmt"
	"os"

	"github.com/gavink97/pgn-tools/internal/global"
)

func getPlayer(cbpFile []byte, playerNo int) (string, error) {
	var offset int

	switch cbpFile[0x18] {
	case 4:
		offset = 32 + (playerNo * 67)
	case 0:
		offset = 28 + (playerNo * 67)
	default:
		global.Logger.Error("Unknown CBP file version")
		os.Exit(1)
	}

	if len(cbpFile) < offset+59 {
		return "", fmt.Errorf("cbp too short: expected atleast %d bytes, got %d", offset+59, len(cbpFile))
	}

	lastNameByte := cbpFile[offset+9 : offset+9+30]
	firstNameByte := cbpFile[offset+39 : offset+39+20]

	lastName := string(bytes.TrimRight(lastNameByte, "\x00\xfe"))
	firstName := string(bytes.TrimRight(firstNameByte, "\x00\xfe"))

	return fmt.Sprintf("%s, %s", lastName, firstName), nil
}
