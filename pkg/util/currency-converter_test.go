package util

import "testing"

func TestGetBitcoinPrice(t *testing.T) {
	_, err := GetBitcoinPrice("CZK")
	if err != nil {
		t.Fatalf(`Cannot get the BTC price. %v`, err)
	}
}
