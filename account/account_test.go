package account

import (
	"account/infrastructure"
	"account/models"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
	service *Service
)

func TestNewServiceHTTPDefaults(t *testing.T) {
	t.Parallel()
	s := NewService(nil)

	if s.http == nil {
		t.Fatal("default HTTP not instantiated")
	}

	if s.http.Client == nil {
		t.Fatal("default HTTP client not instantiated")
	}

	if s.http.BaseURL.String() != defaultBaseURL {
		t.Fatalf("NewService BaseURL is %v, should %v", s.http.BaseURL.String(), defaultBaseURL)
	}

	if s.http.Context == nil {
		t.Fatal("default execution context not instantiated")
	}

	if s.http.UserAgent != userAgent {
		t.Errorf("NewService UserAgent is %v, should be %v", s.http.UserAgent, userAgent)
	}
}

func TestListSuccess(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(OkListResponse))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()
	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)

	service = NewService(hp)
	accountsReceived, err := service.List(1, 2)
	if err != nil {
		t.Fatalf("List failed with error: %v", err)
	}

	resultCount := len(accountsReceived)

	if resultCount != 2 {
		t.Errorf("List returned %d, expected a count of 2", resultCount)
	}
}

func TestGetByIDInvalidArgument(t *testing.T) {
	t.Parallel()
	service = NewService(nil)
	_, err := service.GetByID("")
	expectedErrorMessage := "Invalid id argument"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Fatalf("Expected %s error message", expectedErrorMessage)
	}
}

func TestCreateNoIDInAccountToCreate(t *testing.T) {
	t.Parallel()
	attr := models.Attributes{
		AccountClassification: "Personal",
		AccountNumber:         "10000004",
		BankID:                "400302",
		BankIDCode:            "GBDSC",
		BaseCurrency:          "GBP",
		Bic:                   "NWBKGB42",
		Country:               "GB",
		CustomerID:            "234",
		IBAN:                  "GB28NWBK40030212764204",
	}

	accountToCreate := models.Account{
		Attributes:     attr,
		CreatedOn:      "2020-08-25T21:24:39.999Z",
		ModifiedOn:     "2020-08-25T21:24:39.999Z",
		OrganisationID: "bcd",
		Type:           "accounts",
		Version:        0,
	}
	service = NewService(nil)
	_, err := service.Create(accountToCreate)
	expectedErrorMessage := "ID field is missing, generate new UUID"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected %s error message", expectedErrorMessage)
	}
}

func TestDeleteByIDInvalidArgument(t *testing.T) {
	t.Parallel()
	service = NewService(nil)
	err := service.DeleteByID("")
	expectedErrorMessage := "Invalid id argument"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("Expected %s error message", expectedErrorMessage)
	}
}

func TestListInvalidPagingArguments(t *testing.T) {
	service = NewService(nil)
	t.Parallel() // marks TLog as capable of running in parallel with other tests
	tests := []struct {
		pageNumber int
		pageItems  int
	}{
		{0, 0},
		{0, 1},
		{1, 0},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("page number: %d, page items: %d", tt.pageNumber, tt.pageItems),
			func(t *testing.T) {
				t.Parallel()
				_, err := service.List(tt.pageNumber, tt.pageItems)
				expectedErrorMessage := "pageNumber and pageItem arguments must both be greater than 1"
				if !strings.Contains(err.Error(), expectedErrorMessage) {
					t.Errorf("Expected %s error message", expectedErrorMessage)
				}
			},
		)
	}
}

func TestListErrorResponse(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(BadListAccountResponse))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)
	service = NewService(hp)
	_, err := service.List(1, 2)
	if !strings.Contains(err.Error(), "downstream api error") {
		t.Error("GetByID not non existent account did not return expected error")
	}
}

func TestListSuccessPageItemsMoreThanAvalaibleAccounts(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(OkListResponse))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)
	service = NewService(hp)
	accountsReceived, err := service.List(1, 200)
	if err != nil {
		t.Fatalf("List failed with error: %v", err)
	}

	resultCount := len(accountsReceived)

	if resultCount <= 0 {
		t.Errorf("List returned %d, expected all the accounts", resultCount)
	}
}

func TestListSuccessPageOutOfBounds(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(OkListResponse))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)
	service = NewService(hp)
	accountsReceived, err := service.List(2, 3)
	if err != nil {
		t.Fatalf("List failed with error: %v", err)
	}

	resultCount := len(accountsReceived)

	if resultCount != 0 {
		t.Errorf("List returned %d, expected an empty list", resultCount)
	}
}

