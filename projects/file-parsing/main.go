package main

import (
	filesConverting "file-parsing/files-converting"
	"fmt"
	"os"
)

func main() {
	files := []string{"examples/json.txt", "examples/data.csv", "examples/repeated-json.txt", "examples/custom-binary-le.bin"}
	for i, file := range files {
		textByte, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't read: %s", err)
			os.Exit(2)
		}
		switch i {
		case 0:
			result := filesConverting.JsonToStruct(textByte)
			fmt.Println(result)
		case 1:
			result := filesConverting.CsvToStruct(textByte)
			fmt.Println(result)
		case 2:
			result := filesConverting.RepeatedJsonToStruct(textByte)
			fmt.Println(result)
		case 3:
			filesConverting.CustomBinaryToStruct(textByte)
		}
	}
}
