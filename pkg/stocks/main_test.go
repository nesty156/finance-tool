package stocks

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestScrapePriceWithISIN(t *testing.T) {
	isins := []string{"US88160R1014", "US0378331005", "IE00BK5BQT80", "IE00B5BMR087", "DE0007664039"}

	// Measure the execution time of the function
	startTime := time.Now()

	var wg sync.WaitGroup

	for _, isin := range isins {
		wg.Add(1)
		go func(isin string) {
			defer wg.Done()

			_, _, err := ScrapePriceWithISIN(isin)
			if err != nil {
				t.Errorf("Error during scraping for ISIN %s: %s", isin, err)
			}
		}(isin)
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)

	fmt.Printf("Scraping took: %s\n", elapsedTime)
}
