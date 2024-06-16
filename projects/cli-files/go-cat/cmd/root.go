package cmd

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
)

//go:embed help.txt
var helpFile string

func Execute() error {
	pathArr, numberLine, help := ParsingFlags()

	err := flagHandling(os.Stdout, pathArr, numberLine, help)
	if err != nil {
		return err
	}
	return nil
}

func ParsingFlags() (pathArr []string, numberLine *bool, help *bool) {
	numberLine = flag.Bool("n", false, "Add number in front of each line")
	help = flag.Bool("h", false, "Show help")

	flag.Parse()
	pathArr = flag.Args()
	return
}

func flagHandling(w io.Writer, pathArr []string, numberLine *bool, help *bool) error {
	if *help {
		fmt.Fprintln(w, string(helpFile))
		return nil
	}

	if len(pathArr) == 0 && !*help {
		return fmt.Errorf("did not provide any path")
	}

	for _, path := range pathArr {
		err := OpenFile(w, path, *numberLine)
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenFile(w io.Writer, path string, numberLine bool) error {
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
		io.WriteString(w, suffix)
		io.WriteString(w, scanner.Text())
		io.WriteString(w, "\n")

	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
