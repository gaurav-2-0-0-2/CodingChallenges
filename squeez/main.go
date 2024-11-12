package main

import(
 "bufio"
 "fmt"
 "os"
 "log"
 "io"
 "container/heap"
)

type HuffTree interface {
	Freq() int
}

type LeafNode struct {
	char rune 
	freq int
}

type HuffNode struct {
	freq int
	left_child, right_child HuffTree
}

func (self LeafNode) Freq() int{
	return self.freq
}

func (self HuffNode) Freq() int{
	return self.freq
}

type PriorityQueue []HuffTree

func (pq PriorityQueue) Len() int {return len(pq)}

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

	for trees.Len()>1{
		tree1 := heap.Pop(&trees).(HuffTree)
		tree2 := heap.Pop(&trees).(HuffTree)
		internalNode := HuffNode{tree1.Freq() + tree2.Freq(), tree1, tree2 }
		heap.Push(&trees, internalNode)
	}
	return heap.Pop(&trees).(HuffTree)
}

func CountOccurrences(file io.Reader) map[rune]int {
	m := make(map[rune]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _,char := range(scanner.Text()){
				m[char] += 1
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return m
}

func GenerateCodes(tree HuffTree, prefix []byte, encoder map[rune]string) map[rune]string{
	switch t := tree.(type) {
	case LeafNode:
		encoder[t.char] = string(prefix)
	case HuffNode:
		GenerateCodes(t.left_child, append(prefix, '0'), encoder)
		GenerateCodes(t.right_child, append(prefix, '1'), encoder)
	}
	return encoder
}

func main(){
	args := os.Args[1:]

	Flag := args[0]
	//outputFile := args[2]

	switch Flag {
	case "-e":
		f, err := os.Open(args[1])
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		frequencyMap := CountOccurrences(f)

		huffManTree := BuildTree(frequencyMap)
		fmt.Printf("Encoder Map\n")
		encoderMap := GenerateCodes(huffManTree, []byte{}, make(map[rune]string))
		for char, prefixCodes := range encoderMap {
			fmt.Printf("%c: %s\n", char, prefixCodes)
		}
	case "-o":
		// Decoding part
		fmt.Println("this is an -o flag")
	case "-h":
		fmt.Println("Guide to use Squeez:")
		fmt.Println("squeez [OPTION] [INPUT FILENAME] [OUTPUT FILENAME]")
		fmt.Println("Available Options")
		fmt.Println("-h: Help")
		fmt.Println("-e: Encoding the file")
		fmt.Println("-o: Decoding the file in an output file")
	default:
		fmt.Println("not enough arguments")
		fmt.Println("run this command for help: squeez -h")
	}


}
