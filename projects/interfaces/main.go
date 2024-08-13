package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

type OurByteBuffer struct {
	myBuff []byte
}

// Is there a maximum about of data an OurByteBuffer can store?
// There is no maximum as far as our memory doesn't run out.

// Are there any important performance characteristics? e.g. is it much faster or slower to Write one large amount of data to it than write the same amount of data one byte at a time?
// yes, for buffer (bytes.buffer) is always faster because working with block of data.
// for OurByteBuffer which is a slice, it is better to append a large amount of data instead of one byte at a time to trigger the copy process less or create a slice with estimated size we need.

func (b *OurByteBuffer) Write(text string) (n int, err error) {
	b.myBuff = append(b.myBuff, []byte(text)...)
	return len(text), nil
}

func (b *OurByteBuffer) OurBytes() []byte {
	return b.myBuff
}

// Are there any important memory characteristics? e.g. does an OurByteBuffer always retain all data that was stored in it, or does it free some of its memory after it's been read?
// when we read from buffer/OurByteBuffer, the pointer move on and we don't access to data read but it is not free until all bytes read.
func (b *OurByteBuffer) Read(s []byte) (n int, err error) {
	num := copy(s, b.myBuff)
	b.myBuff = b.myBuff[num:]
	return num, nil
}

// What operations are safe or unsafe to perform concurrently on an OurByteBuffer from different threads?
// I am not sure. It won't be safe in all cases.
// Maybe it was ok if we write to the end of our buffer with one thread and read from start of it with another thread. but there is risk and it is better to finish one operation, then another.(async await)
// But It is unsafe to write to the same buffer with two different thread (or read as well).

type FilteringPipe struct {
	writer io.Writer
}

func NewFilteringPipe(w io.Writer) io.Writer {
	var f FilteringPipe
	f.writer = w
	return &f
}

func isNumber(data byte) bool {
	_, err := strconv.Atoi(string(data))
	return err == nil
}

func (f *FilteringPipe) Write(text []byte) (int, error) {
	// x:=make([]byte,0,len(text))
	first := 0
	length := 0
	for i, data := range text {
		if isNumber(data) {
			l, err := f.writer.Write(text[first:i])
			length += l
			if err != nil {
				return length, err
			}
			first = i + 1
			// x=append(x,data)   // It doing copy
		}
	}
	if first < len(text) {
		l, err := f.writer.Write(text[first:])
		length += l
		if err != nil {
			return l, err
		}
	}
	return length, nil
}

func main() {
	var b bytes.Buffer
	input := "start=1, end=10"
	// filteringPipe := FilteringPipe{writer: &b}
	filteringPipe := NewFilteringPipe(&b)
	_, err := filteringPipe.Write([]byte(input))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println(b.String())
}
