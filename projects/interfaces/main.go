package main

type OurByteBuffer struct {
	myBuff []byte
}

// Is there a maximum about of data an OurByteBuffer can store?
// There is no maximum as far as our memory doesn't run out.
// As far as I learned, buffer (bytes.buffer) can grow magically and always is the same buffer.
// But OurByteBuffer created as a slice bytes and can grow but always not the same slice bytes.

// Are there any important performance characteristics? e.g. is it much faster or slower to Write one large amount of data to it than write the same amount of data one byte at a time?
// yes, for buffer (bytes.buffer) is always faster because working with block of data.
// for OurByteBuffer which is a slice, it is better to append a large amount of data instead of one byte at a time to trigger the copy process less or create a slice with estimated size we need.

func (b *OurByteBuffer) OurWriteString(text string) {
	b.myBuff = append(b.myBuff, []byte(text)...)
}

func (b *OurByteBuffer) OurBytes() []byte {
	return b.myBuff
}

// Are there any important memory characteristics? e.g. does an OurByteBuffer always retain all data that was stored in it, or does it free some of its memory after it's been read?
// when we read from buffer/OurByteBuffer, the pointer move on and we don't access to data read so later garbage collector will frees them as there is no reference to them.
func (b *OurByteBuffer) OurRead(s []byte) {
	num := copy(s, b.myBuff)
	b.myBuff = b.myBuff[num:]
}

// What operations are safe or unsafe to perform concurrently on an OurByteBuffer from different threads?
// I am not sure. It won't be safe in all cases.
// Maybe it was ok if we write to the end of our buffer with one thread and read from start of it with another thread. but there is risk and it is better to finish one operation, then another.(async await)
// But It is unsafe to write to the same buffer with two different thread (or read as well).
