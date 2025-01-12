package main

import (
	filesConverting "file-parsing/types-converting"
	"fmt"
	"os"
)

func main() {
	files := []string{"examples/json.txt", "examples/repeated-json.txt", "examples/data.csv", "examples/custom-binary-le.bin", "examples/custom-binary-be.bin"}
	var result []filesConverting.Score

	for i, file := range files {
		textByte, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't read: %s", err)
			os.Exit(2)
		}

		switch i {
		case 0, 1:
			result = filesConverting.JsonToStruct(textByte)
		case 2:
			result = filesConverting.CsvToStruct(textByte)
		default:
			result = filesConverting.CustomBinaryToStruct(textByte)
		}
		findMinMax(result)
	}
}

func findMinMax(result []filesConverting.Score) {
	min := &result[0]
	max := &result[0]
	for i := 1; i < len(result); i++ {
		if result[i].High_Score < min.High_Score {
			min = &result[i]
		}
		if result[i].High_Score > max.High_Score {
			max = &result[i]
		}
	}
	fmt.Printf("\nMin is: %s \nMax is: %s\n", min.Name, max.Name)
}
