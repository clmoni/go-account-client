package integrationtests

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateUUID - generate ids for test
func GenerateUUID() string {
	uuid, _ := uuid.NewUUID()
	id := uuid.String()
	return id
}

// CreateListOfAccounts - mock ok list response
func CreateListOfAccounts() string {

	return fmt.Sprintf(`{
		"data": [
		{
			"attributes": {
			"account_classification": "Personal",
			"account_number": "10000004",
			"alternative_bank_account_names": null,
			"bank_id": "400302",
			"bank_id_code": "GBDSC",
			"base_currency": "GBP",
			"bic": "NWBKGB42",
			"country": "GB",
			"customer_id": "da968913-79ae-4a1c-8490-5597d28ecf5b",
			"iban": "GB28NWBK40030212764204"
			},
			"created_on": "2020-08-25T21:24:39.999Z",
			"id": "%s",
			"modified_on": "2020-08-25T21:24:39.999Z",
			"organisation_id": "538fd1a0-b62d-4b56-beb8-7836a1fedd2e",
			"type": "accounts",
			"version": 0
		},
		{
			"attributes": {
			"account_classification": "Personal",
			"account_number": "10000004",
			"alternative_bank_account_names": null,
			"bank_id": "400302",
			"bank_id_code": "GBDSC",
			"base_currency": "GBP",
			"bic": "NWBKGB42",
			"country": "GB",
			"customer_id": "d64797ec-c107-4ecf-a1cb-436e70fb85b8",
			"iban": "GB28NWBK40030212764204"
			},
			"created_on": "2020-08-25T21:24:48.358Z",
			"id": "%s",
			"modified_on": "2020-08-25T21:24:48.358Z",
			"organisation_id": "538fd1a0-b62d-4b56-beb8-7836a1fedd2e",
			"type": "accounts",
			"version": 0
		},
		{
			"attributes": {
			"account_classification": "Personal",
			"account_number": "10000004",
			"alternative_bank_account_names": null,
			"bank_id": "400302",
			"bank_id_code": "GBDSC",
			"base_currency": "GBP",
			"bic": "NWBKGB42",
			"country": "GB",
			"customer_id": "eca70c57-4f13-4907-a974-30459363509f",
			"iban": "GB28NWBK40030212764204"
			},
			"created_on": "2020-08-25T21:24:50.979Z",
			"id": "%s",
			"modified_on": "2020-08-25T21:24:50.979Z",
			"organisation_id": "538fd1a0-b62d-4b56-beb8-7836a1fedd2e",
			"type": "accounts",
			"version": 0
		}
		],
		"links": {
		"first": "",
		"last": "",
		"self": ""
		}
	}`, GenerateUUID(), GenerateUUID(), GenerateUUID())
}
