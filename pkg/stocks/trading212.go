package stocks

import (
	"fmt"
	"os"
	"time"

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

func CreateTrading212Portfolio(fileName, portfolioName, currency string) (Portfolio, error) {
	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	txs := []*TradingTransaction{}

	if err := gocsv.UnmarshalFile(csvFile, &txs); err != nil {
		panic(err)
	}

	return TradingTXsToPortfolio(txs, portfolioName, currency), nil
}

func TradingTXsToPortfolio(tradingTXs []*TradingTransaction, portfolioName, currency string) Portfolio {
	portfolio := Portfolio{Name: portfolioName, Currency: currency}
	transactions := []Transaction{}

	for _, tx := range tradingTXs {
		productIndex := -1
		for i, product := range portfolio.Products {
			if product.Name == tx.Name {
				productIndex = i
				break
			}
		}

		if productIndex == -1 {
			newProduct := Product{Name: tx.Name, SymbolISIN: tx.ISIN, Currency: currency}
			if tx.Action == "Market buy" {
				newProduct.Quantity = tx.Shares
				newProduct.Value = tx.Total
			} else if tx.Action == "Market sell" {
				newProduct.Quantity = -tx.Shares
				newProduct.Value = -tx.Total
			}
			portfolio.Products = append(portfolio.Products, newProduct)
		} else {
			if tx.Action == "Market buy" {
				portfolio.Products[productIndex].Quantity += tx.Shares
				portfolio.Products[productIndex].Value += tx.Total
			} else if tx.Action == "Market sell" {
				portfolio.Products[productIndex].Quantity += tx.Shares
				portfolio.Products[productIndex].Value -= tx.Total
			}
			if portfolio.Products[productIndex].Quantity == 0 {
				portfolio.Products[productIndex].Value = 0
			}
		}
		transaction, _ := TradingTXConvert(*tx)
		transactions = append(transactions, transaction)
	}

	portfolio.Transactions = transactions
	return portfolio
}

func TradingTXConvert(tx TradingTransaction) (Transaction, error) {
	layout := "2006-01-02 15:04:05"

	// Parse the date-time string into a time.Time value
	parsedTime, err := time.Parse(layout, tx.Time)
	if err != nil {
		fmt.Println("Error parsing date-time:", err)
		return Transaction{}, err
	}

	return Transaction{
		Action:       tx.Action,
		Time:         parsedTime,
		ISIN:         tx.ISIN,
		Ticker:       tx.Ticker,
		Name:         tx.Name,
		Shares:       tx.Shares,
		Price:        tx.Price,
		Currency:     tx.Currency,
		ExchangeRate: tx.ExchangeRate,
		Total:        tx.Total,
		ChargeAmount: tx.ChargeAmount,
		Notes:        tx.Notes,
		ID:           tx.ID,
		fee:          tx.fee,
	}, nil
}
