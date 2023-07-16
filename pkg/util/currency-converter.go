package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

func GetBitcoinPrice(currency string) (float64, error) {
	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/BTC-%s/spot", currency)
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var cr CoinbaseResponse
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return 0, err
	}
	if price, err := strconv.ParseFloat(cr.Data.Amount, 64); err == nil {
		return price, nil
	}
	return 0, err
}

func GetConvertRatesCZK() ConvertRatesCZK {
	btcCZK, _ := GetBitcoinPrice("CZK")
	btcEUR, _ := GetBitcoinPrice("EUR")
	btcUSD, _ := GetBitcoinPrice("USD")
	return ConvertRatesCZK{BTC: btcCZK, EUR: btcCZK / btcEUR, USD: btcCZK / btcUSD}
}
