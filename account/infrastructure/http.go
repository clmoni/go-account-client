package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type (
	// HTTP - a client that will be used as the basis of a request
	// Naming this with protocol as opposed to the "client"
	// feels more expressive & reduces stutter
	HTTP struct {
		BaseURL   *url.URL
		UserAgent string
		Client    *http.Client
		Context   context.Context
	}

	request struct {
		Data interface{} `json:"data,omitempty"`
	}

	response struct {
		Data         interface{} `json:"data,omitempty"`
		ErrorMessage string      `json:"error_message,omitempty"`
	}
)

// NewHTTP - returns a new client. If a nil httpClient is
// provided the default httpClient will be used.
func NewHTTP(ctx context.Context, httpClient *http.Client, baseURL *url.URL, userAgent string) *HTTP {
	h := &HTTP{
		Client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
		Context:   ctx,
	}
	return h
}

// Get - make get request to the account api
// v: response
func (h *HTTP) Get(path string, v interface{}) error {
	return h.makeRequest("GET", path, nil, v)
}

// Delete - make get request to the account api
// v: response
func (h *HTTP) Delete(path string) error {
	return h.makeRequest("DELETE", path, nil, nil)
}

// Post - make post request to the account api
// data: body of request
// v: response
func (h *HTTP) Post(path string, data interface{}, v interface{}) error {
	return h.makeRequest("POST", path, data, v)
}

func (h *HTTP) makeRequest(method string, path string, data interface{}, v interface{}) error {
	urlStr := path
	rel, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	u := h.BaseURL.ResolveReference(rel)

	resp, err := h.do(u.String(), method, data)
	defer resp.Body.Close()

	err = unmarshalResponse(resp, v)

	return err
}

func (h *HTTP) do(url, method string, data interface{}) (*http.Response, error) {
	req, err := h.createRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	resp, err := h.Client.Do(req.WithContext(h.Context))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *HTTP) createRequest(method, url string, data interface{}) (*http.Request, error) {
	var body io.Reader
	if data != nil {
		b, err := json.Marshal(request{Data: data})
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", h.UserAgent)
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func unmarshalResponse(r *http.Response, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	bodyString := string(body)
	if len(bodyString) > 0 {
		res := &response{Data: v}
		err = json.Unmarshal([]byte(bodyString), res)
		if len(res.ErrorMessage) > 0 {
			err = fmt.Errorf("downstream api error: %s", res.ErrorMessage)
		}
	}

	return err
}
