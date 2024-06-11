package cmd

import (
	"bytes"
	"flag"
	"testing"
)

func TestExecute(t *testing.T) {
	comma := flag.Bool("m", false, "Add comma between the list")
	help := flag.Bool("h", false, "Show help")

	flagTests := []struct {
		flagM    string
		flagH    string
		name     string
		comma    *bool
		help     *bool
		eachPath []string
		expected string
	}{
		{"false", "false", "Listing root files without any flags", comma, help, []string{}, "help.txt \troot.go \troot_test.go \t"},
		{"true", "false", "Listing root files without any flags", comma, help, []string{}, "help.txt, \troot.go, \troot_test.go \t"},
		{"false", "true", "Listing root files with flag h", comma, help, []string{}, "\n"},
		{"false", "false", "Don't list any files for wrong directory", comma, help, []string{"test"}, "stat test: no such file or directory"},
	}

	t.Run("Listing the files in the current path", func(t *testing.T) {
		for _, test := range flagTests {
			flag.Set("m", test.flagM)
			flag.Set("h", test.flagH)
			flag.Parse()
			if *test.help == true {
				content, _ := helpFile.ReadFile("help.txt")
				test.expected = string(content) + test.expected
			}
			var output bytes.Buffer
			got := flagHandling(&output, test.comma, test.help, test.eachPath)
			if output.String() != test.expected && got.Error() != test.expected {
				t.Errorf("got %s but expected %s", output.String(), test.expected)
			}
		}
	})
}
