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
		heap.Push(&trees, &HuffNode{tree1.Freq() + tree2.Freq(), tree1, tree2 })
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

	BuildTree(frequencyMap)
}
