package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/nesty156/finance-tool/pkg/banks"
	"github.com/nesty156/finance-tool/pkg/bitcoin"
	"github.com/nesty156/finance-tool/pkg/stocks"
	"github.com/nesty156/finance-tool/pkg/user"
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
	login           = "log"
)

var (
	spaces = 12
)

func main() {
	// Register a signal handler for the SIGHUP signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)
	go handleSignal(sigs)

	logged := false
	var logOption string

	for {
		if !logged {
			logOption = "[log] log in"
		} else {
			logOption = "[log] log out"
		}
		options := []string{"[m] merge", logOption}

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
		case login:
			if !logged {
				user.Login()
				logged = true
				fmt.Println("Logged in as " + user.LoggedUser)

				filepath := user.LoggedUser + ".json"
				if _, err := os.Stat(filepath); !errors.Is(err, os.ErrNotExist) {
					account, err := util.LoadUserStatsJson(filepath)
					if err != nil {
						log.Printf("Error loading user stats: %v", err)
						return
					}
					ratesCZK := util.GetConvertRatesCZK()
					account.GetStatsInfo(user.ConvertRatesCZK{BTC: ratesCZK.BTC, EUR: ratesCZK.EUR, USD: ratesCZK.USD})
				}

			} else {
				user.Logout()
				logged = false
			}
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
	options := []string{"[t] trezor", "[m] moneta", "[d] degiro", "[cs] ceska sporitelna", "[212] trading212"}
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

func saveStat(name, component, currency string, value float64) {
	if user.LoggedUser != "" {
		var account user.AppAccount
		filepath := user.LoggedUser + ".json"
		if _, err := os.Stat(filepath); !errors.Is(err, os.ErrNotExist) {
			account, err = util.LoadUserStatsJson(filepath)
			if err != nil {
				log.Printf("Error loading user stats: %v", err)
				return
			}
		}

		account.Name = user.LoggedUser
		if account.RemoveStat(name) {
			fmt.Printf("Stat %s updated\n", name)
		}
		account.Stats = append(account.Stats, user.Stat{
			Name:       name,
			Component:  component,
			InsertDate: time.Now(),
			FilePath:   user.LoggedUser,
			Value:      value,
			Currency:   currency,
		})
		util.SaveUserStatsJson(account)
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

	*statement = banks.SortTransactions(*statement)
	value := banks.SumTransactions(*statement)

	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currency)
	util.SaveSoaJson(*statement)

	saveStat(statement.AccountNumber, "airbank", statement.Currency, value)
}

func loadTrezor() {
	fmt.Print("Enter the path to the directory of files: ")
	var dirPath string
	fmt.Scanln(&dirPath)

	stats, err := bitcoin.ConvertTrezorToStatement(dirPath)
	if err != nil {
		log.Printf("Error converting Trezor files: %v", err)
		return
	}

	fmt.Println("Conversion successful")

	for _, stat := range stats {
		saveStat(stat.Name, stat.Component, stat.Currency, stat.Value)
	}
}

func loadMoneta() {
	fmt.Print("Enter the path to the file: ")
	var filePath string
	fmt.Scanln(&filePath)

	statement, err := banks.ParseMonetaStatement(filePath, "moneta")
	if err != nil {
		log.Printf("Error parsing Moneta statement: %v", err)
		return
	}
	util.SaveSoaJson(statement)

	value := banks.SumTransactions(statement)
	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currency)

	saveStat(statement.AccountNumber, "moneta", statement.Currency, value)
}

func loadDegiro() {
	fmt.Print("Enter the path to the file: ")
	var filePath string
	fmt.Scanln(&filePath)

	portfolio, err := stocks.ParseDegiroPortfolio(filePath, "degiro")
	util.SavePortfolioJson(portfolio)
	if err != nil {
		log.Printf("Error parsing Degiro portfolio: %v", err)
		return
	}

	value := stocks.PortfolioValue(portfolio)
	fmt.Printf("Value of your Degiro portfolio is %.2f %s\n", value, "EUR")

	saveStat(portfolio.Name, "degiro", "EUR", value)
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

	statement, err := banks.ParseCeskaSporitelnaStatement(jsonData)
	if err != nil {
		log.Printf("Error parsing Ceska Sporitelna statement: %v", err)
		return
	}
	util.SaveSoaJson(statement)

	value := banks.SumTransactions(statement)
	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currency)

	saveStat(statement.AccountNumber, "ceska sporitelna", statement.Currency, value)
}

func loadTrading212() {
	fmt.Print("Enter the path to the file: ")
	var filePath string
	fmt.Scanln(&filePath)

	txs, err := stocks.ParseTrading212History(filePath)
	if err != nil {
		log.Printf("Error parsing Trading 212 history: %v", err)
		return
	}
	portfolio := stocks.TransactionsToPortfolio(txs, "trading212")
	util.SavePortfolioJson(portfolio)

	value := stocks.PortfolioValue(portfolio)
	fmt.Printf("Value of your Trading 212 portfolio is %.2f %s\n", value, "EUR")

	saveStat(portfolio.Name, "trading212", "EUR", value)
}

// Parse AirBank statement files and merge them
func mergeAirBankStatements(dirPath string) (*banks.StatementOfAccount, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	var contents []banks.StatementOfAccount
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		content, err := banks.ParseAirBankStatement(filePath)
		if err != nil {
			log.Printf("Error loading file %s: %v", filePath, err)
			continue
		}

		fmt.Printf("Content of file %s:\n", filePath)
		fmt.Println(len(content.Transactions))
		contents = append(contents, content)
	}

	statement, err := banks.MergeStatements(contents)
	if err != nil {
		return nil, fmt.Errorf("error creating statement: %v", err)
	}

	return &statement, nil
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
