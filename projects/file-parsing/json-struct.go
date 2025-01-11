package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func jsonStruct(fileName string) []score {
	_, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open: %s", err)
		os.Exit(2)
	}
	text, err := os.ReadFile("examples/json.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read: %s", err)
		os.Exit(2)
	}

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
