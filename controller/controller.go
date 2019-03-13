package controller

import "encoding/binary"

func GetScaling(whole, fract []byte) float32 {
	w := float32(binary.BigEndian.Uint16(whole))
	f := float32(binary.BigEndian.Uint16(fract)) / 65536
	return w + f
}

func Scale(factor float32, val []byte) float32 {
	return (float32(binary.BigEndian.Uint16(val)) * factor) / 32768
}
