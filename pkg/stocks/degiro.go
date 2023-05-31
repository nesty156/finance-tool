package stocks

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseDegiroPortfolio(csvData []byte, portfolioName string) (Portfolio, error) {
	reader := csv.NewReader(strings.NewReader(string(csvData)))
	records, err := reader.ReadAll()
	if err != nil {
		return Portfolio{}, fmt.Errorf("failed to read CSV: %v", err)
	}

	portfolio := Portfolio{Name: portfolioName}
	for _, record := range records[1:] {
		product := Product{}
		product.Name = record[0]
		product.SymbolISIN = record[1]
		quantity, err := strconv.Atoi(record[2])
		if err != nil {
			quantity = 1
		}
		product.Quantity = float64(quantity)
		valueEUR, err := strconv.ParseFloat(strings.ReplaceAll(record[5], ",", "."), 64)
		if err != nil {
			fmt.Println("Failed to parse ValueEUR:", err)
			os.Exit(1)
		}
		product.ValueEUR = valueEUR
		fmt.Println(product)

		portfolio.Products = append(portfolio.Products, product)
	}

	return portfolio, nil
}

func PortfolioValue(portfolio Portfolio) float64 {
	total := 0.0

	// Loop through products and add up the amounts
	for _, product := range portfolio.Products {
		total += product.ValueEUR
	}

	return total
}
