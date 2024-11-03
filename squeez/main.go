package main

import(
 "unicode/utf8"
 "bufio"
 "fmt"
 "os"
 "log"
)

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

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println("count =", utf8.RuneCountInString(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
