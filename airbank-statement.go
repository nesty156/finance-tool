package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dslipak/pdf"
)

type StatementOfAccount struct {
	AccountNumber string
	StartDate     time.Time
	EndDate       time.Time
	Transactions  []Transaction
}

type Transaction struct {
	AccountingDate     time.Time
	ExecutionDate      time.Time
	Type               string
	Code               string
	Name               string
	AccountOrDebitCard string
	Details            string
	AmountCZK          float64
	FeesCZK            float64
}

/* Parser airbank statement of account. */
func parseAirBankStatement(path string) (StatementOfAccount, error) {
	layout := "2. 1. 2006"
	account := StatementOfAccount{}

	r, err := pdf.Open(path)
	if err != nil {
		return account, err
	}

	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for i, word := range row.Content {
				if strings.TrimSpace(word.S) == "Číslo účtu:" {
					account.AccountNumber = row.Content[i+2].S
				}
				if strings.TrimSpace(word.S) == "Období výpisu:" {
					dateStr := row.Content[i+2].S
					dateParts := strings.Split(dateStr, " - ")

					startDate, err := time.Parse(layout, dateParts[0])
					if err != nil {
						fmt.Println("Error parsing start date:", err)
					}

					endDate, err := time.Parse(layout, dateParts[1])
					if err != nil {
						fmt.Println("Error parsing end date:", err)
					}
					account.StartDate = startDate
					account.EndDate = endDate
				}

				if strings.TrimSpace(word.S) == "Zaúčtování" {
					offsetStart := 20
					offsetEnd := 60
					for j := 0; j < 30; j++ {
						transactionRow := pdf.Row{Position: 0, Content: row.Content[i+offsetStart : i+offsetEnd]}
						transaction, offset, err := createTransaction(transactionRow)
						if err != nil {
							break
						}
						offsetStart += offset
						offsetEnd += offset
						account.Transactions = append(account.Transactions, transaction)
					}
				}
			}
		}
	}
	return account, nil
}

func createTransaction(row pdf.Row) (Transaction, int, error) {
	layout := "02.01.2006" // the layout string to parse the date, which follows the pattern "day.month.year"
	offset := 0
	AccountingDate, err := time.Parse(layout, strings.TrimSpace(row.Content[0].S))
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return Transaction{}, 0, err
	}
	ExecutionDate, err := time.Parse(layout, strings.TrimSpace(row.Content[2].S))
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return Transaction{}, 0, err
	}
	Type, Code := "", ""
	if row.Content[10].S == " " {
		Type = strings.TrimSpace(row.Content[6].S)
		Code = strings.TrimSpace(row.Content[8].S)
	} else {
		Type += strings.TrimSpace(row.Content[6].S)
		Type += " " + strings.TrimSpace(row.Content[8].S)
		Code = strings.TrimSpace(row.Content[10].S)
		offset += 2
	}
	Name, AccountOrDebitCard := "", ""
	if row.Content[10+offset].S == " " && row.Content[14+offset].S == " " {
		AccountOrDebitCard = strings.TrimSpace(row.Content[12+offset].S)
		offset -= 2
	} else {
		Name = strings.TrimSpace(row.Content[12+offset].S)
		AccountOrDebitCard = strings.TrimSpace(row.Content[14+offset].S)
		if row.Content[16+offset].S != " " {
			AccountOrDebitCard += "\n" + strings.TrimSpace(row.Content[16+offset].S)
			offset += 2
		}
	}
	Details := ""
	if row.Content[18+offset].S != " " {
		Details = strings.TrimSpace(row.Content[18+offset].S)
		offset += 2
		if row.Content[18+offset].S != " " {
			Details += "\n" + strings.TrimSpace(row.Content[18+offset].S)
			offset += 2
			if row.Content[18+offset].S != " " {
				Details += "\n" + strings.TrimSpace(row.Content[18+offset].S)
				offset += 2
				if row.Content[18+offset].S != " " {
					Details += "\n" + strings.TrimSpace(row.Content[18+offset].S)
					offset += 2
				}
			}
		}
	}

	Amount := strings.Replace(row.Content[20+offset].S, " ", "", -1)
	Amount = strings.Replace(Amount, ",", ".", -1)
	AmountCZK, err := strconv.ParseFloat(Amount, 64)
	if err != nil {
		fmt.Println("Error parsing float:", err)
		return Transaction{}, 0, err
	}
	Fees := strings.Replace(row.Content[24+offset].S, " ", "", -1)
	Fees = strings.Replace(Fees, ",", ".", -1)
	FeesCZK, err := strconv.ParseFloat(Fees, 64)
	if err != nil {
		fmt.Println("Error parsing float:", err)
		return Transaction{}, 0, err
	}

	result := Transaction{
		AccountingDate:     AccountingDate,
		ExecutionDate:      ExecutionDate,
		Type:               Type,
		Code:               Code,
		Name:               Name,
		AccountOrDebitCard: AccountOrDebitCard,
		Details:            Details,
		AmountCZK:          AmountCZK,
		FeesCZK:            FeesCZK,
	}

	/* Lenght of transation. */
	count := 0
	for i, content := range row.Content {
		if content.S == " " {
			count++
			if count == 6 {
				return result, i + 2, nil
			}
		}
	}

	return result, 4, nil
}
