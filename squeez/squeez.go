package main

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"io"
	"unicode/utf8"
)

const HEADER_SEPARATER = "\nEND OF HEADER\n"
const HELP_MESSAGE = "Guide to use Squeez:\n" +
	"squeez [OPTION] [INPUT FILENAME] [OUTPUT FILENAME]\n" +
	"Available Options\n -h: Help\n -e: Encoding the file\n -o: Decoding the file in an output file\n"

type Writer struct {
	data []byte
	current byte
	count uint
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

type HuffTree interface {
	Freq() int
}

type LeafNode struct {
	char rune
	freq int
}

type HuffNode struct {
	freq                    int
	left_child, right_child HuffTree
}

func (self LeafNode) Freq() int {
	return self.freq
}

func (self HuffNode) Freq() int {
	return self.freq
}

type PriorityQueue []HuffTree

func (pq PriorityQueue) Len() int { return len(pq) }

// This is a min heap: item with lowest priority comes first
// So when we call Pop() on pq we will get the lowest frequency character
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq() < pq[j].Freq()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(HuffTree))
}

func (pq *PriorityQueue) Pop() (x interface{}) {
	n := len(*pq)
	x = (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return
}

func BuildTree(m map[rune]int) HuffTree {
	trees := make(PriorityQueue, 0)
	for char, freq := range m {
		trees = append(trees, LeafNode{char, freq})
	}

	heap.Init(&trees)

	for trees.Len() > 1 {
		tree1 := heap.Pop(&trees).(HuffTree)
		tree2 := heap.Pop(&trees).(HuffTree)
		internalNode := HuffNode{tree1.Freq() + tree2.Freq(), tree1, tree2}
		heap.Push(&trees, internalNode)
	}
	return heap.Pop(&trees).(HuffTree)
}

func GenerateCodes(tree HuffTree, prefix []byte, encoder map[rune]string, count int) (int, map[rune]string){
	switch t := tree.(type) {
	case LeafNode:
		encoder[t.char] = string(prefix)
		count = count + 2
	case HuffNode:
		GenerateCodes(t.left_child, append(prefix, '0'), encoder, count)
		GenerateCodes(t.right_child, append(prefix, '1'), encoder, count)
	default: 
		count = count + 1
	}
	return count, encoder
}

func WriteToFile(fileContent interface{}, filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	switch v := fileContent.(type) {
	case string:
		_, write_err := f.WriteString(v)
		if write_err != nil {
			log.Fatal(err)
		}
	case []byte:
		_, err := f.Write(v)
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("fileContent must be of type string or []byte")
	}

	fmt.Println("Data written successfully to file")
}

func WriteHeader(filename string, encoderMap map[rune]string) {
	mapJSON := make(map[rune]string)
	for key, value := range encoderMap {
		mapJSON[key] = value
	}

	jsonData, err := json.Marshal(mapJSON)
	if err != nil {
		log.Fatal(err)
	}

	err_writing_header := os.WriteFile(filename, jsonData, 0644)
	if err_writing_header != nil {
		log.Fatal(err)
	}

	WriteToFile(HEADER_SEPARATER, filename)
}

func CountOccurences(file *os.File) map[rune]int {
		m := make(map[rune]int)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			for _, char := range scanner.Text() {
				m[char] += 1
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		return m
}

func EncodeFile(file *os.File, size int, outputFile string, encoderMap map[rune]string) {
	bitWriter := CreateBitWriter()

	_, fileErr := file.Seek(0, io.SeekStart)
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	err := encodeBits(bufio.NewReader(file), &bitWriter, encoderMap)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, utf8.RuneLen(rune(size)))
	utf8.EncodeRune(b, rune(size))

	data := make([]byte, 0)
	data = append(data, b...)
	data = append(data, bitWriter.Bytes()...)

	WriteToFile(data, outputFile)
}

func encodeBits(reader *bufio.Reader, bitWriter *Writer, encoderMap map[rune]string)error{
	for {
		r, _, err := reader.ReadRune()

		if err != nil {
			if err == io.EOF{
				break
			}
			log.Fatal(err)	
		}

		err = encodeBitCode(encoderMap, r, bitWriter)
		if err != nil {
			log.Fatal(err)	
		}
	}
	return nil
}

func encodeBitCode(encoderMap map[rune]string, r rune, bitWriter *Writer) (error) {
	for _, value := range encoderMap[r] {
		err := bitWriter.WriteBitFromChar(value)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func DecodeFile(file *os.File, encoderMap map[rune]string){

}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Error: Insufficient arguments.")
		fmt.Println("run this command for help: squeez -h")
		return
	}

	Flag := args[0]


	switch Flag {
	case "-e":
		if len(args) < 3 {
			fmt.Println("Error: Missing input or output file.")
			fmt.Println("Run this command for help: squeez -h")
			return
		}
		f, err := os.Open(args[1])
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		frequencyMap := CountOccurences(f)

		huffManTree := BuildTree(frequencyMap)

		count := 0
		size, encoderMap := GenerateCodes(huffManTree, []byte{}, make(map[rune]string), count)

		WriteHeader(args[2], encoderMap)
		EncodeFile(f, size, args[2], encoderMap)
	case "-o":
		fmt.Println("this is an -o flag")
	case "-h":
		fmt.Printf(HELP_MESSAGE)
	default:
		fmt.Println("Unrecognized flag:", Flag)
		fmt.Println("run this command for help: squeez -h")
	}
}
