package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseMonetaStatement(t *testing.T) {
	xmlFile := filepath.Join("moneta", "fake-statement04-2019.xml")
	xmlData, err := ioutil.ReadFile(xmlFile)
	if err != nil {
		t.Fatalf("Failed to read file %s: %s", xmlFile, err)
	}
	_, err = parseMonetaStatement(xmlData)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}
