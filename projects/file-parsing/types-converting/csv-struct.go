package filesConverting

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CsvToStruct(textByte []byte) []Score {
	result := []Score{}
	s := Score{}

	r := csv.NewReader(strings.NewReader(string(textByte)))
	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read: %s", err)
	}
	nameIndex := -1
	highScoreIndex := -1
	i := 0
	for nameIndex == -1 || highScoreIndex == -1 || i < len(records[0]) {
		if records[0][i] == "name" {
			nameIndex = i
		} else if records[0][i] == "high score" {
			highScoreIndex = i
		}
		i++
	}

	for i := 1; i < len(records); i++ {
		highScore, err := strconv.Atoi(records[i][highScoreIndex])
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't parse the high score to type int: %s\n", err)
			continue
		}
		s.Name = records[i][nameIndex]
		s.HighScore = int32(highScore)
		result = append(result, s)
	}
	return result
}
