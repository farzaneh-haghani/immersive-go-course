package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"
)

func sendRequest(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("We can't get them the weather!")
		fmt.Fprintln(os.Stderr, "Failed to initialize.!", err)
		os.Exit(2)
	}
	if resp == nil {
		fmt.Fprint(os.Stderr, "The server did not reply anything, which here is considered an error!")
		os.Exit(52)
	}
	return resp
}

func parseRetryAfterToNumber(retryAfter string) time.Duration {
	if sleepTime, err := time.ParseDuration(retryAfter + "s"); err != nil {
		return parseRetryAfterToString(retryAfter)
	} else {
		return sleepTime
	}
}

func parseRetryAfterToString(retryAfter string) time.Duration {
	sleepTime := time.Second
	if timeStamp, err := time.Parse(http.TimeFormat, retryAfter); err == nil {
		sleepTime = time.Until(timeStamp)
	}
	return sleepTime
}

func main() {
	var sleepTime time.Duration
	for {
		resp := sendRequest("http://localhost:8080")
		if resp.StatusCode == 200 {
			scanner := bufio.NewScanner(resp.Body)
			for i := 0; scanner.Scan(); i++ {
				fmt.Println(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				panic(err)
			}
			resp.Body.Close()
			os.Exit(0)
		} else if resp.StatusCode == 429 {
			// If I receive 429 but can't find ["Retry-After"] so I determine 5s to sleep
			// because you mentioned more than 5s give up.
			if resp.Header["Retry-After"] == nil {
				// I added with before sleep time to prevent sleeping long time if all before were less than 5s
				sleepTime += 5 * time.Second
			} else {
				retryAfter := resp.Header.Get("Retry-After")
				sleepTime += parseRetryAfterToNumber(retryAfter)
			}
			if sleepTime > time.Second && sleepTime < 5*time.Second {
				fmt.Println("Things may be a bit slow because we're doing a retry!")
				time.Sleep(sleepTime)
			} else {
				fmt.Println("We can't get them the weather!")
				resp.Body.Close()
				fmt.Fprintln(os.Stderr, "HTTP page not retrieved. Returned another error with the HTTP error code being 400 or above.")
				os.Exit(22)
			}
		} else {
			resp.Body.Close()
			fmt.Fprintln(os.Stderr, "Failed to initialize!")
			os.Exit(2)
		}
	}
}

// strconv.Atoi(str)                                 To convert a string representing an int number to int type
// strconv.ParseFloat(str, 64)                       To convert a string representing a float 64 number to float 64 type
// time.Duration(float64(time.Second) * sleepTime)   To convert float 64 to time.Duration
// time.Parse(http.TimeFormat, timeStr)              To convert time stamp (Thu, 25 Apr 2024 16:05:40 GMT ) to time.Time type (2024-04-25 16:05:40 +0000 UTC)
// time.Now().Add(time.Second)                       To add one second to current time with time.Time type
// time.Second                                       To return 1s with time.Duration type - we can multiple it by number of second we want
// endTime.Sub(startTime)                            To get time.Duration type from two times with time.Time type
// Time.Since(startTime)                             To get time.Duration type from a start time to now (time.Now())
// Time.Until(endTime)                               To get time.Duration type from now to an end time
// fmt.Printf("%#v", resp.Header)
