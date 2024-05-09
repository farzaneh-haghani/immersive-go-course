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
	path := pathArr[0]
	isValid, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s: no such file or directory", path)
	}
	directory := isValid.IsDir()
	if directory {
		return fmt.Errorf("%s: is a directory", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%s: is not readable", path)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
