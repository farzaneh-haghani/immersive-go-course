package cmd

import (
	"embed"
	"flag"
	"fmt"
	"os"
)

//go:embed help.txt
var helpFile embed.FS

func Execute() error {
	comma, help, pathArr := ParsingFlags()

	err := flagHandling(comma, help, pathArr)
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

func flagHandling(comma *bool, help *bool, pathArr []string) error {
	if *help {
		content, err := helpFile.ReadFile("help.txt")
		if err != nil {
			return err
		}
		fmt.Println(string(content))
	}

	if len(pathArr) == 0 && !*help {
		path := "."
		err := DirectoryOutput(path, *comma)
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
			if *comma && i < len(pathArr)-1 {
				suffix = ", "
			} else {
				suffix = "  "
			}
			os.Stdout.WriteString(path)
			os.Stdout.WriteString(suffix)
		} else {

			if len(pathArr) > 1 {
				pathStr := "\n" + path + ":\n"
				os.Stdout.WriteString(pathStr) // situation: go-ls assets cmd
			}
			err := DirectoryOutput(path, *comma)
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

func DirectoryOutput(path string, comma bool) error {

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
		os.Stdout.WriteString(list[i].Name())
		os.Stdout.WriteString(suffix)
	}
	fmt.Print("\n")
	return nil
}
