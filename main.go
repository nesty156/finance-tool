package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/nesty156/finance-tool/pkg/bitcoin"
	stat "github.com/nesty156/finance-tool/pkg/statement"
	"github.com/nesty156/finance-tool/pkg/stocks"
	"github.com/nesty156/finance-tool/pkg/util"
)

const (
	airbank         = "a"
	trezor          = "t"
	moneta          = "m"
	degiro          = "d"
	ceskasporitelna = "cs"
	trading         = "212"
	merge           = "m"
	load            = "l"
)

var (
	spaces = 12
)

func main() {
	// Register a signal handler for the SIGHUP signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)
	go handleSignal(sigs)

	options := []string{"[m] merge"}
	for {
		prompt := "Choose from [l] load\n"
		for _, option := range options {
			prompt += fmt.Sprintf("%*s%s\n", spaces, "", option)
		}
		prompt += "Enter your choice: "
		fmt.Printf(prompt)
		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case load:
			loadUser()
		case merge:
			mergePortfolios()
		default:
			fmt.Println("Invalid choice")
		}
	}
}

// SIGHUP signal handler
func handleSignal(sigs chan os.Signal) {
	<-sigs
	fmt.Println("Exiting...")
	os.Exit(0)
}

func loadUser() {
	options := []string{"[t] trezor", "[m] moneta", "[d] degiro", "[c] ceska sporitelna", "[212] trading212"}
	prompt := "Choose from [a] airbank\n"
	for _, option := range options {
		prompt += fmt.Sprintf("%*s%s\n", spaces, "", option)
	}
	prompt += "Enter your choice: "
	fmt.Printf(prompt)
	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case airbank:
		loadAirBank()
	case trezor:
		loadTrezor()
	case moneta:
		loadMoneta()
	case degiro:
		loadDegiro()
	case ceskasporitelna:
		loadCeskaSporitelna()
	case trading:
		loadTrading212()
	default:
		fmt.Println("Invalid choice")
	}
}

func loadAirBank() {
	fmt.Print("Enter the path to the directory of files: ")
	var dirPath string
	fmt.Scanln(&dirPath)

	statement, err := mergeAirBankStatements(dirPath)
	if err != nil {
		log.Printf("Error merging AirBank statements: %v", err)
		return
	}

	fmt.Println(statement.AccountNumber)
	fmt.Println(statement.StartDate)
	fmt.Println(statement.EndDate)
	fmt.Println(len(statement.Transactions))

	*statement = stat.SortTransactions(*statement)
	value := stat.SumTransactions(*statement)

	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currnecy)
	util.SaveSoaJson(*statement)
}

func loadTrezor() {
	fmt.Print("Enter the path to the directory of files: ")
	var dirPath string
	fmt.Scanln(&dirPath)

	_, err := convertTrezorToStatement(dirPath)
	if err != nil {
		log.Printf("Error converting Trezor files: %v", err)
		return
	}

	fmt.Println("Conversion successful")
}

func loadMoneta() {
	fmt.Print("Enter the path to the file: ")
	var filePath string
	fmt.Scanln(&filePath)

	xmlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}

	statement, err := stat.ParseMonetaStatement(xmlData)
	if err != nil {
		log.Printf("Error parsing Moneta statement: %v", err)
		return
	}
	util.SaveSoaJson(statement)

	value := stat.SumTransactions(statement)
	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currnecy)
}

func loadDegiro() {
	fmt.Print("Enter the path to the file: ")
	var filePath string
	fmt.Scanln(&filePath)

	csvData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}

	portfolio, err := stocks.ParseDegiroPortfolio(csvData, "degiro")
	util.SavePortfolioJson(portfolio)
	if err != nil {
		log.Printf("Error parsing Degiro portfolio: %v", err)
		return
	}

	value := stocks.PortfolioValue(portfolio)
	fmt.Printf("Value of your Degiro portfolio is %.2f %s\n", value, "EUR")
}

func loadCeskaSporitelna() {
	fmt.Print("Enter the path to the file: ")
	var filePath string
	fmt.Scanln(&filePath)

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}

	statement, err := stat.ParseCeskaSporitelnaStatement(jsonData)
	if err != nil {
		log.Printf("Error parsing Ceska Sporitelna statement: %v", err)
		return
	}
	util.SaveSoaJson(statement)

	value := stat.SumTransactions(statement)
	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currnecy)
}

func loadTrading212() {
	fmt.Print("Enter the path to the file: ")
	var filePath string
	fmt.Scanln(&filePath)

	csvData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}

	txs, err := stocks.ParseTrading212History(csvData)
	if err != nil {
		log.Printf("Error parsing Trading 212 history: %v", err)
		return
	}
	portfolio := stocks.TransactionsToPortfolio(txs, "trading212")
	util.SavePortfolioJson(portfolio)

	value := stocks.PortfolioValue(portfolio)
	fmt.Printf("Value of your Trading 212 portfolio is %.2f %s\n", value, "EUR")
}

// Parse AirBank statement files and merge them
func mergeAirBankStatements(dirPath string) (*stat.StatementOfAccount, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	var contents []stat.StatementOfAccount
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		content, err := stat.ParseAirBankStatement(filePath)
		if err != nil {
			log.Printf("Error loading file %s: %v", filePath, err)
			continue
		}

		fmt.Printf("Content of file %s:\n", filePath)
		fmt.Println(len(content.Transactions))
		contents = append(contents, content)
	}

	statement, err := stat.MergeStatements(contents)
	if err != nil {
		return nil, fmt.Errorf("error creating statement: %v", err)
	}

	return &statement, nil
}

// TODO: Move to pkg/bitcoin
// Convert Trezor files to Statement of Account format
func convertTrezorToStatement(dirPath string) (*stat.StatementOfAccount, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	btcCZK, err := util.GetBitcoinPrice("CZK")
	if err != nil {
		return nil, fmt.Errorf("error getting bitcoin price: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		account, err := bitcoin.ParseBtcAccount(filePath)
		if err != nil {
			log.Printf("Error loading file %s: %v", filePath, err)
			continue
		}

		statement, err := account.ConvertToStatementOfAccount()
		if err != nil {
			log.Printf("Error converting bitcoin account to statement %s: %v", filePath, err)
			continue
		}

		fmt.Println(statement.AccountNumber)
		fmt.Println(statement.StartDate)
		fmt.Println(statement.EndDate)
		fmt.Println(len(statement.Transactions))

		value := stat.SumTransactions(*statement)

		fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value*btcCZK, "CZK")
		util.SaveSoaJson(*statement)
	}

	return nil, nil
}

func mergePortfolios() {
	fmt.Print("Enter the first file: ")
	var filePath1 string
	fmt.Scanln(&filePath1)

	source, err := util.LoadPortfolioJson(filePath1)
	if err != nil {
		log.Printf("Error loading portfolio: %v", err)
		return
	}
	fmt.Println("File loaded")

	fmt.Print("Enter the second file: ")
	var filePath2 string
	fmt.Scanln(&filePath2)

	destination, err := util.LoadPortfolioJson(filePath2)
	if err != nil {
		log.Printf("Error loading portfolio: %v", err)
		return
	}
	fmt.Println("File loaded")

	merged := stocks.MergePortfolios(*source, *destination)
	util.SavePortfolioJson(merged)

	value := stocks.PortfolioValue(merged)
	fmt.Printf("Value of your merged portfolio is %.2f %s\n", value, "EUR")
}