func TestListSuccessStartGreaterThanAccounts(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(OkListResponse))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)
	service = NewService(hp)
	accountsReceived, err := service.List(10, 1)
	if err != nil {
		t.Fatalf("List failed with error: %v", err)
	}

	resultCount := len(accountsReceived)

	if resultCount != 0 {
		t.Errorf("List returned %d, expected an empty list", resultCount)
	}
}

func TestGetByIDSuccess(t *testing.T) {
	t.Parallel()
	id := "abc"
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(OkGetResponse))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)
	service = NewService(hp)
	accountReceived, err := service.GetByID(id)

	if err != nil {
		t.Fatalf("GetByID failed with error: %v", err)
	}

	attr := models.Attributes{
		AccountClassification: "Personal",
		AccountNumber:         "10000004",
		BankID:                "400302",
		BankIDCode:            "GBDSC",
		BaseCurrency:          "GBP",
		Bic:                   "NWBKGB42",
		Country:               "GB",
		CustomerID:            "234",
		IBAN:                  "GB28NWBK40030212764204",
	}

	want := &models.Account{
		Attributes:     attr,
		CreatedOn:      "2020-08-25T21:24:39.999Z",
		ID:             "abc",
		ModifiedOn:     "2020-08-25T21:24:39.999Z",
		OrganisationID: "abc",
		Type:           "accounts",
		Version:        0,
	}

	if !reflect.DeepEqual(accountReceived, want) {
		t.Errorf("GetByID returned %+v instead of %+v", accountReceived, want)
	}
}

func TestGetByIDNonExistentAccount(t *testing.T) {
	t.Parallel()
	id := "abc"
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(NonExistentAccountResponse))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()
	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)

	service = NewService(hp)
	_, err := service.GetByID(id)

	if !strings.Contains(err.Error(), "downstream api error") {
		t.Error("GetByID not non existent account did not return expected error")
	}
}

func TestCreateSuccess(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(OkCreateResponse))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)
	service = NewService(hp)
	attr := models.Attributes{
		AccountClassification: "Personal",
		AccountNumber:         "10000004",
		BankID:                "400302",
		BankIDCode:            "GBDSC",
		BaseCurrency:          "GBP",
		Bic:                   "NWBKGB42",
		Country:               "GB",
		CustomerID:            "234",
		IBAN:                  "GB28NWBK40030212764204",
	}

	accountToCreate := models.Account{
		Attributes:     attr,
		CreatedOn:      "2020-08-25T21:24:39.999Z",
		ID:             "bcd",
		ModifiedOn:     "2020-08-25T21:24:39.999Z",
		OrganisationID: "bcd",
		Type:           "accounts",
		Version:        0,
	}
	accountReceived, err := service.Create(accountToCreate)

	if err != nil {
		t.Fatalf("GetByID failed with error: %v", err)
	}

	want := &models.Account{
		Attributes:     attr,
		CreatedOn:      "2020-08-25T21:24:39.999Z",
		ID:             "bcd",
		ModifiedOn:     "2020-08-25T21:24:39.999Z",
		OrganisationID: "bcd",
		Type:           "accounts",
		Version:        0,
	}

	if !reflect.DeepEqual(accountReceived, want) {
		t.Errorf("GetByID returned %+v instead of %+v", accountReceived, want)
	}
}

func TestCreateFailure(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(BadCreateResponse))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)
	service = NewService(hp)
	attr := models.Attributes{
		AccountClassification: "Personal",
		AccountNumber:         "10000004",
		BankID:                "400302",
		BankIDCode:            "GBDSC",
		BaseCurrency:          "GBP",
		Bic:                   "NWBKGB42",
		Country:               "GB",
		CustomerID:            "234",
		IBAN:                  "GB28NWBK40030212764204",
	}

	accountToCreate := models.Account{
		Attributes:     attr,
		CreatedOn:      "2020-08-25T21:24:39.999Z",
		ID:             "bcd",
		ModifiedOn:     "2020-08-25T21:24:39.999Z",
		OrganisationID: "bcd",
		Type:           "accounts",
		Version:        0,
	}
	_, err := service.Create(accountToCreate)
	if !strings.Contains(err.Error(), "downstream api error") {
		t.Error("create failure did not return expected error")
	}
}

func TestDeleteSuccess(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)
	service = NewService(hp)
	err := service.DeleteByID("bcd")
	if err != nil {
		t.Errorf("DeleteByID failed with error: %v", err)
	}
}

func TestDeleteFailure(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(BadDeleteResponse))
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	hp := infrastructure.NewHTTP(ctx, httpClient, nil, userAgent)
	service = NewService(hp)
	err := service.DeleteByID("bcd")
	if !strings.Contains(err.Error(), "downstream api error") {
		t.Errorf("DeleteByID failed with error: %v", err)
	}
}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	c := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return c, s.Close
}
