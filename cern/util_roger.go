package cern

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/dpotapov/go-spnego"
)

// Roger client manages requests to the Roger service
type Roger struct {
	URL *url.URL
}

type RogerRequest struct {
	Hostname   string `json:"hostname"`
	AppState   string `json:"appstate,omitempty"`
	Expires    string `json:"expires,omitempty"`
	Message    string `json:"message,omitempty"`
	AppAlarmed *bool  `json:"app_alarmed,omitempty"`
	HwAlarmed  *bool  `json:"hw_alarmed,omitempty"`
	NcAlarmed  *bool  `json:"nc_alarmed,omitempty"`
	OsAlarmed  *bool  `json:"os_alarmed,omitempty"`
}

type RogerResponse struct {
	RogerRequest
	ExpiresDt       string `json:"expires_dt"`
	UpdateTime      string `json:"update_time"`
	UpdateTimeDt    string `json:"update_time_dt"`
	UpdatedBy       string `json:"updated_by"`
	UpdatedByPuppet bool   `json:"updated_by_puppet"`
}

// NewRogerClient constructs a new client configuration
func NewRogerClient(endpoint string) (*Roger, error) {
	url, err := url.Parse(endpoint)

	if err != nil {
		return nil, err
	}

	url, err = resolveURL(url)
	if err != nil {
		return nil, err
	}

	return &Roger{
		URL: url,
	}, nil
}

// Get roger state for a given host
func (r Roger) Get(ctx context.Context, hostname string) (*RogerResponse, error) {
	rogerRequest := RogerRequest{
		Hostname: hostname,
	}
	return r.do(ctx, rogerRequest, "GET")
}

// Create a roger state
func (r Roger) Create(ctx context.Context, rogerRequest RogerRequest) error {
	_, err := r.do(ctx, rogerRequest, "POST")
	return err
}

// Update a roger state
func (r Roger) Update(ctx context.Context, rogerRequest RogerRequest) error {
	_, err := r.do(ctx, rogerRequest, "PUT")
	return err
}

// Update a roger state
func (r Roger) Delete(ctx context.Context, hostname string) (*RogerResponse, error) {
	rogerRequest := RogerRequest{
		Hostname: hostname,
	}
	return r.do(ctx, rogerRequest, "DELETE")
}

// Do a request to roger
func (r Roger) do(ctx context.Context, rogerRequest RogerRequest, method string) (*RogerResponse, error) {
	client := http.Client{
		Transport: &spnego.Transport{},
	}

	url := fmt.Sprintf("%s/roger/v1/state/%s/", r.URL, rogerRequest.Hostname)
	if method == "POST" {
		url = fmt.Sprintf("%s/roger/v1/state/", r.URL)
	}
	log.Printf("[DEBUG] Request url constructed as follows: %s", url)

	var requestData []byte
	if method != "GET" && method != "DELETE" {
		requestData, _ = json.Marshal(rogerRequest)
		log.Printf("[DEBUG] Request data: %s", string(requestData[:]))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}
	if method != "GET" && method != "DELETE" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "deflate")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rogerResponse RogerResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, HTTPError{url, resp.StatusCode, string(body[:])}
	}

	if method != "POST" && method != "DELETE" && method != "PUT" {
		err = json.Unmarshal(body, &rogerResponse)
		if err != nil {
			return nil, err
		}
	}

	return &rogerResponse, nil
}
