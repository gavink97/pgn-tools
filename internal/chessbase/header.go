package chessbase

import (
	"encoding/binary"
	"fmt"
	"time"
)

func getRatings(cbhRecord []byte) ([]int, error) {
	if len(cbhRecord) < 35 {
		return []int{0, 0}, fmt.Errorf("cbhRecord too short: expected atleast 35 bytes, got %d", len(cbhRecord))
	}

	white := binary.BigEndian.Uint16([]byte{cbhRecord[31], cbhRecord[32]})
	black := binary.BigEndian.Uint16([]byte{cbhRecord[33], cbhRecord[34]})
	return []int{int(white), int(black)}, nil
}

func getRoundSubround(cbhRecord []byte) (string, error) {
	if len(cbhRecord) < 30 {
		return "", fmt.Errorf("cbhRecord too short: expected atleast 30 bytes, got %d", len(cbhRecord))
	}

	round := int(cbhRecord[29])
	subround := int(cbhRecord[30])

	if subround != 0 {
		return fmt.Sprintf("%d.%d", round, subround), nil
	}

	return fmt.Sprintf("%d", round), nil
}

func getResult(cbhRecord []byte) (string, error) {
	if len(cbhRecord) < 27 {
		return "", fmt.Errorf("cbhRecord too short: expected atleast 27 bytes, got %d", len(cbhRecord))
	}

	switch cbhRecord[27] {
	case 2:
		return "1-0", nil
	case 1:
		return "1/2-1/2", nil
	case 0:
		return "0-1", nil
	default:
		return "", nil
	}
}

func getDate(cbhRecord []byte) (string, error) {
	if len(cbhRecord) < 27 {
		return "", fmt.Errorf("cbhRecord too short: expected atleast 27 bytes, got %d", len(cbhRecord))
	}

	data := []byte{0, cbhRecord[24], cbhRecord[25], cbhRecord[26]}
	date := int(binary.BigEndian.Uint32(data))

	year := int((date & MASK_YEAR) >> 9)
	month := int((date & MASK_MONTH) >> 5)
	day := int(date & MASK_DAY)

	// fix if the date is incomplete, add ?? characters instead

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006.01.02"), nil
}

func getWhiteOffset(cbhRecord []byte) (int, error) {
	if len(cbhRecord) < 12 {
		return 0, fmt.Errorf("cbhRecord too short: expected atleast 12 bytes, got %d", len(cbhRecord))
	}

	data := []byte{0, cbhRecord[9], cbhRecord[10], cbhRecord[11]}
	white := binary.BigEndian.Uint32(data)
	return int(white), nil
}

func getBlackOffset(cbhRecord []byte) (int, error) {
	if len(cbhRecord) < 15 {
		return 0, fmt.Errorf("cbhRecord too short: expected atleast 15 bytes, got %d", len(cbhRecord))
	}

	data := []byte{0, cbhRecord[12], cbhRecord[13], cbhRecord[14]}
	black := binary.BigEndian.Uint32(data)
	return int(black), nil
}

func getTournamentOffset(cbhRecord []byte) (int, error) {
	if len(cbhRecord) < 18 {
		return 0, fmt.Errorf("cbhRecord too short: expected atleast 18 bytes, got %d", len(cbhRecord))
	}

	data := []byte{0, cbhRecord[15], cbhRecord[16], cbhRecord[17]}
	tournament := binary.BigEndian.Uint32(data)
	return int(tournament), nil
}

func getGameOffset(cbhRecord []byte) (int, error) {
	if len(cbhRecord) < 5 {
		return 0, fmt.Errorf("cbhRecord too short: expected atleast 5 bytes, got %d", len(cbhRecord))
	}

	tournament := binary.BigEndian.Uint32(cbhRecord[1:5])
	return int(tournament), nil
}

func isMarkedDeleted(cbhRecord []byte) bool {
	return int((MASK_MARKED_FOR_DELETION&int(cbhRecord[0]))>>7) == 1
}

func isGame(cbhRecord []byte) bool {
	return int(MASK_IS_GAME&int(cbhRecord[0])) == 1
}
