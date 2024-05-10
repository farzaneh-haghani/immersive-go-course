package cmd

import (
	"fmt"
	"os"
)

func Execute() error {
	pathArr := os.Args[1:]

	var path string
	var comma bool

	if len(pathArr) == 0 {
		path = "."
		directoryOutput(path, comma)
		return nil
	}

	if pathArr[0] == "-m" {
		comma = true
	} else {
		comma = false
	}

	if len(pathArr) == 1 && comma {
		path = "."
		directoryOutput(path, comma)
		return nil
	}

	var i int
	if comma {
		i = 1
	} else {
		i = 0
	}

	for i < len(pathArr) {

		isPathValid, err := os.Stat(pathArr[i])
		if err != nil {
			return fmt.Errorf("%s no such file or directory", pathArr[i])
		}

		isDirectory := isPathValid.IsDir()
		if !isDirectory {
			if comma && i < len(pathArr)-1 {
				os.Stdout.Write([]byte(pathArr[i] + ", "))
			} else {
				os.Stdout.Write([]byte(pathArr[i] + "  "))
			}
		} else {
			dirPath1 := len(pathArr) > 1 && !comma
			dirPath2 := len(pathArr) > 2 && comma
			if dirPath1 || dirPath2 {
				os.Stdout.Write([]byte("\n" + pathArr[i] + ":\n")) // situation: go-ls assets cmd
			}
			directoryOutput(pathArr[i], comma)
			if i < len(pathArr)-1 {
				fmt.Print("\n")
			}
		}
		i++
	}
	fmt.Print("\n")
	return nil
}

func directoryOutput(path string, comma bool) error {

	list, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i := 0; i < len(list); i++ {
		if comma && i < len(list)-1 {
			os.Stdout.Write([]byte(list[i].Name() + ",\t"))
		} else {
			os.Stdout.Write([]byte(list[i].Name() + " \t"))
		}
	}
	fmt.Print("\n")
	return nil
}
