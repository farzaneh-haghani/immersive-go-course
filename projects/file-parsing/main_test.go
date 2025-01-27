package main

import (
	"bytes"
	"errors"
	filesConverting "file-parsing/types-converting"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToStruct(t *testing.T) {

	files := []string{"examples/json.txt", "examples/repeated-json.txt", "examples/data.csv", "examples/custom-binary-le.bin", "examples/custom-binary-be.bin"}
	var got []filesConverting.Score

	t.Run("Checking error when can't read a file", func(t *testing.T) {
		myError := errors.New("open amples/json.txt: no such file or directory")
		_, err := readingFiles("amples/json.txt")

		require.EqualError(t, err, myError.Error())
	})

	t.Run("Parsing different types to struct", func(t *testing.T) {
		for _, file := range files {
			textByte, _ := readingFiles(file)
			fileNameExtension := filepath.Ext(file)
			switch fileNameExtension {
			case ".txt":
				got = filesConverting.JsonToStruct(textByte)
			case ".csv":
				got = filesConverting.CsvToStruct(textByte)
			default:
				got = filesConverting.CustomBinaryToStruct(textByte)
			}

			require.Equal(t, filesConverting.Want(), got)
		}
	})

	t.Run("When no score provided", func(t *testing.T) {
		var b bytes.Buffer
		findMinMax(&b, []filesConverting.Score{})

		require.Equal(t, "No score provided\n", b.String())
	})

	t.Run("When just one score provided", func(t *testing.T) {
		var b bytes.Buffer
		findMinMax(&b, []filesConverting.Score{{Name: "Aya", HighScore: 10}})

		require.Equal(t, "Min is: Aya\n", b.String())
	})

	t.Run("Getting min and max score", func(t *testing.T) {
		var b bytes.Buffer
		findMinMax(&b, got)

		require.Equal(t, "\nMin is: Charlie \nMax is: Prisha\n", b.String())
	})

}
