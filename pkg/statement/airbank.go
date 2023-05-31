package statement

import (
	"fmt"
	"strings"
	"time"

	"github.com/dslipak/pdf"
)

type StatementOfAccount struct {
	AccountNumber string
	StartDate     time.Time
	EndDate       time.Time
	Currnecy      string
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
	Amount             float64
	Fee                float64
}

/* Parser airbank statement of account. */
func ParseAirBankStatement(path string) (StatementOfAccount, error) {
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
					account.AccountNumber = strings.ReplaceAll(row.Content[i+2].S, " ", "")

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

				case "Měna:":
					account.Currnecy = row.Content[i+2].S

				case "Zaúčtování":
					start, end := 20, 60
					for j := 0; j < 30; j++ {
						if i+end > len(row.Content) {
							end = len(row.Content) - i
						}
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
