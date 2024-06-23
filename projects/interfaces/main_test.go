package main

import (
	"bytes"
	"testing"
)

func TestBuffer(t *testing.T) {
	t.Run("Make OurByteBuffer buffer named b containing some bytes, calling b.Bytes() returns the same bytes I created it with.", func(t *testing.T) {
		text := "Hello"
		var b OurByteBuffer
		b.OurWriteString(text)
		got := b.OurBytes()
		want := []byte(text)

		if !bytes.Equal(got, want) {
			t.Errorf("got %s but expect %s", got, want)
		}
	})

	t.Run("Write some extra bytes to that buffer using b.Write(), a call to b.Bytes() returns both the initial bytes and the extra bytes.", func(t *testing.T) {
		text1 := "Hello"
		text2 := " world"
		var b OurByteBuffer
		b.OurWriteString(text1)
		b.OurWriteString(text2)
		got := b.OurBytes()
		want := []byte(text1 + text2)

		if !bytes.Equal(got, want) {
			t.Errorf("got %s but expect %s", got, want)
		}
	})

	t.Run("If I call b.Read() with a slice big enough to read all of the bytes in the buffer, all of the bytes are read.", func(t *testing.T) {
		text := "Hello world"
		var b OurByteBuffer
		b.OurWriteString(text)
		s := make([]byte, len(text))
		b.OurRead(s)
		want := []byte(text)

		if !bytes.Equal(s, want) {
			t.Errorf("got %s but expect %s", s, want)
		}
	})

	t.Run("Call b.Read() with a slice smaller than the contents of the buffer, some of the bytes are read. If I call it again, the next bytes are read.", func(t *testing.T) {
		text := "Hello world"
		var b OurByteBuffer
		b.OurWriteString(text)
		s := make([]byte, len(text)-5)
		b.OurRead(s)
		b.OurRead(s)
		want := []byte("world ")

		if !bytes.Equal(s, want) {
			t.Errorf("got %s but expect %s", s, want)
		}
	})
}
