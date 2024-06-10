package cmd

import (
	"bytes"
	"testing"
)

func TestExecute(t *testing.T) {
	flagTest := []struct {
		name     string
		comma    bool
		help     bool
		eachPath []string
		expected string
	}{
		{"Listing root files without any flags", false, false, []string{}, "help.txt        root.go         root_test.go"},
		{"Don't list any files for wrong directory", false, false, []string{"test"}, "test: no such file or directory"},
	}

	t.Run("Listing the files in the current path", func(t *testing.T) {
		// comma:=flag.Bool("m",false,"Add comma between the list")
		// help:=flag.Bool("h",false,"Show help")
		// flag.Set("m","false")
		// flag.Set("h","false")
		// flag.Parse()

		for _, e := range flagTest {
			var output bytes.Buffer
			flagHandling(&output, e.comma, e.help, e.eachPath)
			if output.String() != e.expected {
				t.Errorf("got %s", output.String())
			}
		}
	})
}
