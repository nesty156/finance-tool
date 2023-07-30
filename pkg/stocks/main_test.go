package stocks

import "testing"

func TestScrapePriceWithISIN(t *testing.T) {
	price, currency, err := ScrapePriceWithISIN("US88160R1014")
	if err != nil {
		t.Fatalf("Scraping failed: %s", err)
	}
	if price == 0 && currency == "" {
		t.Fatalf("Scraping failed: %s", err)
	}
}
