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
				wordStr := strings.TrimSpace(word.S)
				switch wordStr {
				case "Číslo účtu:":
					account.AccountNumber = row.Content[i+2].S

				case "Období výpisu:":
					dateStr := row.Content[i+2].S
					dateParts := strings.Split(dateStr, " - ")
					if startDate, err := time.Parse(layout, dateParts[0]); err != nil {
						fmt.Println("Error parsing start date:", err)
					} else {
						account.StartDate = startDate
					}

					if endDate, err := time.Parse(layout, dateParts[1]); err != nil {
						fmt.Println("Error parsing end date:", err)
					} else {
						account.EndDate = endDate
					}

				case "Zaúčtování":
					start, end := 20, 60
					for j := 0; j < 30; j++ {
						transactionRow := pdf.Row{Position: 0, Content: row.Content[i+start : i+end]}
						if transaction, offset, err := createTransaction(transactionRow); err != nil {
							break
						} else {
							account.Transactions = append(account.Transactions, transaction)
							start += offset
							end += offset
						}
					}
				}
			}
		}
	}
	return account, nil
}

/* From row of statement and creates transaction with all data */
func createTransaction(row pdf.Row) (Transaction, int, error) {
	offset := 0
	const numFields = 6
	layout := "02.01.2006"

	parseDate := func(date string) (time.Time, error) {
		return time.Parse(layout, strings.TrimSpace(date))
	}

	parseFloat := func(value string) (float64, error) {
		value = strings.ReplaceAll(value, " ", "")
		value = strings.ReplaceAll(value, ",", ".")
		return strconv.ParseFloat(value, 64)
	}

	/* Automatically adds offset */
	parseField := func(index int) string {
		return strings.TrimSpace(row.Content[index+offset].S)
	}

	var err error
	temp := Transaction{}
	temp.AccountingDate, err = parseDate(parseField(0))
	if err != nil {
		return Transaction{}, 0, err
	}

	temp.ExecutionDate, err = parseDate(parseField(2))
	if err != nil {
		return Transaction{}, 0, err
	}

	if row.Content[10].S == " " {
		temp.Type = parseField(6)
		temp.Code = parseField(8)
	} else {
		temp.Type = parseField(6) + " " + parseField(8)
		temp.Code = parseField(10)
		offset += 2
	}

	if row.Content[10+offset].S == " " && row.Content[14+offset].S == " " {
		temp.AccountOrDebitCard = parseField(12)
		offset -= 2
	} else {
		temp.Name = parseField(12)
		temp.AccountOrDebitCard = parseField(14)
		if row.Content[16+offset].S != " " {
			temp.AccountOrDebitCard += "\n" + parseField(16)
			offset += 2
		}
	}

	if row.Content[18+offset].S != " " {
		temp.Details = parseField(18)
		offset += 2
		for row.Content[18+offset].S != " " {
			temp.Details += "\n" + parseField(18)
			offset += 2
		}
	}

	temp.AmountCZK, err = parseFloat(parseField(20))
	if err != nil {
		return Transaction{}, 0, err
	}

	temp.FeesCZK, err = parseFloat(parseField(24))
	if err != nil {
		return Transaction{}, 0, err
	}

	count := 0
	for i := range row.Content {
		if row.Content[i].S == " " {
			count++
			if count == numFields {
				return temp, i + 2, nil
			}
		}
	}

	return temp, 4, nil
}
