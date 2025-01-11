package filesConverting

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CsvToStruct(textByte []byte) []score {
	result := []score{}

	r := csv.NewReader(strings.NewReader(string(textByte)))
	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read: %s", err)
	}

	for i := 1; i < len(records); i++ {
		highScore, err := strconv.Atoi(records[i][1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't convert to number: %s", err)
		}
		s := score{Name: records[i][0], High_score: highScore}
		result = append(result, s)
	}
	return result
}
