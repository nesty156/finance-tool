package stocks

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nesty156/finance-tool/pkg/converter"
)

type Portfolio struct {
	Name         string
	Currency     string
	Products     []Product
	Transactions []Transaction
}

type Product struct {
	Name       string
	SymbolISIN string
	Quantity   float64
	Value      float64
	Currency   string
}

type Transaction struct {
	Action       string
	Time         time.Time
	ISIN         string
	Ticker       string
	Name         string
	Shares       float64
	Price        float64
	Currency     string
	ExchangeRate float64
	Total        float64
	ChargeAmount float64
	Notes        string
	ID           string
	fee          float64
}

func ScrapePriceWithISIN(isin string) (float64, string, error) {
	// Get the URL for the ISIN
	url := fmt.Sprintf("https://www.marketscreener.com/search/?q=%s", isin)
	// Make an HTTP request to fetch the page content
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching the page: %s", err)
	}
	defer response.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatalf("Error parsing the HTML: %s", err)
	}

	// Find the price element (you need to inspect the website's HTML structure to find the correct selector)
	target := doc.Find("td[aria-label='Price'] span.last")
	priceStr := target.Nodes[0].FirstChild.Data
	currency := strings.TrimSpace(target.Nodes[0].NextSibling.FirstChild.FirstChild.Data)

	// Convert the price string to a float
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Fatalf("Error converting price to float: %s", err)
	}

	fmt.Printf("Price: %.2f %s\n", price, currency)

	return price, currency, nil
}

func PortfolioValue(portfolio Portfolio) float64 {
	total := 0.0

	// Loop through products and add up the amounts
	for _, product := range portfolio.Products {
		shareValue, currency, _ := ScrapePriceWithISIN(product.SymbolISIN)
		if currency != portfolio.Currency && portfolio.Currency == "EUR" {
			convertRate := converter.GetConvertRatesEUR()
			if currency == "USD" {
				shareValue = shareValue * convertRate.USD
			} else if currency == "CZK" {
				shareValue = shareValue * convertRate.CZK
			}
		} else if currency != portfolio.Currency && portfolio.Currency == "CZK" {
			convertRate := converter.GetConvertRatesCZK()
			if currency == "USD" {
				shareValue = shareValue * convertRate.USD
			} else if currency == "EUR" {
				shareValue = shareValue * convertRate.EUR
			}
		}
		product.Value = shareValue * product.Quantity
		total += product.Value
	}

	return total
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
			newProduct := Product{Name: sourceProduct.Name, SymbolISIN: sourceProduct.SymbolISIN, Quantity: sourceProduct.Quantity, Value: sourceProduct.Value}
			destination.Products = append(destination.Products, newProduct)
		} else {
			destination.Products[productIndex].Quantity += sourceProduct.Quantity
			destination.Products[productIndex].Value += sourceProduct.Value
		}
	}

	return destination
}
