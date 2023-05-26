package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type TradingTransaction struct {
	action       string
	time         string
	isin         string
	ticker       string
	name         string
	shares       float64
	price        float64
	currency     string
	exchangeRate float64
	total        float64
	chargeAmount float64
	notes        string
	id           string
	fee          float64
}

func parseTrading212History(csvData []byte) ([]TradingTransaction, error) {
	reader := csv.NewReader(strings.NewReader(string(csvData)))
	reader.Comma = ','

	var history []TradingTransaction

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading record:", err)
			continue
		}

		if record[0] == "Market buy" {
			shares, err := strconv.ParseFloat(record[5], 64)
			if err != nil {
				fmt.Println("Error parsing shares:", err)
				continue
			}

			price, err := strconv.ParseFloat(record[6], 64)
			if err != nil {
				fmt.Println("Error parsing price:", err)
				continue
			}

			exchangeRate, err := strconv.ParseFloat(record[8], 64)
			if err != nil {
				fmt.Println("Error parsing exchange rate:", err)
				continue
			}

			total, err := strconv.ParseFloat(record[10], 64)
			if err != nil {
				fmt.Println("Error parsing total:", err)
				continue
			}

			history = append(history, TradingTransaction{
				action:       record[0],
				time:         record[1],
				isin:         record[2],
				ticker:       record[3],
				name:         record[4],
				shares:       shares,
				price:        price,
				currency:     record[7],
				exchangeRate: exchangeRate,
				total:        total,
				chargeAmount: 0,
				notes:        record[12],
				id:           record[13],
				fee:          0,
			})
		} else if record[0] == "Market sell" {
			shares, err := strconv.ParseFloat(record[5], 64)
			if err != nil {
				fmt.Println("Error parsing shares:", err)
				continue
			}

			price, err := strconv.ParseFloat(record[6], 64)
			if err != nil {
				fmt.Println("Error parsing price:", err)
				continue
			}

			exchangeRate, err := strconv.ParseFloat(record[8], 64)
			if err != nil {
				fmt.Println("Error parsing exchange rate:", err)
				continue
			}

			total, err := strconv.ParseFloat(record[10], 64)
			if err != nil {
				fmt.Println("Error parsing total:", err)
				continue
			}

			history = append(history, TradingTransaction{
				action:       record[0],
				time:         record[1],
				isin:         record[2],
				ticker:       record[3],
				name:         record[4],
				shares:       -shares,
				price:        price,
				currency:     record[7],
				exchangeRate: exchangeRate,
				total:        total,
				chargeAmount: 0,
				notes:        record[12],
				id:           record[13],
				fee:          0,
			})
		} else if record[0] == "Deposit" {
			total, err := strconv.ParseFloat(record[11], 64)
			if err != nil {
				fmt.Println("Error parsing total:", err)
				continue
			}

			history = append(history, TradingTransaction{
				action:       record[0],
				time:         record[1],
				isin:         "",
				ticker:       "",
				name:         "",
				shares:       0,
				price:        0,
				currency:     "",
				exchangeRate: 0,
				total:        total,
				chargeAmount: 0,
				notes:        record[12],
				id:           record[13],
				fee:          0,
			})
		}
	}

	return history, nil
}

func TransactionsToPortfolio(transactions []TradingTransaction, portfolioName string) Portfolio {
	portfolio := Portfolio{Name: portfolioName}

	for _, transaction := range transactions {
		productIndex := -1
		for i, product := range portfolio.Products {
			if product.Name == transaction.name {
				productIndex = i
				break
			}
		}

		if productIndex == -1 {
			newProduct := Product{Name: transaction.name, SymbolISIN: transaction.isin}
			if transaction.action == "Market buy" {
				newProduct.Quantity = transaction.shares
				newProduct.ValueEUR = transaction.total
			} else if transaction.action == "Market sell" {
				newProduct.Quantity = -transaction.shares
				newProduct.ValueEUR = -transaction.total
			}
			portfolio.Products = append(portfolio.Products, newProduct)
		} else {
			if transaction.action == "Market buy" {
				portfolio.Products[productIndex].Quantity += transaction.shares
				portfolio.Products[productIndex].ValueEUR += transaction.total
			} else if transaction.action == "Market sell" {
				portfolio.Products[productIndex].Quantity += transaction.shares
				portfolio.Products[productIndex].ValueEUR -= transaction.total
			}
			if portfolio.Products[productIndex].Quantity == 0 {
				portfolio.Products[productIndex].ValueEUR = 0
			}
		}
	}

	return portfolio
}

func MergePortfolios(destination Portfolio, source Portfolio) Portfolio {
	for _, sourceProduct := range source.Products {
		productIndex := -1
		for i, destinationProduct := range destination.Products {
			if destinationProduct.SymbolISIN == sourceProduct.SymbolISIN {
				productIndex = i
				break
			}
		}

		if productIndex == -1 {
			newProduct := Product{Name: sourceProduct.Name, SymbolISIN: sourceProduct.SymbolISIN, Quantity: sourceProduct.Quantity, ValueEUR: sourceProduct.ValueEUR}
			destination.Products = append(destination.Products, newProduct)
		} else {
			destination.Products[productIndex].Quantity += sourceProduct.Quantity
			destination.Products[productIndex].ValueEUR += sourceProduct.ValueEUR
		}
	}

	return destination
}
