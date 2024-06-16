package cmd

import (
	"bytes"
	"flag"
	"testing"
)

func TestFlagHandling(t *testing.T) {
	numberLine := flag.Bool("n", false, "Add number in front of each line")
	help := flag.Bool("h", false, "Show help")

	flagTests := []struct {
		flagN      string
		flagH      string
		name       string
		numberLine *bool
		help       *bool
		fileNames  []string
		expect     string
	}{
		{"false", "true", "Read the help file by using h flag", numberLine, help, []string{}, "\n"},
		{"false","false","Read the help file without h flag",numberLine,help,[]string{"help.txt"},"     The options are as follows:\n\n     -b      Number the non-blank output lines, starting at 1.\n\n     -e      Display non-printing characters (see the -v option), and display a\n             dollar sign (‘$’) at the end of each line.\n\n     -l      Set an exclusive advisory lock on the standard output file\n             descriptor.  This lock is set using fcntl(2) with the F_SETLKW\n             command.  If the output file is already locked, cat will block\n             until the lock is acquired.\n"},
		{"true","false","Read the help file with numbers and without h flag",numberLine,help,[]string{"help.txt"},"     1       The options are as follows:\n     2  \n     3       -b      Number the non-blank output lines, starting at 1.\n     4  \n     5       -e      Display non-printing characters (see the -v option), and display a\n     6               dollar sign (‘$’) at the end of each line.\n     7  \n     8       -l      Set an exclusive advisory lock on the standard output file\n     9               descriptor.  This lock is set using fcntl(2) with the F_SETLKW\n     10               command.  If the output file is already locked, cat will block\n     11               until the lock is acquired.\n"},
	}

	t.Run("Reading the files", func(t *testing.T) {
		for _, test := range flagTests {
			flag.Set("n", test.flagN)
			flag.Set("h", test.flagH)
			flag.Parse()
			if *help {
				test.expect = string(helpFile) + test.expect
			}
			var output bytes.Buffer
			flagHandler(&output, test.fileNames, test.numberLine, test.help)
			if output.String() != test.expect {
				t.Errorf("got '%s' but expected '%s'", output.String(), test.expect)
			}
		}
	})
}
