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
	fileNames, numberLine, help := ParsingFlags()

	err := flagHandler(os.Stdout, fileNames, *numberLine, *help)
	if err != nil {
		return err
	}
	return nil
}

func ParsingFlags() (fileNames []string, numberLine *bool, help *bool) {
	numberLine = flag.Bool("n", false, "Add number in front of each line")
	help = flag.Bool("h", false, "Show help")

	flag.Parse()
	fileNames = flag.Args()
	return
}

func flagHandler(w io.Writer, fileNames []string, numberLine bool, help bool) error {
	if help {
		fmt.Fprintln(w, string(helpFile))
		return nil
	}

	if len(fileNames) == 0 && !help {
		return fmt.Errorf("did not provide any path")
	}

	for _, fileName := range fileNames {
		err := OutputFile(w, fileName, numberLine)
		if err != nil {
			return err
		}
	}
	return nil
}

func OutputFile(w io.Writer, fileName string, numberLine bool) error {
	file, err := os.Open(fileName)
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
