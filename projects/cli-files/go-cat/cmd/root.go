package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func Execute() error {

	numberLine := flag.Bool("n", false, "Add number in front of each line")
	help := flag.Bool("h", false, "Show help")
	helpFilePath := "assets/help.txt"

	flag.Parse()

	var pathArr []string

	if *numberLine && *help {
		pathArr = os.Args[3:]
	} else if !*numberLine && !*help {
		pathArr = os.Args[1:]
	} else {
		pathArr = os.Args[2:]
	}

	if *help {
		err := OpenFile(helpFilePath, false)
		if err != nil {
			return err
		}
	}

	if len(pathArr) == 0 && !*help {
		return fmt.Errorf("did not provide any path")
	}

	for _, path := range pathArr {
		err := OpenFile(path, *numberLine)
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenFile(path string, numberLine bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		if numberLine {
			text := fmt.Sprint("     ", i+1, "  ", scanner.Text()+"\n")
			os.Stdout.Write([]byte(text))
		} else {
			os.Stdout.Write([]byte(scanner.Text() + "\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
