package lib

import (
	"errors"
	"io"
)

const (
	MaxVarintLen16 = 3
	MaxVarintLen32 = 5
	MaxVarintLen64 = 10
)

func UintToBuf(xx uint64) []byte {
	scratchBytes := make([]byte, MaxVarintLen64)
	nn := PutUvarint(scratchBytes, xx)
	return scratchBytes[:nn]
}

func PutUvarint(buf []byte, x uint64) int {
	i := 0
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return i + 1
}

func Uvarint(buf []byte) (uint64, int) {
	var x uint64
	var s uint
	for i, b := range buf {
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, -(i + 1) // overflow
			}
			return x | uint64(b)<<s, i + 1
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, 0
}

func IntToBuf(xx int64) []byte {
	scratchBytes := make([]byte, MaxVarintLen64)
	nn := PutVarint(scratchBytes, xx)
	return scratchBytes[:nn]
}

func PutVarint(buf []byte, x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return PutUvarint(buf, ux)
}

func Varint(buf []byte) (int64, int) {
	ux, n := Uvarint(buf) // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, n
}

var overflow = errors.New("binary: varint overflows a 64-bit integer")

func ReadUvarint(r io.Reader) (uint64, []byte) {
	var x uint64
	var s uint
	buf := []byte{0x00}
	history := []byte{}
	for i := 0; ; i++ {
		nn, err := io.ReadFull(r, buf)
		history = append(history, buf[0])
		if err != nil || nn != 1 {
			return x, history
		}
		b := buf[0]
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return x, history
			}
			return x | uint64(b)<<s, history
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

func ReadVarint(r io.Reader) (int64, []byte) {
	ux, history := ReadUvarint(r) // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, history
}
