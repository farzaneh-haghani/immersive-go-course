package main

import (
	"net/http"
	"testing"
	"time"
)

func TestParseRetryAfterString(t *testing.T) {
	t.Run("Changing a timestamp as a string to time.Duration - using parse ", func(t *testing.T) {
		timestampStr := time.Now().UTC().Add(time.Duration(4) * time.Second).Format(http.TimeFormat)
		timestampParse, _ := time.Parse(http.TimeFormat, timestampStr)
		got := parseRetryAfterString(timestampStr)
		want := time.Until(timestampParse)

		if got != want {
			t.Errorf("got %s want %d", got, want)
		}
	})

	t.Run("Changing a timestamp as a string to time.Duration - without pars", func(t *testing.T) {
		timestampStr := time.Now().UTC().Add(time.Duration(4) * time.Second).Format(http.TimeFormat)
		got := parseRetryAfterNumber(timestampStr)

		if got < 3*time.Second || got > 4*time.Second {
			t.Errorf("got %s want a time between 3 and 4 ", got)
		}
	})
}
