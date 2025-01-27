package filesConverting

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func JsonToStruct(textByte []byte) []Score {
	result := []Score{}
	s := Score{}

	if string(textByte[0]) == "[" {
		if err := json.Unmarshal(textByte, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Can't marshal: %s\n", err)
		}
		return result
	}
	textTrimmed := strings.Trim(string(textByte), "[]\n")
	textSliced := strings.Split(string(textTrimmed), "\n")

	for _, text := range textSliced {
		if string(text[0]) == "#" {
			continue
		}
		if err := json.Unmarshal([]byte(text), &s); err != nil {
			fmt.Fprintf(os.Stderr, "can't unmarshal: %s", err)
		} else {
			result = append(result, s)
		}
	}
	return result
}
