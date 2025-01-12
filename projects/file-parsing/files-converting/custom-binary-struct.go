package filesConverting

import (
	"encoding/binary"
)

func CustomBinaryToStruct(textByte []byte) []score {
	result := []score{}
	s := score{}

	i := 2
	for i < len(textByte) {
		j := i + 4
		if textByte[0] == 255 && textByte[1] == 254 { // little endian byte order

			s.High_score = int32(binary.LittleEndian.Uint32(textByte[i:j]))
		}
		if textByte[0] == 254 && textByte[1] == 255 { // big endian byte order
			s.High_score = int32(binary.BigEndian.Uint32(textByte[i:j]))
		}
		b := make([]byte, len(textByte))
		for textByte[j] != 0 {
			b = append(b, textByte[j])
			j++
		}
		s.Name = string(b)
		i = j + 1
		result = append(result, s)
	}

	return result
}
