package main

import (
	"fmt"
	"go-cat/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(2)
	}
}
