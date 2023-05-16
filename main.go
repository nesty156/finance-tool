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

// SIGHUP signal handler
func handleSignal(sigs chan os.Signal) {
	<-sigs
	fmt.Println("Exiting...")
	os.Exit(0)
}

// Prompt user for input
func promptUser() string {
	options := []string{"[t] trezor", "[m] moneta", "[d] degiro"}
	spaces := 12
	prompt := "Choose from [a] airbank\n"
	for _, option := range options {
		prompt += fmt.Sprintf("%*s%s\n", spaces, "", option)
	}
	prompt += "Enter your choice: "
	fmt.Print(prompt)
	var choice string
	fmt.Scanln(&choice)
	return choice
}

// Parse AirBank statement files and merge them
func mergeAirBankStatements(dirPath string) (*StatementOfAccount, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	var contents []StatementOfAccount
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

	statement, err := mergeStatements(contents)
	if err != nil {
		return nil, fmt.Errorf("error creating statement: %v", err)
	}

	return &statement, nil
}

// Convert Trezor files to Statement of Account format
func convertTrezorToStatement(dirPath string) (*StatementOfAccount, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	btcCZK, err := getBitcoinPrice("CZK")
	if err != nil {
		return nil, fmt.Errorf("error getting bitcoin price: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		account, err := parseBtcAccount(filePath)
		if err != nil {
			log.Printf("Error loading file %s: %v", filePath, err)
			continue
		}

		statement, err := account.convertToStatementOfAccount()
		if err != nil {
			log.Printf("Error converting bitcoin account to statement %s: %v", filePath, err)
			continue
		}

		fmt.Println(statement.AccountNumber)
		fmt.Println(statement.StartDate)
		fmt.Println(statement.EndDate)
		fmt.Println(len(statement.Transactions))

		value := sumTransactions(*statement)

		fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value*btcCZK, "CZK")
		saveToJson(*statement)
	}
	fmt.Scanln()

	return nil, nil
}

func main() {
	// Register a signal handler for the SIGHUP signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)
	go handleSignal(sigs)

	for {
		choice := promptUser()
		switch choice {
		case "a":
			fmt.Print("You chose AirBank. Enter the path to the directory of files: ")
			var dirPath string
			fmt.Scanf("%s", &dirPath)
			statement, err := mergeAirBankStatements(dirPath)
			if err != nil {
				log.Printf("Error merging statements: %v", err)
				continue
			}
			fmt.Println(statement.AccountNumber)
			fmt.Println(statement.StartDate)
			fmt.Println(statement.EndDate)
			fmt.Println(len(statement.Transactions))
			*statement = sortTransactions(*statement)
			value := sumTransactions(*statement)
			fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currnecy)
			saveToJson(*statement)
			fmt.Scanln()
		case "t":
			fmt.Print("You chose Trezor. Enter the path to the directory of files: ")
			var dirPath string
			fmt.Scanf("%s", &dirPath)
			_, err := convertTrezorToStatement(dirPath)
			if err != nil {
				log.Printf("Error converting Trezor files: %v", err)
				continue
			}
		case "m":
			fmt.Print("You chose Moneta. Enter the path to the file: ")
			var filePath string
			fmt.Scanf("%s", &filePath)
			xmlData, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file: %v", err)
				continue
			}
			statement, err := parseMonetaStatement(xmlData)
			if err != nil {
				log.Printf("Error parsing Moneta statement: %v", err)
				continue
			}
			saveToJson(statement)
			value := sumTransactions(statement)
			fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currnecy)
		case "d":
			fmt.Print("You chose Degiro. Enter the path to the file: ")
			var filePath string
			fmt.Scanf("%s", &filePath)
			csvData, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file: %v", err)
				continue
			}
			portfolio, err := parseDegiroPortfolio(csvData)
			if err != nil {
				log.Printf("Error parsing Degiro portfolio: %v", err)
				continue
			}
			value := portfolioValue(portfolio)
			fmt.Printf("Value of your degiro portfolio is %.2f %s\n", value, "EUR")
		default:
			fmt.Println("Invalid choice")
		}
	}
}
