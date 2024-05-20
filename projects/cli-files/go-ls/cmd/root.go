package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func Execute() error {

	comma := flag.Bool("m", false, "Add comma between the list")
	help := flag.Bool("h", false, "Show help")
	helpFilePath := "../go-cat/assets/help.txt"
	flag.Parse()

	var pathArr []string
	if *comma && *help {
		pathArr = os.Args[3:]
	} else if !*comma && !*help {
		pathArr = os.Args[1:]
	} else {
		pathArr = os.Args[2:]
	}

	if *help {
		err := OpenFile(helpFilePath)
		if err != nil {
			return err
		}
	}

	if len(pathArr) == 0 && !*help {
		path := "."
		err := DirectoryOutput(path, *comma)
		if err != nil {
			return err
		}
		return nil
	}

	for i, path := range pathArr {

		isPathValid, err := os.Stat(path)
		if err != nil {
			return err
		}

		isDirectory := isPathValid.IsDir()
		if !isDirectory {
			var suffix string
			if *comma && i < len(pathArr)-1 {
				suffix = ", "
			} else {
				suffix = "  "
			}
			os.Stdout.Write([]byte(path + suffix))
		} else {

			if len(pathArr) > 1 {
				os.Stdout.Write([]byte("\n" + path + ":\n")) // situation: go-ls assets cmd
			}
			err := DirectoryOutput(path, *comma)
			if err != nil {
				return err
			}
			if i < len(pathArr)-1 {
				fmt.Print("\n")
			}
		}
		i++
	}
	fmt.Print("\n")
	return nil
}

func DirectoryOutput(path string, comma bool) error {

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

func OpenFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	isPathValid, err := os.Stat(path)
	if err != nil {
		return err
	}

	isDirectory := isPathValid.IsDir()
	if isDirectory {
		return fmt.Errorf("%s: is a directory", path)
	}

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		os.Stdout.Write([]byte(scanner.Text() + "\n"))
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
