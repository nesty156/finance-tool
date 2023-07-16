package user

import (
	"fmt"
	"time"
)

var LoggedUser string
var DefaultCurrency string = "CZK"

type AppAccount struct {
	Name  string
	Stats []Stat
}

type Stat struct {
	Name       string
	Component  string
	InsertDate time.Time
	FilePath   string
	Value      float64
	Currency   string
}

type ConvertRatesCZK struct {
	BTC float64
	USD float64
	EUR float64
}

func Login() {
	fmt.Print("Enter your name: ")
	var name string
	fmt.Scanln(&name)
	LoggedUser = name
}

func Logout() {
	LoggedUser = ""
}

func (a *AppAccount) GetStatsInfo(ratesCZK ConvertRatesCZK) {
	total := 0.0
	for _, stat := range a.Stats {
		// insert date is month old
		if stat.InsertDate.AddDate(0, 1, 0).Before(time.Now()) {
			fmt.Printf("UPDATE THIS: Stat %s [%s] is month old\n", stat.Name, stat.Component)
		}

		// stats
		fmt.Printf("Stat %s [%s] is %.2f %s", stat.Name, stat.Component, stat.Value, stat.Currency)

		// convert value to CZK
		if stat.Currency != DefaultCurrency {
			var valueCZK float64
			switch stat.Currency {
			case "USD":
				valueCZK = stat.Value * ratesCZK.USD
			case "EUR":
				valueCZK = stat.Value * ratesCZK.EUR
			case "BTC":
				valueCZK = stat.Value * ratesCZK.BTC
			}
			fmt.Printf(" Converted %.2f %s\n", valueCZK, DefaultCurrency)
			total += valueCZK
		} else {
			fmt.Printf("\n")
			total += stat.Value
		}
	}
	// estimated networth
	fmt.Printf("Estimated networth is %.2f %s\n", total, DefaultCurrency)
	return
}

// Check if stat is already included in the user account and removes it.
func (a *AppAccount) RemoveStat(name string) bool {
	for i, stat := range a.Stats {
		if stat.Name == name {
			a.Stats = append(a.Stats[:i], a.Stats[i+1:]...)
			return true
		}
	}
	return false
}
