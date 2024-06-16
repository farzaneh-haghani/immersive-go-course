package cmd

import (
	"bytes"
	"testing"
)

func TestFlagHandling(t *testing.T) {

	flagTests := []struct {
		name       string
		comma      bool
		help       bool
		pathArr    []string
		wantOutput string
	}{
		{
			name:       "Listing root files without any flags",
			comma:      false,
			help:       false,
			pathArr:    []string{},
			wantOutput: "help.txt \troot.go \troot_test.go \t",
		},
		{
			name:       "Listing files for cmd and assets directories without any flags",
			comma:      false,
			help:       false,
			pathArr:    []string{"../cmd", "../assets"},
			wantOutput: "\n../cmd:\nhelp.txt \troot.go \troot_test.go \t\n../assets:\ndew.txt \tfor_you.txt \train.txt \t",
		},
		{
			name:       "Listing root files with comma between them",
			comma:      true,
			help:       false,
			pathArr:    []string{},
			wantOutput: "help.txt, \troot.go, \troot_test.go \t",
		},
		{
			name:       "Showing help",
			comma:      false,
			help:       true,
			pathArr:    []string{},
			wantOutput: "\n",
		},
		{
			name:       "Don't list any files for wrong directory",
			comma:      false,
			help:       false,
			pathArr:    []string{"test"},
			wantOutput: "stat test: no such file or directory",
		},
	}

	t.Run("Listing the files", func(t *testing.T) {
		for _, test := range flagTests {
			if test.help {
				test.wantOutput = string(helpFile) + test.wantOutput
			}
			var output bytes.Buffer
			got := flagHandler(&output, test.comma, test.help, test.pathArr)
			if output.String() != test.wantOutput && got.Error() != test.wantOutput {
				t.Errorf("got error %s and output %s but expected %s", got, output.String(), test.wantOutput)
			}
		}
	})
}
