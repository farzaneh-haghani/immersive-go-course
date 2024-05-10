package main

import (
	"fmt"
	"go-ls/cmd"
	"os"
)

func main(){
	err:=cmd.Execute()
	if err!=nil{
		fmt.Fprintln(os.Stderr,err)
		os.Exit(2)
	}
}