package cmd

import (
	"fmt"
	"io/fs"
	"os"
)

func Execute() error {
	pathArr := os.Args[1:]
	var list []fs.DirEntry
	var err error

	if len(pathArr) == 0 {
		list, err = os.ReadDir(".")
		if err != nil {
			return err
		}
		output(list)
		return nil
	}

	// situation: go-ls assets cmd
	for i := 0; i < len(pathArr); i++ {

		isPathValid, err := os.Stat(pathArr[i])
		if err != nil {
			return fmt.Errorf("%s no such file or directory", pathArr[i])
		}

		isDirectory := isPathValid.IsDir()
		if !isDirectory {
			os.Stdout.Write([]byte(pathArr[i] + "\t"))
		} else {
			list, err = os.ReadDir(pathArr[i])
			if err != nil {
				return err
			}

			if len(pathArr) > 1 {
				fmt.Println(pathArr[i] + ":")
			}
			output(list)
			if i < len(pathArr)-1 {
				fmt.Print("\n")
			}
		}
	}
	fmt.Print("\n")
	return nil
}

func output(list []fs.DirEntry) {
	for i := 0; i < len(list); i++ {
		fileName := fmt.Sprint(list[i].Name(), "\t")
		os.Stdout.Write([]byte(fileName))
	}
	fmt.Print("\n")
}
