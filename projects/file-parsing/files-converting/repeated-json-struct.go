package filesConverting

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func RepeatedJsonToStruct(textByte []byte) []score {
	result := []score{}
	s := score{}

	textTrimmed := strings.Trim(string(textByte), "\n")
	textSliced := strings.Split(string(textTrimmed), "\n")

	for _, text := range textSliced {
		if string(text[0]) != "#" {
			err := json.Unmarshal([]byte(text), &s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "can't unmarshal: %s", err)
			}
			result = append(result, s)
		}
	}
	return result
}
