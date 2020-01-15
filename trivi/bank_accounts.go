package trivi

import (
	"net/http"
)

func (t *triviApi) ListBankAccounts() (err error, bankAccountsResponse BankAccountsResponse) {
	err = t.execute(http.MethodGet, TRIVI_BANK_ACCOUNTS, nil, &bankAccountsResponse)
	return
}

type BankAccountsResponse struct {
	BankAccounts []BankAccounts `json:"bankAccounts"`
	PagesCount   int            `json:"pagesCount"`
	CurrentPage  int            `json:"currentPage"`
	PageSize     int            `json:"pageSize"`
}

type BankAccounts struct {
	ID                    int         `json:"id"`
	Name                  string      `json:"name"`
	Currency              string      `json:"currency"`
	AccountNo             string      `json:"accountNo"`
	AccountCode           string      `json:"accountCode"`
	AccountSwift          string      `json:"accountSwift"`
	AccountIban           string      `json:"accountIban"`
	StatementFormatID     int         `json:"statementFormatId"`
	PaymentOrderFormatID  int         `json:"paymentOrderFormatId"`
	BankID                int         `json:"bankId"`
	InitialAmount         float64     `json:"initialAmount"`
	IsMain                interface{} `json:"isMain"`
	ProcessingState       int         `json:"processingState"`
	ProcessingMessage     string      `json:"processingMessage"`
	ProcessedFinAccount   string      `json:"processedFinAccount"`
	BankName              string      `json:"bankName"`
	UnmatchedPaymentCount int         `json:"unmatchedPaymentCount"`
}
