package account

import (
	"account/infrastructure"
	"account/models"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	version        = "0.1"
	apiVersion     = "v1"
	userAgent      = "go-account/" + version
	defaultBaseURL = "http://localhost:8080"
)

// Service - handles communication with the account endpoint
type Service struct {
	http *infrastructure.HTTP
}

// NewService - initialise the service along with the client it will use for
// making request to the account api
func NewService(h *infrastructure.HTTP) *Service {
	return &Service{
		http: validateInjectedClientOrDefault(h),
	}
}

// GetByID - get new account by ID
func (s *Service) GetByID(id string) (*models.Account, error) {
	if len(id) <= 0 {
		return nil, errors.New("Invalid id argument")
	}
	account := &models.Account{}
	getAccountPath := fmt.Sprintf("%s/organisation/accounts/%s", apiVersion, id)
	err := s.http.Get(getAccountPath, account)
	return account, err
}

// Create - create a new account
func (s *Service) Create(a models.Account) (*models.Account, error) {
	account := &models.Account{}
	if len(a.ID) <= 0 {
		return nil, errors.New("ID field is missing, generate new UUID")
	}
	createAccountPath := fmt.Sprintf("%s/organisation/accounts", apiVersion)
	err := s.http.Post(createAccountPath, a, account)
	return account, err
}

// DeleteByID - get new account by ID
func (s *Service) DeleteByID(id string) error {
	if len(id) <= 0 {
		return errors.New("Invalid id argument")
	}
	getAccountPath := fmt.Sprintf("%s/organisation/accounts/%s?version=0", apiVersion, id)
	err := s.http.Delete(getAccountPath)
	return err
}

// List - list paged accounts
func (s *Service) List(pageNumber, pageItems int) ([]models.Account, error) {
	if pageNumber <= 0 || pageItems <= 0 {
		return nil, errors.New("pageNumber and pageItem arguments must both be greater than 1")
	}
	start := (pageNumber - 1) * pageItems
	stop := start + pageItems

	accounts := &[]models.Account{}
	createAccountPath := fmt.Sprintf("%s/organisation/accounts", apiVersion)
	err := s.http.Get(createAccountPath, accounts)
	if err != nil {
		return nil, err
	}
	pagedAccounts := pageAccounts(start, stop, *accounts)
	return pagedAccounts, err
}

// pageAccounts - paging here because the paging functionality of the downstream acount api isnt working
// (or maybe I couldn't get it to work, could't find any documentation of it anywhere)
// this is inefficient because the paging is not being done at the database level, instead
// i'm having to get all the account (could potentially be a huge list), dereference
// and then splicing it. I'm still quite new to golang but this, to my mind doesn't feel performant
func pageAccounts(start, stop int, accounts []models.Account) []models.Account {
	if start > len(accounts) {
		start = len(accounts)
	}

	if stop > len(accounts) {
		stop = len(accounts)
	}

	return accounts[start:stop]
}

func validateInjectedClientOrDefault(h *infrastructure.HTTP) *infrastructure.HTTP {
	if h == nil {
		h = infrastructure.NewHTTP(nil, nil, nil, "")
	}
	if h.BaseURL == nil {
		h.BaseURL, _ = url.Parse(defaultBaseURL)
	}

	if h.Client == nil {
		h.Client = http.DefaultClient
		h.Client.Timeout = 100 * time.Second
	}

	if h.Context == nil {
		h.Context = context.TODO()
	}
	h.UserAgent = userAgent

	return h
}
