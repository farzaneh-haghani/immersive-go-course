package main

import (
	"bytes"
	"testing"
)

func TestBuffer(t *testing.T) {
	bufferTest := []struct {
		name       string
		text1      string
		text2      string
		readLength int
		readOutput string
	}{
		{
			name:       "Write string on ourBuffer",
			text1:      "Hello",
			readLength: 3,
			readOutput: "lol",
		},
		{
			name:       "Write numbers on ourBuffer",
			text1:      "123",
			readLength: 3,
			readOutput: "123",
		},
		{
			name:       "Write nothing on ourBuffer",
			text1:      "",
			readLength: 0,
			readOutput: "",
		},
		{
			name:       "Write extra string",
			text1:      "Hello",
			text2:      " world",
			readLength: 3,
			readOutput: "lol",
		},
		{
			name:       "Write extra numbers",
			text1:      "Hello",
			text2:      "123",
			readLength: 3,
			readOutput: "lol",
		},
		{
			name:       "Add empty string to ourBuffer",
			text1:      "Hello",
			text2:      "",
			readLength: 3,
			readOutput: "lol",
		},
	}
	t.Run("Make OurByteBuffer named b containing some bytes, calling b.Bytes() returns the same bytes I created it with", func(t *testing.T) {
		for _, test := range bufferTest {
			var b OurByteBuffer
			b.OurWriteString(test.text1)
			got := b.OurBytes()
			want := []byte(test.text1)

			if !bytes.Equal(got, want) {
				t.Errorf("got %s but expect %s", got, want)
			}
		}
	})

	t.Run("Write some extra bytes to that buffer using b.Write(), a call to b.Bytes() returns both the initial bytes and the extra bytes", func(t *testing.T) {
		for _, test := range bufferTest {
			var b OurByteBuffer
			b.OurWriteString(test.text1)
			b.OurWriteString(test.text2)
			got := b.OurBytes()
			want := []byte(test.text1 + test.text2)

			if !bytes.Equal(got, want) {
				t.Errorf("got %s but expect %s", got, want)
			}
		}
	})

	t.Run("If I call b.Read() with a slice big enough to read all of the bytes in the buffer, all of the bytes are read.", func(t *testing.T) {
		for _, test := range bufferTest {
			var b OurByteBuffer
			b.OurWriteString(test.text1)
			s := make([]byte, len(test.text1))
			b.OurRead(s)
			want := []byte(test.text1)

			if !bytes.Equal(s, want) {
				t.Errorf("got %s but expect %s", s, want)
			}
		}
	})

	t.Run("Call b.Read() with a slice smaller than the contents of the buffer, some of the bytes are read. If I call it again, the next bytes are read.", func(t *testing.T) {
		for _, test := range bufferTest {
			var b OurByteBuffer
			b.OurWriteString(test.text1)
			s := make([]byte, test.readLength)
			b.OurRead(s)
			// s = make([]byte, test.readLength) I wanted to reset my buffer here but I got error "got `lo` but expect `lo`" ??
			b.OurRead(s)
			want := []byte(test.readOutput)

			if !bytes.Equal(s, want) {
				t.Errorf("got %s but expect %s", s, test.readOutput)
			}
		}
	})
}

func TestFilteringPipe(t *testing.T) {
	filteringPipeTest := []struct {
		name       string
		input      string
		wantOutput string
		length     int
	}{
		{
			name:       "mixing string and numbers",
			input:      "start=1, end=10",
			wantOutput: "start=, end=",
			length:     12,
		},
		{
			name:       "just numbers",
			input:      "123",
			wantOutput: "",
			length:     0,
		},
		{
			name:       "just string",
			input:      "Hello",
			wantOutput: "Hello",
			length:     5,
		},
	}

	t.Run("Implement write method for the interface", func(t *testing.T) {
		for _, test := range filteringPipeTest {
			var b bytes.Buffer
			filteringPipe := NewFilteringPipe(&b)
			length, _ := filteringPipe.Write([]byte(test.input))
			if b.String() != test.wantOutput || length != test.length {
				t.Errorf("got %s but expected %s", b.String(), test.wantOutput)
			}
		}
	})
}
