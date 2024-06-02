package cmd

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"os"
)

//go:embed help.txt
var helpFile embed.FS

func Execute() error {

	numberLine := flag.Bool("n", false, "Add number in front of each line")
	help := flag.Bool("h", false, "Show help")

	flag.Parse()
	pathArr := flag.Args()

	if *help {
		content, err := helpFile.ReadFile("help.txt")
		if err != nil {
			return err
		}
		fmt.Println(string(content))
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
		var suffix string
		if numberLine {
			suffix = fmt.Sprint("     ", i+1, "  ")
		} else {
			suffix = ""
		}
		os.Stdout.WriteString(suffix)
		os.Stdout.WriteString(scanner.Text())
		os.Stdout.WriteString("\n")

	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
