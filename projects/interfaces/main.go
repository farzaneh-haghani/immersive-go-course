package main

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
)

func main() {
	// b:=bytes.NewBufferString("Hi")
	var b bytes.Buffer
	fmt.Printf("%#v\n", b)

	b.WriteString("test") // Add to the end of buffer - Write directly
	fmt.Printf("1%#v  %s\n", b, b.String())

	var w io.Writer = &b
	w.Write([]byte("test")) // Write through interface
	fmt.Printf("2%#v  %s\n", b, b.String())

	c := b.Bytes() // Read all bytes

	x := make([]byte, 4)
	l, _ := b.Read(x) // Read some bytes directly from from buffer so the pointer move on and those data is not accessible anymore.
	// b.UnreadByte()

	var r1 io.Reader = bytes.NewReader(c) // Make a copy from our buffer to don't change original buffer when reading
	l1, _ := r1.Read(x)

	var r2 io.Reader = &b // Read through interface but give buffer address to it so the pointer move on and those data is not accessible anymore.
	l2, _ := r2.Read(x)

	fmt.Println(b, b.Bytes(), c, x, l, l1, l2)

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	fmt.Printf("%d - ", m1.Frees)
	fmt.Printf("%d - ", m1.HeapAlloc)

	b.Reset()

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	fmt.Printf("%d - ", m2.Frees)
	fmt.Printf("%d - ", m2.HeapAlloc)

	runtime.GC() // running manually garbage collector

	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)
	fmt.Printf("%d - ", m3.Frees)     // free space increased
	fmt.Printf("%d - ", m3.HeapAlloc) // Heap allocation decrease due to frees buffer
}

// RAM: variables / data structure / buffer
// Stack: part of RAM - fast - allocation: Lifo - free: after finish the func
// Heap: part of RAM - slower than stack - allocation: dynamic - free: manually or garbage collector
// Array: fix size - allocate in compile time - not necessarily in stack
// Slice: getting size dynamically - allocate in run time(execute) - header will be in stack and data will be in heap specially when created by make()

// arr:=[3]int{1,2,3}  arr maybe in stack or heap
// s:=arr[1:4]  s in stack (because pointing to stack)
// x:=make([]int,3)  x in heap

// buffer: it is a structure data to store data as a block of bytes or string between I/O instead of reading byte by byte from hard disk
// one thread can read from buffer, another thread can write to the buffer
// input buffer: collect data from keyboard / hard disk / network and write on buffer
// output buffer: save before sending to screen / hard disk / network (e.g for testing)
// bytes: is a package working with buffer - using slice
// bytes.buffer: is a struct in this package with Write method
// io: is a package having some interfaces for I/O
// io.Writer: is an interface with a write method

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with [ErrTooLarge].

// Write to a buffer always uses this method, whether we use it directly or through interface

// var b bytes.Buffer

// func (b *Buffer) WriteString(s string) (n int, err error) {
// 	b.lastRead = opInvalid
// 	m, ok := b.tryGrowByReslice(len(s))
// 	if !ok {
// 		m = b.grow(len(s))
// 	}
// 	return copy(b.buf[m:], s), nil
// }

// This method is define a standard and uses it for sending to the right destination.
// var w io.Writer=&b
// func WriteString(w Writer, s string) (n int, err error) {
// 	if sw, ok := w.(StringWriter); ok {
// 		return sw.WriteString(s)
// 	}
// 	return w.Write([]byte(s))
// }
