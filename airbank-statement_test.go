package main

import "testing"

/* Parser airbank statement of account. */
func TestParseAirBankStatement(t *testing.T) {
	want := 1
	ret, err := parseAirBankStatement()
	if want != ret || err != nil {
		t.Fatalf(`Parsing failed.`)
	}
}
