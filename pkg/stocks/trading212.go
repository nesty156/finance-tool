package stocks

import (
	"os"

	"github.com/gocarina/gocsv"
)

type TradingTransaction struct {
	Action       string  `csv:"Action"`
	Time         string  `csv:"Time"`
	ISIN         string  `csv:"ISIN"`
	Ticker       string  `csv:"Ticker"`
	Name         string  `csv:"Name"`
	Shares       float64 `csv:"No. of shares"`
	Price        float64 `csv:"Price / share"`
	Currency     string  `csv:"Currency (Price / share)"`
	ExchangeRate float64 `csv:"Exchange rate"`
	Total        float64 `csv:"Total (EUR)"`
	ChargeAmount float64 `csv:"Charge amount (EUR)"`
	Notes        string  `csv:"Notes"`
	ID           string  `csv:"ID"`
	fee          float64 `csv:"Currency conversion fee (EUR)"`
}

func ParseTrading212History(fileName string) ([]*TradingTransaction, error) {
	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	history := []*TradingTransaction{}

	if err := gocsv.UnmarshalFile(csvFile, &history); err != nil { // Load clients from file
		panic(err)
	}

	return history, nil
}

func TransactionsToPortfolio(transactions []*TradingTransaction, portfolioName string) Portfolio {
	portfolio := Portfolio{Name: portfolioName}

	for _, transaction := range transactions {
		productIndex := -1
		for i, product := range portfolio.Products {
			if product.Name == transaction.Name {
				productIndex = i
				break
			}
		}

		if productIndex == -1 {
			newProduct := Product{Name: transaction.Name, SymbolISIN: transaction.ISIN}
			if transaction.Action == "Market buy" {
				newProduct.Quantity = transaction.Shares
				newProduct.ValueEUR = transaction.Total
			} else if transaction.Action == "Market sell" {
				newProduct.Quantity = -transaction.Shares
				newProduct.ValueEUR = -transaction.Total
			}
			portfolio.Products = append(portfolio.Products, newProduct)
		} else {
			if transaction.Action == "Market buy" {
				portfolio.Products[productIndex].Quantity += transaction.Shares
				portfolio.Products[productIndex].ValueEUR += transaction.Total
			} else if transaction.Action == "Market sell" {
				portfolio.Products[productIndex].Quantity += transaction.Shares
				portfolio.Products[productIndex].ValueEUR -= transaction.Total
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
