package main

import "fmt"

type score struct {
	Name       string `json:"name"`
	High_score int    `json:"high_score"`
}

func main() {
	result := jsonStruct("examples/json.txt")
	fmt.Println(result)
	result = csvStruct("examples/data.csv")
}
