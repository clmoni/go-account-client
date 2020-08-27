package models

// Attributes - account atrributes
type Attributes struct {
	AccountClassification       string `json:"account_classification,omitempty"`
	AccountNumber               string `json:"account_number,omitempty"`
	AlternativeBankAccountNames string `json:"alternative_bank_account_names,omitempty"`
	BankID                      string `json:"bank_id,omitempty"`
	BankIDCode                  string `json:"bank_id_code,omitempty"`
	BaseCurrency                string `json:"base_currency,omitempty"`
	Bic                         string `json:"bic"`
	Country                     string `json:"country,omitempty"`
	CustomerID                  string `json:"customer_id,omitempty"`
	IBAN                        string `json:"iban,omitempty"`
}
