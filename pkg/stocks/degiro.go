package stocks

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type DegiroTransaction struct {
	Date         string  `csv:"Datum"`
	Time         string  `csv:"Čas"`
	ISIN         string  `csv:"ISIN"`
	Name         string  `csv:"Produkt"`
	Shares       float64 `csv:"Počet"`
	Price        float64 `csv:"Hodnota"`
	ExchangeRate float64 `csv:"Směnný kurz"`
	Total        float64 `csv:"Celkem"`
	ID           string  `csv:"ID objednávky"`
	fee          float64 `csv:"Transaction and/or third"`
}

func CreateDegiroPortfolio(fileName, portfolioName, currency string) (Portfolio, error) {
	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	txs := []*DegiroTransaction{}

	if err := gocsv.UnmarshalFile(csvFile, &txs); err != nil {
		panic(err)
	}

	return DegiroTXsToPortfolio(txs, portfolioName, currency), nil
}

func DegiroTXsToPortfolio(tradingTXs []*DegiroTransaction, portfolioName, currency string) Portfolio {
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
			newProduct := Product{Name: tx.Name, SymbolISIN: tx.ISIN, Quantity: tx.Shares, Value: -tx.Total, Currency: currency}
			portfolio.Products = append(portfolio.Products, newProduct)
		} else {
			portfolio.Products[productIndex].Quantity += tx.Shares
			portfolio.Products[productIndex].Value -= tx.Total
			if portfolio.Products[productIndex].Quantity == 0 {
				portfolio.Products[productIndex].Value = 0
			}
		}
		transaction, _ := DegiroTXConvert(*tx, "EUR")
		transactions = append(transactions, transaction)
	}

	portfolio.Transactions = transactions
	return portfolio
}

func DegiroTXConvert(tx DegiroTransaction, currency string) (Transaction, error) {
	layout := "02-01-2006 15:04"

	// Parse the date-time string into a time.Time value
	parsedTime, err := time.Parse(layout, tx.Date+" "+tx.Time)
	if err != nil {
		fmt.Println("Error parsing date-time:", err)
		return Transaction{}, err
	}
	return Transaction{
		Time:         parsedTime,
		ISIN:         tx.ISIN,
		Name:         tx.Name,
		Shares:       tx.Shares,
		Price:        tx.Price,
		Currency:     "EUR",
		ExchangeRate: tx.ExchangeRate,
		Total:        tx.Total,
		ID:           tx.ID,
		fee:          tx.fee,
	}, nil
}
