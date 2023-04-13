package main

import (
	"testing"
	"time"

	"github.com/dslipak/pdf"
)

/* Delete afterwards */
func TestParseAirBankStatement(t *testing.T) {
	want := 1
	_, err := parseAirBankStatement("eur-vypis01-2023.pdf")
	if want != 1 || err != nil {
		t.Fatalf(`Parsing failed.`)
	}
}

/* Parser airbank statement of account. */
func TestCreateTransaction(t *testing.T) {
	testCases := []struct {
		name          string
		row           pdf.Row
		expectedTrans Transaction
		expectedError error
	}{
		{
			name: "Name, Number, 2-row details",
			row: pdf.Row{Position: 0, Content: pdf.TextHorizontal{
				pdf.Text{S: "01.03.2023"},
				pdf.Text{S: ""},
				pdf.Text{S: "01.03.2023"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "Payment"},
				pdf.Text{S: ""},
				pdf.Text{S: "12345678"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "John Doe"},
				pdf.Text{S: ""},
				pdf.Text{S: "1234-5678-9012-3456"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "PFG1234"},
				pdf.Text{S: ""},
				pdf.Text{S: "Payment for goods"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "3 000,00"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "0,00"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
			}},
			expectedTrans: Transaction{
				AccountingDate:     time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				ExecutionDate:      time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				Type:               "Payment",
				Code:               "12345678",
				Name:               "John Doe",
				AccountOrDebitCard: "1234-5678-9012-3456",
				Details:            "PFG1234\nPayment for goods",
				AmountCZK:          3000.0,
				FeesCZK:            0.0,
			},
			expectedError: nil,
		},
		{
			name: "No Name, Number, 2-row details",
			row: pdf.Row{Position: 0, Content: pdf.TextHorizontal{
				pdf.Text{S: "01.03.2023"},
				pdf.Text{S: ""},
				pdf.Text{S: "01.03.2023"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "Payment"},
				pdf.Text{S: ""},
				pdf.Text{S: "12345678"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "1234-5678-9012-3456"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "PFG1234"},
				pdf.Text{S: ""},
				pdf.Text{S: "Payment for goods"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "3 000,00"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "0,00"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
			}},
			expectedTrans: Transaction{
				AccountingDate:     time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				ExecutionDate:      time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				Type:               "Payment",
				Code:               "12345678",
				Name:               "",
				AccountOrDebitCard: "1234-5678-9012-3456",
				Details:            "PFG1234\nPayment for goods",
				AmountCZK:          3000.0,
				FeesCZK:            0.0,
			},
			expectedError: nil,
		},
		{
			name: "No Name, Number, 3-row details",
			row: pdf.Row{Position: 0, Content: pdf.TextHorizontal{
				pdf.Text{S: "01.03.2023"},
				pdf.Text{S: ""},
				pdf.Text{S: "01.03.2023"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "Payment"},
				pdf.Text{S: ""},
				pdf.Text{S: "12345678"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "1234-5678-9012-3456"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "PFG1234"},
				pdf.Text{S: ""},
				pdf.Text{S: "Payment for goods"},
				pdf.Text{S: ""},
				pdf.Text{S: "from - chytry.house"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "3 000,00"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "0,00"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
			}},
			expectedTrans: Transaction{
				AccountingDate:     time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				ExecutionDate:      time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				Type:               "Payment",
				Code:               "12345678",
				Name:               "",
				AccountOrDebitCard: "1234-5678-9012-3456",
				Details:            "PFG1234\nPayment for goods\nfrom - chytry.house",
				AmountCZK:          3000.0,
				FeesCZK:            0.0,
			},
			expectedError: nil,
		},
		{
			name: "Name, Number, No details",
			row: pdf.Row{Position: 0, Content: pdf.TextHorizontal{
				pdf.Text{S: "01.03.2023"},
				pdf.Text{S: ""},
				pdf.Text{S: "01.03.2023"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "Payment"},
				pdf.Text{S: ""},
				pdf.Text{S: "12345678"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "John Doe"},
				pdf.Text{S: ""},
				pdf.Text{S: "1234-5678-9012-3456"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "3 000,00"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "0,00"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
			}},
			expectedTrans: Transaction{
				AccountingDate:     time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				ExecutionDate:      time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				Type:               "Payment",
				Code:               "12345678",
				Name:               "John Doe",
				AccountOrDebitCard: "1234-5678-9012-3456",
				Details:            "",
				AmountCZK:          3000.0,
				FeesCZK:            0.0,
			},
			expectedError: nil,
		},
		{
			name: "2-row type, Name, 2-row Number, 4-row details",
			row: pdf.Row{Position: 0, Content: pdf.TextHorizontal{
				pdf.Text{S: "01.03.2023"},
				pdf.Text{S: ""},
				pdf.Text{S: "01.03.2023"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "Fast"},
				pdf.Text{S: ""},
				pdf.Text{S: "Payment"},
				pdf.Text{S: ""},
				pdf.Text{S: "12345678"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "John Doe"},
				pdf.Text{S: ""},
				pdf.Text{S: "1234-5678-9012-3456 /"},
				pdf.Text{S: ""},
				pdf.Text{S: "1234-5678"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "PFG1234"},
				pdf.Text{S: ""},
				pdf.Text{S: "Payment"},
				pdf.Text{S: ""},
				pdf.Text{S: "for"},
				pdf.Text{S: ""},
				pdf.Text{S: "goods"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "3 000,00"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
				pdf.Text{S: "0,00"},
				pdf.Text{S: ""},
				pdf.Text{S: " "},
				pdf.Text{S: ""},
			}},
			expectedTrans: Transaction{
				AccountingDate:     time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				ExecutionDate:      time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				Type:               "Fast Payment",
				Code:               "12345678",
				Name:               "John Doe",
				AccountOrDebitCard: "1234-5678-9012-3456 /\n1234-5678",
				Details:            "PFG1234\nPayment\nfor\ngoods",
				AmountCZK:          3000.0,
				FeesCZK:            0.0,
			},
			expectedError: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trans, _, err := createTransaction(tc.row)
			if err != nil {
				t.Errorf("Test case %s failed: %v", tc.name, err)
			}
			if err != tc.expectedError {
				t.Errorf("Test case %s failed: Expected %v, but got %v", tc.name, tc.expectedError, err)
			}
			if trans != tc.expectedTrans {
				t.Errorf("Test case %s failed: Expected %v, but got %v", tc.name, tc.expectedTrans, trans)
			}
		})
	}
}
