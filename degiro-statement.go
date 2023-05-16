package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Portfolio struct {
	Products []Product
}

type Product struct {
	Name       string
	SymbolISIN string
	Quantity   int
	ValueEUR   float64
}

func parseDegiroPortfolio(csvData []byte) (Portfolio, error) {
	reader := csv.NewReader(strings.NewReader(string(csvData)))
	records, err := reader.ReadAll()
	if err != nil {
		return Portfolio{}, fmt.Errorf("failed to read CSV: %v", err)
	}

	portfolio := Portfolio{}
	for _, record := range records[1:] {
		product := Product{}
		product.Name = record[0]
		product.SymbolISIN = record[1]
		quantity, err := strconv.Atoi(record[2])
		if err != nil {
			quantity = 1
		}
		product.Quantity = quantity
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
