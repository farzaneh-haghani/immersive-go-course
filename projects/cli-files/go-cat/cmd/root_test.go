package cmd

import (
	"bytes"
	"testing"
)

func TestFlagHandling(t *testing.T) {

	flagTests := []struct {
		name       string
		numberLine bool
		help       bool
		fileNames  []string
		wantOutput string
	}{
		{
			name:       "Read the help file by using h flag",
			numberLine: false,
			help:       true,
			fileNames:  []string{},
			wantOutput: "\n",
		},
		{
			name:       "Read the help file without h flag",
			numberLine: false,
			help:       false,
			fileNames:  []string{"help.txt"},
			wantOutput: "     The options are as follows:\n\n     -b      Number the non-blank output lines, starting at 1.\n\n     -e      Display non-printing characters (see the -v option), and display a\n             dollar sign (‘$’) at the end of each line.\n\n     -l      Set an exclusive advisory lock on the standard output file\n             descriptor.  This lock is set using fcntl(2) with the F_SETLKW\n             command.  If the output file is already locked, cat will block\n             until the lock is acquired.\n",
		},
		{
			name:       "Read the help file with numbers and without h flag",
			numberLine: true,
			help:       false,
			fileNames:  []string{"help.txt"},
			wantOutput: "     1       The options are as follows:\n     2  \n     3       -b      Number the non-blank output lines, starting at 1.\n     4  \n     5       -e      Display non-printing characters (see the -v option), and display a\n     6               dollar sign (‘$’) at the end of each line.\n     7  \n     8       -l      Set an exclusive advisory lock on the standard output file\n     9               descriptor.  This lock is set using fcntl(2) with the F_SETLKW\n     10               command.  If the output file is already locked, cat will block\n     11               until the lock is acquired.\n",
		},
		{
			name:       "Can not read the file which doesn't exist",
			numberLine: true,
			help:       false,
			fileNames:  []string{"test.txt"},
			wantOutput: "open test.txt: no such file or directory",
		},
	}

	t.Run("Reading the files", func(t *testing.T) {
		for _, test := range flagTests {
			if test.help {
				test.wantOutput = string(helpFile) + test.wantOutput
			}
			var output bytes.Buffer
			got := flagHandler(&output, test.fileNames, test.numberLine, test.help)
			if output.String() != test.wantOutput && got.Error() != test.wantOutput {
				t.Errorf("got error %s and output %s but expected %s", got, output.String(), test.wantOutput)
			}
		}
	})
}
