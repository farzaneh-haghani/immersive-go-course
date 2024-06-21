package main

type OurByteBuffer struct {
	myBuff []byte
}

func (b *OurByteBuffer) OurWriteString(text string) {
	b.myBuff = append(b.myBuff, []byte(text)...)
}

func (b *OurByteBuffer) OurBytes() []byte {
	return b.myBuff
}

func (b *OurByteBuffer) OurRead(s []byte){
	num:=copy(s,b.myBuff)
	b.myBuff=b.myBuff[num:]
}
