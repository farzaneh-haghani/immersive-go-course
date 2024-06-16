package cmd

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
)

//go:embed help.txt
var helpFile string

func Execute() error {
	comma, help, pathArr := ParsingFlags()

	err := flagHandler(os.Stdout, *comma, *help, pathArr)
	if err != nil {
		return err
	}
	return nil
}

func ParsingFlags() (comma *bool, help *bool, pathArr []string) {
	comma = flag.Bool("m", false, "Add comma between the list")
	help = flag.Bool("h", false, "Show help")
	flag.Parse()
	pathArr = flag.Args()
	return
}

func flagHandler(w io.Writer, comma bool, help bool, pathArr []string) error {
	if help {
		fmt.Fprintln(w, string(helpFile))
		return nil
	}

	if len(pathArr) == 0 {
		path := "."
		err := DirectoryOutput(w, path, comma)
		if err != nil {
			return err
		}
	}

	for i, path := range pathArr {

		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		isDirectory := fileInfo.IsDir()
		if !isDirectory {
			var suffix string
			if comma && i < len(pathArr)-1 {
				suffix = ", "
			} else {
				suffix = "  "
			}
			io.WriteString(w, path)
			io.WriteString(w, suffix)

		} else {

			if len(pathArr) > 1 {
				pathStr := "\n" + path + ":\n"
				io.WriteString(w, pathStr) // situation: go-ls assets cmd
			}
			err := DirectoryOutput(w, path, comma)
			if err != nil {
				return err
			}
			if i < len(pathArr)-1 {
				fmt.Print("\n")
			}
		}
	}
	fmt.Print("\n")
	return nil
}

func DirectoryOutput(w io.Writer, path string, comma bool) error {
	list, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i := 0; i < len(list); i++ {
		var suffix string
		if comma && i < len(list)-1 {
			suffix = ", \t"
		} else {
			suffix = " \t"
		}
		io.WriteString(w, list[i].Name())
		io.WriteString(w, suffix)
	}
	fmt.Print("\n")
	return nil
}
