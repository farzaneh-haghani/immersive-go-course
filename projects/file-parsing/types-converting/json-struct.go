package filesConverting

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Score struct {
	Name       string
	High_Score int32
}

func JsonToStruct(textByte []byte) []Score {
	result := []Score{}
	s := Score{}

	textTrimmed := strings.Trim(string(textByte), "[]\n")
	textSliced := strings.Split(string(textTrimmed), "\n")

	for _, text := range textSliced {
		text = strings.Trim(text, ",")
		if string(text[0]) != "#" {
			if err := json.Unmarshal([]byte(text), &s); err != nil {
				fmt.Fprintf(os.Stderr, "can't unmarshal: %s", err)
				os.Exit(2)
			}
			result = append(result, s)
		}
	}
	return result
}
