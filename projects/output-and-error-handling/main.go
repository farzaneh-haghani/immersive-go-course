package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	for {
		resp, err := http.Get("http://localhost:8080")
		if err != nil {
			fmt.Println("Server not working at the moment!")
			fmt.Fprint(os.Stderr, "Server not working at the moment! ", err, "\n")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 {
			// fmt.Printf("\n%#v\n", resp.Header)
			fmt.Println("Please wait!")
			retryAfterStr := resp.Header["Retry-After"][0]
			// retryAfterNum, numErr := strconv.Atoi(retryAfterStr)
			retryAfterNum, numErr := strconv.ParseFloat(retryAfterStr, 64)

			var retryDuration time.Duration
			if numErr != nil {
				result, timeErr := time.Parse(http.TimeFormat, retryAfterStr)
				if timeErr != nil {
					fmt.Fprint(os.Stderr, "Retry after a while (1s): ", timeErr, "\n")
					result = time.Now().Add(time.Second)
				}else{
					fmt.Fprint(os.Stderr, "Retry after a time stamp: ", retryAfterStr, "\n")
				}
				retryDuration = -time.Since(result)
			} else {
				fmt.Fprint(os.Stderr, "Retry after ", retryAfterNum, "s \n")
				retryDuration = time.Duration(float64(time.Second) * retryAfterNum)
				fmt.Println(retryAfterNum,retryDuration)

			}
			// fmt.Println(retryDuration)
			time.Sleep(retryDuration)

		} else if resp.StatusCode == 200 {

			scanner := bufio.NewScanner(resp.Body)
			for i := 0; scanner.Scan() && i < 5; i++ {
				fmt.Println(scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				panic(err)
			}
			return
		}
	}
}
