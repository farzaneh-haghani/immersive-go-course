package main

import (
	"bytes"
	"testing"
)

func TestBuffer(t *testing.T) {
	t.Run("If I make a buffer named b containing some bytes, calling b.Bytes() returns the same bytes I created it with.", func(t *testing.T) {
		text := "Hello"
		var b bytes.Buffer
		b.WriteString(text)
		got := b.Bytes()
		want := []byte(text)

		// if string(got)!=string(want){
		if !bytes.Equal(got, want) {
			t.Errorf("got %s but expect %s", got, want)
		}
	})

	t.Run("If I write some extra bytes to that buffer using b.Write(), a call to b.Bytes() returns both the initial bytes and the extra bytes.", func(t *testing.T) {
		text1 := "Hello"
		text2 := " world"
		var b bytes.Buffer
		b.WriteString(text1)
		b.WriteString(text2)
		got := b.Bytes()
		want := []byte(text1 + text2)

		if !bytes.Equal(got, want) {
			t.Errorf("got %s but expect %s", got, want)
		}
	})

	t.Run("If I call b.Read() with a slice big enough to read all of the bytes in the buffer, all of the bytes are read.", func(t *testing.T) {
		text := "Hello world"
		var b bytes.Buffer
		b.WriteString(text)
		s := make([]byte, len(text))
		b.Read(s)
		want := []byte(text)

		if !bytes.Equal(s, want) {
			t.Errorf("got %s but expect %s", s, want)
		}

	})

	t.Run("If I call b.Read() with a slice smaller than the contents of the buffer, some of the bytes are read. If I call it again, the next bytes are read.", func(t *testing.T) {
		text := "Hello world"
		var b bytes.Buffer
		b.WriteString(text)
		s := make([]byte, 6)
		b.Read(s)
		b.Read(s)
		want := []byte("world ")

		if !bytes.Equal(s, want) {
			t.Errorf("got %s but expect %s", s, want)
		}
	})
}
