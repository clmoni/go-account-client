package integrationtests

import (
	"account"
	"account/infrastructure"
	"account/models"
	"context"
	"encoding/json"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type (
	request struct {
		Data interface{} `json:"data,omitempty"`
	}
)

func getAccountAPIBaseURL() *url.URL {
	accountAPIBaseURL := os.Getenv("ACCOUNT_API_BASE_URL")
	if len(accountAPIBaseURL) == 0 {
		return nil
	}
	url, _ := url.Parse(accountAPIBaseURL)
	return url
}

func setup() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	baseURL := getAccountAPIBaseURL()
	h := infrastructure.NewHTTP(ctx, nil, baseURL, "")
	accountsCreated := []string{}
	accountService := account.NewService(h)

	accounts := &[]models.Account{}
	res := &request{Data: accounts}
	json.Unmarshal([]byte(CreateListOfAccounts()), res)
	for _, account := range *accounts {
		a, err := accountService.Create(account)
		if err != nil {
			return nil, err
		}
		accountsCreated = append(accountsCreated, a.ID)
	}

	return accountsCreated, nil
}

func teardown(accountsCreated []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	baseURL := getAccountAPIBaseURL()
	h := infrastructure.NewHTTP(ctx, nil, baseURL, "")
	accountService := account.NewService(h)

	for _, id := range accountsCreated {
		accountService.DeleteByID(id)
	}
}

func TestListOk(t *testing.T) {
	t.Parallel()
	accountsCreated, err := setup()
	if err != nil {
		t.Fatalf("Unable to set up, error occurred: %v", err)
	}
	defer teardown(accountsCreated)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	baseURL := getAccountAPIBaseURL()
	h := infrastructure.NewHTTP(ctx, nil, baseURL, "")
	accountService := account.NewService(h)

	accounts, err := accountService.List(1, 2)

	if err != nil {
		t.Fatalf("List failed with error: %v", err)
	}
	resultCount := len(accounts)
	if resultCount != 2 {
		t.Errorf("List returned %d, expected 2 in account list", resultCount)
	}
}

func TestGetOk(t *testing.T) {
	t.Parallel()
	accountsCreated, err := setup()
	if err != nil {
		t.Fatalf("Unable to set up, error occurred: %v", err)
	}
	defer teardown(accountsCreated)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	baseURL := getAccountAPIBaseURL()
	h := infrastructure.NewHTTP(ctx, nil, baseURL, "")
	accountService := account.NewService(h)

	account, err := accountService.GetByID(accountsCreated[0])

	if err != nil {
		t.Fatalf("GetByID failed with error: %v", err)
	}
	if account.ID != accountsCreated[0] {
		t.Errorf("GetByID returned %s, expected %s in account list", account.ID, accountsCreated[0])
	}
}

func TestGetNotFound(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	baseURL := getAccountAPIBaseURL()
	h := infrastructure.NewHTTP(ctx, nil, baseURL, "")
	accountService := account.NewService(h)

	_, err := accountService.GetByID("538fd1a0-b62d-4b56-beb8-7836a1fedd2e")

	if !strings.Contains(err.Error(), "downstream api error") {
		t.Error("GetByID not non existent account did not return expected error")
	}
}

func TestCreateSuccess(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := GenerateUUID()
	attr := models.Attributes{
		AccountClassification: "Personal",
		AccountNumber:         "10000004",
		BankID:                "400302",
		BankIDCode:            "GBDSC",
		BaseCurrency:          "GBP",
		Bic:                   "NWBKGB42",
		Country:               "GB",
		CustomerID:            "da968913-79ae-4a1c-8490-5597d28ecf5b",
		IBAN:                  "GB28NWBK40030212764204",
	}

	accountToCreate := models.Account{
		Attributes:     attr,
		CreatedOn:      "2020-08-25T21:24:39.999Z",
		ID:             id,
		ModifiedOn:     "2020-08-25T21:24:39.999Z",
		OrganisationID: "da968913-79ae-4a1c-8490-5597d28ecf5b",
		Type:           "accounts",
		Version:        0,
	}

	baseURL := getAccountAPIBaseURL()
	h := infrastructure.NewHTTP(ctx, nil, baseURL, "")
	accountService := account.NewService(h)
	accountReceived, err := accountService.Create(accountToCreate)

	if err != nil {
		t.Fatalf("GetByID failed with error: %v", err)
	}

	want := &models.Account{
		Attributes:     attr,
		CreatedOn:      "2020-08-25T21:24:39.999Z",
		ID:             id,
		ModifiedOn:     "2020-08-25T21:24:39.999Z",
		OrganisationID: "da968913-79ae-4a1c-8490-5597d28ecf5b",
		Type:           "accounts",
		Version:        0,
	}

	if !reflect.DeepEqual(accountReceived.Attributes, want.Attributes) {
		t.Errorf("GetByID returned %+v instead of %+v", accountReceived, want)
	}
}

func TestCreateFailure(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	attr := models.Attributes{
		AccountClassification: "Personal",
		AccountNumber:         "10000004",
		BankID:                "400302",
		BankIDCode:            "GBDSC",
		BaseCurrency:          "GBP",
		Bic:                   "NWBKGB42",
		Country:               "GB",
		CustomerID:            "da968913-79ae-4a1c-8490-5597d28ecf5b",
		IBAN:                  "GB28NWBK40030212764204",
	}

	accountToCreate := models.Account{
		Attributes:     attr,
		CreatedOn:      "2020-08-25T21:24:39.999Z",
		ID:             "da968913",
		ModifiedOn:     "2020-08-25T21:24:39.999Z",
		OrganisationID: "da968913-79ae-4a1c-8490-5597d28ecf5b",
		Type:           "accounts",
		Version:        0,
	}
	baseURL := getAccountAPIBaseURL()
	h := infrastructure.NewHTTP(ctx, nil, baseURL, "")
	accountService := account.NewService(h)
	_, err := accountService.Create(accountToCreate)

	if !strings.Contains(err.Error(), "downstream api error") {
		t.Error("GetByID not non existent account did not return expected error")
	}
}

func TestDeleteSuccess(t *testing.T) {
	t.Parallel()
	accountsCreated, err := setup()
	if err != nil {
		t.Fatalf("Unable to set up, error occurred: %v", err)
	}
	defer teardown(accountsCreated)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	baseURL := getAccountAPIBaseURL()
	h := infrastructure.NewHTTP(ctx, nil, baseURL, "")
	accountService := account.NewService(h)

	err = accountService.DeleteByID(accountsCreated[0])
	if err != nil {
		t.Errorf("DeleteByID failed with error: %v", err)
	}
}
