package filesConverting

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type score struct {
	Name       string
	High_score int32
}

func JsonToStruct(textByte []byte) []score {
	result := []score{}
	s := score{}

	textTrimmed := strings.Trim(string(textByte), "[]\n")
	textSliced := strings.Split(string(textTrimmed), "\n")

	for _, text := range textSliced {
		text = strings.Trim(text, ",")
		if err := json.Unmarshal([]byte(text), &s); err != nil {
			fmt.Fprintf(os.Stderr, "can't unmarshal: %s", err)
			os.Exit(2)
		}
		result = append(result, s)
	}
	return result
}
