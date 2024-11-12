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
	
	fmt.Printf("Initial Priority Queue:\n")
	for _, tree := range trees {
		switch t := tree.(type) {
		case LeafNode: 
			fmt.Printf("Leaf:%c Freq: %d\n", t.char, t.freq)
		}
	}

	for trees.Len()>1{
		tree1 := heap.Pop(&trees).(HuffTree)
		tree2 := heap.Pop(&trees).(HuffTree)
		internalNode := &HuffNode{tree1.Freq() + tree2.Freq(), tree1, tree2 }
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
		prefix = prefix[:len(prefix)-1]

		GenerateCodes(t.right_child, append(prefix, '1'), encoder)
		prefix = prefix[:len(prefix)-1]
	}
	return encoder
}

func main(){
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Provide argument")
	}

	f, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	frequencyMap := CountOccurrences(f)

	fmt.Printf("Frequency Map\n")
	for char, freq := range frequencyMap {
		fmt.Printf("%c: %d\n", char, freq)
	}

	huffManTree := BuildTree(frequencyMap)
	fmt.Printf("Encoder Map\n")
	encoderMap := GenerateCodes(huffManTree, []byte{}, make(map[rune]string))
	for char, prefixCode := range encoderMap {
		fmt.Printf("char: %c, value: %s\n", char, prefixCode)
	}
}
