package account

const (
	// OkGetResponse - mock ok get response
	OkGetResponse = `{
		"data": {
		  "attributes": {
			"account_classification": "Personal",
			"account_number": "10000004",
			"alternative_bank_account_names": null,
			"bank_id": "400302",
			"bank_id_code": "GBDSC",
			"base_currency": "GBP",
			"bic": "NWBKGB42",
			"country": "GB",
			"customer_id": "234",
			"iban": "GB28NWBK40030212764204"
		  },
		  "created_on": "2020-08-25T21:24:39.999Z",
		  "id": "abc",
		  "modified_on": "2020-08-25T21:24:39.999Z",
		  "organisation_id": "abc",
		  "type": "accounts",
		  "version": 0
		},
		"links": {
		  "self": "/v1/organisation/accounts/abc"
		}
	  }`

	// OkCreateResponse - mock ok create response
	OkCreateResponse = `{
		"data": {
		  "attributes": {
			"account_classification": "Personal",
			"account_number": "10000004",
			"alternative_bank_account_names": null,
			"bank_id": "400302",
			"bank_id_code": "GBDSC",
			"base_currency": "GBP",
			"bic": "NWBKGB42",
			"country": "GB",
			"customer_id": "234",
			"iban": "GB28NWBK40030212764204"
		  },
		  "created_on": "2020-08-25T21:24:39.999Z",
		  "id": "bcd",
		  "modified_on": "2020-08-25T21:24:39.999Z",
		  "organisation_id": "bcd",
		  "type": "accounts",
		  "version": 0
		},
		"links": {
		  "self": "/v1/organisation/accounts/bcd"
		}
	  }`

	// BadDeleteResponse - mock bad delete request response error
	BadDeleteResponse = `{"error_message":"id is not a valid uuid"}`

	// BadCreateResponse - mock bad create request response error
	BadCreateResponse = `{"error_message":"Account cannot be created as it violates a duplicate constraint"}`

	// NonExistentAccountResponse - mock non existent account get response error
	NonExistentAccountResponse = `{"error_message":"record 626e880a-e719-11ea-8eaa-8c85903c0c70 does not exist"}`

	// OkListResponse - mock ok list response
	OkListResponse = `{
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
			  "customer_id": "234",
			  "iban": "GB28NWBK40030212764204"
			},
			"created_on": "2020-08-25T21:24:39.999Z",
			"id": "626e880a-e719-11ea-8eaa-8c85903c0c20",
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
			  "customer_id": "234",
			  "iban": "GB28NWBK40030212764204"
			},
			"created_on": "2020-08-25T21:24:48.358Z",
			"id": "67687384-e719-11ea-8ee4-8c85903c0c20",
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
			  "customer_id": "234",
			  "iban": "GB28NWBK40030212764204"
			},
			"created_on": "2020-08-25T21:24:50.979Z",
			"id": "68f76b60-e719-11ea-90b6-8c85903c0c20",
			"modified_on": "2020-08-25T21:24:50.979Z",
			"organisation_id": "538fd1a0-b62d-4b56-beb8-7836a1fedd2e",
			"type": "accounts",
			"version": 0
		  }
		],
		"links": {
		  "first": "/v1/organisation/accounts?page%5Bnumber%5D=first",
		  "last": "/v1/organisation/accounts?page%5Bnumber%5D=last",
		  "self": "/v1/organisation/accounts?page%5Bnumber%5D=last"
		}
	  }`

	// BadListAccountResponse - mock listing error (if listing was to return some error response)
	BadListAccountResponse = `{"error_message":"some downstream error"}`
)
