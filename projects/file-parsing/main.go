package main

import (
	"bytes"
	filesConverting "file-parsing/types-converting"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	files := []string{"examples/json.txt", "examples/repeated-json.txt", "examples/data.csv", "examples/custom-binary-le.bin", "examples/custom-binary-be.bin"}
	var result []filesConverting.Score

	for _, file := range files {
		textByte, err := readingFiles(file)
		if err != nil {
			os.Exit(2)
		}
		fileNameExtension := filepath.Ext(file)

		switch fileNameExtension {
		case ".txt":
			result = filesConverting.JsonToStruct(textByte)
		case ".csv":
			result = filesConverting.CsvToStruct(textByte)
		default:
			result = filesConverting.CustomBinaryToStruct(textByte)
		}

		var b bytes.Buffer
		findMinMax(&b, result)
	}
}

func readingFiles(file string) ([]byte, error) {
	textByte, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read: %s", err)
		return nil, err
	}
	return textByte, nil
}

func findMinMax(w io.Writer, result []filesConverting.Score) {
	if len(result) == 0 {
		fmt.Fprintln(w, "No score provided")
		return
	} else if len(result) == 1 {
		fmt.Fprintf(w, "Min is: %v\n", result[0].Name)
		return
	}
	min := &result[0]
	max := &result[0]
	for i := 1; i < len(result); i++ {
		if result[i].HighScore < min.HighScore {
			min = &result[i]
		}
		if result[i].HighScore > max.HighScore {
			max = &result[i]
		}
	}
	fmt.Fprintf(w, "\nMin is: %s \nMax is: %s\n", min.Name, max.Name)
}
