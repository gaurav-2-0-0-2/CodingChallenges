package main

import(
 "bufio"
 "fmt"
 "os"
 "log"
)

func main(){
	m := make(map[rune]int)
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Provide argument")
	}

	f, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for _,char := range(scanner.Text()){
			_, exists := m[char] 
			if exists {
				m[char] += 1
			}else{
				m[char] = 1
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for char, count := range m {
		fmt.Printf("%c: %d\n", char, count)
	}
}
