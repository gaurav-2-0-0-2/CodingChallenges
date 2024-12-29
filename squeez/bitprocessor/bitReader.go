package bitprocessor

type Reader struct {
	data   []byte
	offset int
	count  int
	length int
}

func CreateBitReader(data []byte) Reader {
	return Reader{data: data, length: len(data)}
}
