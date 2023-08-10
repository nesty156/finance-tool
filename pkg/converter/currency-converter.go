package converter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CoinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}

type ConvertRatesCZK struct {
	BTC float64
	USD float64
	EUR float64
}

type ConvertRatesEUR struct {
	BTC float64
	USD float64
	CZK float64
}

func GetBitcoinPrice(currency string) (float64, error) {
	if currency == "BTC" {
		return 1, nil
	}
	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/BTC-%s/spot", currency)
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var bodyBuilder strings.Builder
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		bodyBuilder.WriteString(line)
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	var cr CoinbaseResponse
	err = json.Unmarshal([]byte(bodyBuilder.String()), &cr)
	if err != nil {
		return 0, err
	}
	if price, err := strconv.ParseFloat(cr.Data.Amount, 64); err == nil {
		return price, nil
	}
	return 0, err
}

func GetConvertRate(fromCurrency, toCurrency string) float64 {
	from, _ := GetBitcoinPrice(fromCurrency)
	to, _ := GetBitcoinPrice(toCurrency)
	return to / from
}
