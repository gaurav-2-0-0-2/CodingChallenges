package bitprocessor

import (
	"fmt"
)

type Writer struct {
	data    []byte
	current byte
	count   uint
}

func CreateBitWriter() Writer {
	return Writer{
		data:    make([]byte, 0),
		current: 0,
		count:   0,
	}
}

func (writer *Writer) WriteBitFromChar(bit rune) error {
	switch bit {
	case '1':
		writer.current = writer.current<<1 | 1
	case '0':
		writer.current = writer.current << 1
	default:
		return fmt.Errorf("Bit must be 0 or 1")
	}

	writer.count++

	if writer.count == 8 {
		writer.appendByte()
	}

	return nil
}

func (writer *Writer) appendByte() {
	writer.data = append(writer.data, writer.current)
	writer.current = byte(0)
	writer.count = 0
}

func (writer *Writer) Bytes() []byte {
	return writer.data
}
