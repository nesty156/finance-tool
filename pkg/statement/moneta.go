package statement

import (
	"encoding/xml"
	"time"
)

/* Parser moneta statement of account. */
func ParseMonetaStatement(xmlData []byte) (StatementOfAccount, error) {
	var soa StatementOfAccount
	var result struct {
		Header struct {
			Account struct {
				Number    string  `xml:"number,attr"`
				Currency  string  `xml:"currency,attr"`
				Balance   float64 `xml:"stm-bgn-bal"`
				DebitTov  float64 `xml:"debit-tov"`
				CreditTov float64 `xml:"credit-tov"`
			} `xml:"account"`
			Stmt struct {
				Date       string `xml:"date,attr"`
				TrnCnt     int    `xml:"trn-cnt,attr"`
				Periodicty string `xml:"periodicity-description,attr"`
			} `xml:"stmt"`
		} `xml:"header"`
		Transactions []struct {
			ID        string   `xml:"id,attr"`
			AccountNo string   `xml:"other-account-number,attr"`
			DatePost  string   `xml:"date-post,attr"`
			DateEff   string   `xml:"date-eff,attr"`
			Amount    float64  `xml:"amount,attr"`
			Messages  []string `xml:"trn-messages>trn-message"`
		} `xml:"transactions>transaction"`
	}
	if err := xml.Unmarshal(xmlData, &result); err != nil {
		return soa, err
	}
	soa.AccountNumber = result.Header.Account.Number
	soa.Currnecy = result.Header.Account.Currency
	soa.EndDate, _ = time.Parse("2006-01-02", result.Header.Stmt.Date)
	soa.StartDate = soa.EndDate.AddDate(0, -1, 1) // Assume monthly statement
	soa.Transactions = make([]Transaction, len(result.Transactions))
	for i, t := range result.Transactions {
		accountingDate, _ := time.Parse("2006-01-02", t.DatePost)
		executionDate, _ := time.Parse("2006-01-02", t.DateEff)
		var details string
		for i := 1; i < len(t.Messages); i++ {
			if i == 1 {
				details += t.Messages[i]
			} else {
				details += "\n" + t.Messages[i]
			}
		}

		soa.Transactions[i] = Transaction{
			AccountingDate:     accountingDate,
			ExecutionDate:      executionDate,
			Type:               "",
			Code:               t.ID,
			Name:               t.Messages[0],
			AccountOrDebitCard: t.AccountNo,
			Details:            details,
			Amount:             t.Amount,
			Fee:                0.0,
		}

	}
	return soa, nil
}
