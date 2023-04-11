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
func parseAirBankStatement() (int, error) {
	content, err := readPdf("vypis03-2023.pdf") // Read local pdf file
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
	return 1, nil
}

func readPdf(path string) (StatementOfAccount, error) {
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
			println(">>>> row: ", row.Position)
			for i, word := range row.Content {
				if strings.TrimSpace(word.S) == "Číslo účtu:" {
					account.AccountNumber = row.Content[i+2].S
				}

				if strings.TrimSpace(word.S) == "Zaúčtování" {
					offsetStart := 20
					offsetEnd := 70
					for j := 0; j < 20; j++ {
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
	//TODO: Transaction every column is divided by row.Content[i] == " " after (6th space we have end of transaction)
	layout := "02.01.2006" // the layout string to parse the date, which follows the pattern "day.month.year"
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
	Type := strings.TrimSpace(row.Content[6].S)
	Code := strings.TrimSpace(row.Content[8].S)
	Name := strings.TrimSpace(row.Content[12].S)
	AccountOrDebitCard := strings.TrimSpace(row.Content[14].S)
	Details := strings.TrimSpace(row.Content[18].S)
	Details += strings.TrimSpace(row.Content[20].S)
	AmountCZK, FeesCZK := 0.0, 0.0
	if AmountCZK, err := strconv.ParseFloat(strings.TrimSpace(row.Content[24].S), 64); err == nil {
		fmt.Println(AmountCZK)
	}
	if FeesCZK, err := strconv.ParseFloat(strings.TrimSpace(row.Content[28].S), 64); err == nil {
		fmt.Println(FeesCZK)
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
		_, err := time.Parse(layout, strings.TrimSpace(content.S))
		if err == nil {
			count++
			if count == 4 {
				return result, i - 2, nil
			}
		}
	}

	return result, 2, nil
}
