package main

import (
	"testing"
	"time"
)

func TestParseRetryAfterNumber(t *testing.T) {
	t.Run("Changing a number as a string type to time.Duration type", func(t *testing.T) {
		got := parseRetryAfterNumber("5")
		want := 5 * time.Second

		if got != want {
			t.Errorf("got %s want %d", got, want)
		}
	})

	t.Run("Returning 1s when there was an error in the value of Retry-After ", func(t *testing.T) {
		got := parseRetryAfterNumber("4 ")
		want := time.Second

		if got != want {
			t.Errorf("got %s want %d", got, want)
		}
	})
}
