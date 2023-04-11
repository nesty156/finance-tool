package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Print("Enter the name of the file you want to load:")

	var fileName string
	fmt.Scanf("%s", &fileName)

	file, err := os.Open(fileName) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(file)
}
