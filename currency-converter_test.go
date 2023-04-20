package main

import "testing"

func TestGetBitcoinPrice(t *testing.T) {
	_, err := getBitcoinPrice("CZK")
	if err != nil {
		t.Fatalf(`Cannot get the BTC price. %v`, err)
	}
}
