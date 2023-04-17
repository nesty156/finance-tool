package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	// Register a signal handler for the SIGHUP signal
	sighupCh := make(chan os.Signal, 1)
	signal.Notify(sighupCh, syscall.SIGHUP)

	fmt.Print("Enter the path to the directory of files:")
	var dirPath string
	fmt.Scanf("%s", &dirPath)

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: No two same files (checksum or smth)
	contents := []StatementOfAccount{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		content, err := parseAirBankStatement(filePath)
		if err != nil {
			log.Printf("Error loading file %s: %v", filePath, err)
			continue
		}

		fmt.Printf("Content of file %s:\n", filePath)
		fmt.Println(len(content.Transactions))
		contents = append(contents, content)
	}
	fmt.Println(len(contents))
	statement, err := mergeStatements(contents)
	if err != nil {
		log.Printf("Error creating statement: %v", err)
	}
	fmt.Println(statement.AccountNumber)
	fmt.Println(statement.StartDate)
	fmt.Println(statement.EndDate)
	fmt.Println(len(statement.Transactions))

	statement = sortTransactions(statement)
	value := sumTransactions(statement)
	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currnecy)

	saveToJson(statement)

	// Wait for the SIGHUP signal
	<-sighupCh

	fmt.Println("Exiting...")
	os.Exit(0)
}
