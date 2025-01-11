package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func jsonStruct(text []byte) []score {

	textTrim := strings.Trim(string(text), "[]\n")
	textSlice := strings.Split(string(textTrim), "\n")

	result := []score{}
	s := score{}
	for _, v := range textSlice {
		v = strings.Trim(v, ",")
		if err := json.Unmarshal([]byte(v), &s); err != nil {
			fmt.Fprintf(os.Stderr, "can't unmarshal: %s", err)
			os.Exit(2)
		}
		result = append(result, s)
	}
	return result
}
