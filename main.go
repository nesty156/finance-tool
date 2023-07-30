package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nesty156/finance-tool/pkg/banks"
	"github.com/nesty156/finance-tool/pkg/bitcoin"
	"github.com/nesty156/finance-tool/pkg/converter"
	"github.com/nesty156/finance-tool/pkg/stocks"
	"github.com/nesty156/finance-tool/pkg/user"
	"github.com/nesty156/finance-tool/pkg/util"
)

const (
	airbank         = "a"
	revolut         = "r"
	moneta          = "m"
	ceskasporitelna = "cs"
	degiro          = "d"
	trading         = "212"
	trezor          = "t"
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
					ratesCZK := converter.GetConvertRatesCZK()
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
	options := []string{"[r] revolut", "[m] moneta", "[cs] ceska sporitelna", "[d] degiro", "[212] trading212", "[t] trezor"}
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
	case revolut:
		loadRevolut()
	case moneta:
		loadMoneta()
	case ceskasporitelna:
		loadCeskaSporitelna()
	case degiro:
		loadDegiro()
	case trading:
		loadTrading212()
	case trezor:
		loadTrezor()
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

func userInput() (filePath, accounName, currency string) {
	fmt.Print("Enter the path to the file: ")
	fmt.Scanln(&filePath)
	fmt.Print("Enter the name of the account: ")
	fmt.Scanln(&accounName)
	fmt.Print("Enter the currency of the account: ")
	fmt.Scanln(&currency)
	return filePath, accounName, currency
}

func loadAirBank() {
	filePath, accounName, currency := userInput()
	statement, err := banks.CreateAirBankStatement(filePath, accounName, currency)
	if err != nil {
		log.Printf("Error creating Airbank statement: %v", err)
		return
	}
	util.SaveSoaJson(statement)

	value := banks.SumTransactions(statement)
	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currency)

	saveStat(statement.AccountNumber, accounName, statement.Currency, value)
}

func loadRevolut() {
	filePath, accounName, currency := userInput()
	statement, err := banks.CreateRevolutStatement(filePath, accounName, currency)
	if err != nil {
		log.Printf("Error creating Airbank statement: %v", err)
		return
	}
	util.SaveSoaJson(statement)

	value := banks.SumTransactions(statement)
	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currency)

	saveStat(statement.AccountNumber, accounName, statement.Currency, value)
}

func loadMoneta() {
	fmt.Print("Enter the path to the file: ")
	var filePath string
	fmt.Scanln(&filePath)

	statement, err := banks.CreateMonetaStatement(filePath, "moneta")
	if err != nil {
		log.Printf("Error creating Moneta statement: %v", err)
		return
	}
	util.SaveSoaJson(statement)

	value := banks.SumTransactions(statement)
	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currency)

	saveStat(statement.AccountNumber, "moneta", statement.Currency, value)
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

	statement, err := banks.CreateCSStatement(jsonData)
	if err != nil {
		log.Printf("Error creating Ceska Sporitelna statement: %v", err)
		return
	}
	util.SaveSoaJson(statement)

	value := banks.SumTransactions(statement)
	fmt.Printf("Value of account %s is %.2f %s\n", statement.AccountNumber, value, statement.Currency)

	saveStat(statement.AccountNumber, "ceska sporitelna", statement.Currency, value)
}

func loadDegiro() {
	filePath, accounName, currency := userInput()

	portfolio, err := stocks.CreateDegiroPortfolio(filePath, accounName, currency)
	util.SavePortfolioJson(portfolio)
	if err != nil {
		log.Printf("Error creating Degiro portfolio: %v", err)
		return
	}

	value := stocks.PortfolioValue(portfolio)
	fmt.Printf("Value of your Degiro portfolio is %.2f %s\n", value, "EUR")

	saveStat(portfolio.Name, accounName, currency, value)
}

func loadTrading212() {
	filePath, accounName, currency := userInput()

	portfolio, err := stocks.CreateTrading212Portfolio(filePath, accounName, currency)
	if err != nil {
		log.Printf("Error creating Trading 212 fortfolio: %v", err)
		return
	}
	util.SavePortfolioJson(portfolio)

	value := stocks.PortfolioValue(portfolio)
	fmt.Printf("Value of your Trading 212 portfolio is %.2f %s\n", value, "EUR")

	saveStat(portfolio.Name, accounName, currency, value)
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
