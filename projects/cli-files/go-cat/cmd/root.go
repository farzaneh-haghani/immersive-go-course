package cmd

import (
	"bufio"
	"fmt"
	"os"
)

func Execute() error {
	pathArr := os.Args[1:]

	if len(pathArr) == 0 {
		return fmt.Errorf("did not provide any path")
	}

	var i int
	if pathArr[0] == "-n" {
		i = 1
	} else {
		i = 0
	}

	for i < len(pathArr) {
		path := pathArr[i]

		isPathValid, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("%s: no such file or directory", path)
		}

		isDirectory := isPathValid.IsDir()
		if isDirectory {
			return fmt.Errorf("%s: is a directory", path)
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("%s: is not readable", path)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for i := 0; scanner.Scan(); i++ {
			if pathArr[0] == "-n" {
				text := fmt.Sprint("     ", i+1, "  ", scanner.Text()+"\n")
				os.Stdout.Write([]byte(text))
			} else {
				os.Stdout.Write([]byte(scanner.Text() + "\n"))
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
		i++
	}
	return nil
}
