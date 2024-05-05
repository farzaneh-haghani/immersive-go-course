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
		fmt.Println("We can't get the weather!")
		fmt.Fprintln(os.Stderr, "Failed to initialize.!", err)
		os.Exit(2)
	}
	if resp == nil {
		fmt.Fprintln(os.Stderr, "The server did not reply anything, which here is considered an error!")
		os.Exit(52)
	}
	return resp
}

func handleStatusCode(resp *http.Response, sleepTime *time.Duration) {
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		scanner := bufio.NewScanner(resp.Body)
		for i := 0; scanner.Scan(); i++ {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		os.Exit(0)
	} else if resp.StatusCode == 429 {
		// If I receive 429 but can't find ["Retry-After"] so I determine 1s to sleep. It act same as a while
		retryAfter := resp.Header.Get("Retry-After")
		// early return two possibility cases
		if retryAfter == "" || retryAfter == "a while" {
			// I added with the previous sleep times to prevent sleeping for too long and more than 5s give up.
			*sleepTime += time.Second
		} else {

			*sleepTime = parseRetryAfterNumber(retryAfter)
		}
		if *sleepTime > 5*time.Second {
			fmt.Println("We can't get the weather!")
			fmt.Fprintln(os.Stderr, "Sleep time is more than 5 second so HTTP page not retrieved. Returned another error with the HTTP error code being 400 or above.")
			os.Exit(22)
		} else if *sleepTime > time.Second {
			fmt.Println("Things may be a bit slow because we're doing a retry!")
		}
		time.Sleep(*sleepTime)
	} else {
		fmt.Fprintln(os.Stderr, "Failed to initialize!")
		os.Exit(2)
	}
}

func parseRetryAfterNumber(retryAfter string) (*time.Duration,err) {
	if sleepTime, err := time.ParseDuration(retryAfter + "s"); err != nil {
		return parseRetryAfterString(retryAfter),nil
	} else {
		return nil,err
	}
}

func parseRetryAfterString(retryAfter string) time.Duration {
	sleepTime := time.Second
	if timeStamp, err := time.Parse(http.TimeFormat, retryAfter); err == nil {
		sleepTime = time.Until(timeStamp)
	} else {
		fmt.Fprintln(os.Stderr, "There was an error in the value of Retry-After", err)
	}
	return sleepTime
}

func main() {
	var sleepTime time.Duration
	for {
		resp := sendRequest("http://localhost:8080")
		handleStatusCode(resp, &sleepTime)
	}
}

// strconv.Atoi(str)                                 To convert a string representing an int number to int type
// strconv.ParseFloat(str, 64)                       To convert a string representing a float 64 number to float 64 type
// time.Duration(number) * time.Second               To convert a number to time.Duration in second - without multiply it's in nanosecond
// time.Parse(http.TimeFormat, timeStr)              To convert time stamp (Thu, 25 Apr 2024 16:05:40 GMT ) to time.Time type (2024-04-25 16:05:40 +0000 UTC)
// time.ParseDuration("durationStr")                 To convert a string representing a duration to time.Duration
// time.Now().Add(time.Second)                       To add one second to current time with time.Time type
// time.Second                                       To return 1s with time.Duration type - we can multiple it by number of second we want
// endTime.Sub(startTime)                            To get time.Duration type from two times with time.Time type
// Time.Since(startTime)                             To get time.Duration type from a start time to now (time.Now())
// Time.Until(endTime)                               To get time.Duration type from now to an end time
// fmt.Printf("%#v", resp.Header)
