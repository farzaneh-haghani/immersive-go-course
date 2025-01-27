package filesConverting

import (
	"encoding/binary"
	"fmt"
	"os"
)

func CustomBinaryToStruct(textByte []byte) []Score {
	result := []Score{}
	s := Score{}

	i := 2
	if textByte[0] == 255 && textByte[1] == 254 {
		for i < len(textByte) {
			j := i + 4
			s.HighScore = int32(binary.LittleEndian.Uint32(textByte[i:j]))
			name := ""
			for textByte[j] != 0 {
				name += string(textByte[j])
				j++
			}
			s.Name = name
			i = j + 1
			result = append(result, s)
		}
	} else if textByte[0] == 254 && textByte[1] == 255 {
		for i < len(textByte) {
			j := i + 4
			s.HighScore = int32(binary.BigEndian.Uint32(textByte[i:j]))
			name := ""
			for textByte[j] != 0 {
				name += string(textByte[j])
				j++
			}
			s.Name = name
			i = j + 1
			result = append(result, s)
		}
	} else {
		fmt.Fprintln(os.Stderr, "It's not little endian or big endian to parse!")
	}
	return result
}
