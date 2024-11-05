package main

import(
 "bufio"
 "fmt"
 "os"
 "log"
 "io"
 "container/heap"
)

type Item struct {
	key rune
	priority int
	index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {return len(pq)}

// This is a min heap: item with lowest priority comes first
// So when we call Pop() on pq we will get the lowest frequency character
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n 
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq	
	n := len(old)
	item := old[n-1]
	old[0] = old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0: n-1]
	heap.Fix(pq, 0)
	return item
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

	for char, count := range frequencyMap {
		fmt.Printf("%c: %d\n", char, count)
	}

	// build a binary tree from the occurrence table
	// minimum external path weight
	// if frequency of the letter is higher than the leaf node representing that letter will be on lesser depth
	// or a letter with high weight should have a low depth
	// weighted path length of a leaf: weight * depth 

	// Process: Building Huffman tree for n letters
	// creating a priority queue (min-heap)
	// lower frequency on the root node
	pq := make(PriorityQueue, len(frequencyMap))
	i := 0
	for char, frequency := range frequencyMap{
		pq[i] = &Item{
			key: char,
			priority: frequency,
			index: i,
		}
		i++
	}
	heap.Init(&pq)
	
	if pq.Len() > 0 {
		item1 := heap.Pop(&pq).(*Item)
		fmt.Printf("Removed first lowest frequency element: %c: %d\n", item1.key, item1.priority)
	}

}
