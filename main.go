package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Print("Enter the name of the file you want to load:")

	var fileName string
	fmt.Scanf("%s", &fileName)

	content, err := parseAirBankStatement(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(content.AccountNumber)
	fmt.Println(content.StartDate)
	fmt.Println(content.EndDate)
	fmt.Print(len(content.Transactions))
}
