package main

import (
	"fmt"
	"os"
)

type score struct {
	Name       string `json:"name"`
	High_score int    `json:"high_score"`
}

func main() {
	files := []string{"examples/json.txt", "examples/data.csv"}
	for i, file := range files {
		text, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't read: %s", err)
			os.Exit(2)
		}
		switch i {
		case 0:
			result := jsonStruct(text)
			fmt.Println(result)
		case 1:
			result := csvStruct(text)
			fmt.Println(result)
		}

	}
}
