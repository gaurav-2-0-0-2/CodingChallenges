package main

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const HEADER_SEPARATER = "\nEND OF HEADER\n"
const HELP_MESSAGE = "Guide to use Squeez:\n" +
	"squeez [OPTION] [INPUT FILENAME] [OUTPUT FILENAME]\n" +
	"Available Options\n -h: Help\n -e: Encoding the file\n -o: Decoding the file in an output file\n"

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

func GenerateCodes(tree HuffTree, prefix []byte, encoder map[rune]string) map[rune]string {
	switch t := tree.(type) {
	case LeafNode:
		encoder[t.char] = string(prefix)
	case HuffNode:
		GenerateCodes(t.left_child, append(prefix, '0'), encoder)
		GenerateCodes(t.right_child, append(prefix, '1'), encoder)
	}
	return encoder
}

func WriteToFile(fileContent interface{}, filename string) {
	fmt.Println("fileContent to write to file: ", fileContent)
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

func packBitsIntoByte(bitstring string, outputFile string){
	var src []byte = []byte(bitstring)
	var dst []byte = make([]byte, (len(src) + 7) / 8)
	var bitMask byte = 1
	bitCounter := 0
	for b := 0; b < len(bitstring)/8; b++ {
		for bit := 0; bit < 8; bit++ {
			if bitCounter < len(src){
				dst[b] |= (src[bitCounter] & bitMask) << (7 - bit)
				bitCounter++
			}else{
				break
			}
		}
	}
	WriteToFile(dst, outputFile)
}

func EncodeFile(file *os.File, outputFile string, encoderMap map[rune]string) {
		var encodedData string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			for _, char := range scanner.Text() {
				encodedData += encoderMap[char]	
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		packBitsIntoByte(encodedData, outputFile)
}

func DecodeFile(outputFile string){

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

		fmt.Println("Frequency Map")
		for char, freq := range frequencyMap{
			fmt.Printf("%c: %d\n",char, freq)
		}

		huffManTree := BuildTree(frequencyMap)

		fmt.Printf("Encoder Map\n")
		encoderMap := GenerateCodes(huffManTree, []byte{}, make(map[rune]string))

		for char, prefixcode := range encoderMap{
			fmt.Printf("%c: %s\n",char, prefixcode)
		}
		
		WriteHeader(args[2], encoderMap)

		// Rewind the file for the second pass
		_, err = f.Seek(0, 0) // Reset file pointer to the beginning
		if err != nil {
			log.Fatal(err)
		}
		EncodeFile(f, args[2], encoderMap)
	case "-o":
		// Decoding part
		fmt.Println("this is an -o flag")
	case "-h":
		fmt.Printf(HELP_MESSAGE)
	default:
		fmt.Println("Unrecognized flag:", Flag)
		fmt.Println("run this command for help: squeez -h")
	}
}
