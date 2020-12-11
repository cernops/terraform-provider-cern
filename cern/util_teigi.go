package cern

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gitlab.cern.ch/batch-team/negotiate"
)

// Secret defines the Teigi response structure
type Secret struct {
	Encoding      string `json:"encoding"`
	Hostgroup     string `json:"hostgroup"`
	Secret        string `json:"secret"`
	Skey          string `json:"skey"`
	UpdateTime    string `json:"update_time"`
	UpdateTimeStr string `json:"update_time_str"`
	UpdatedBy     string `json:"updated_by"`
}

// Teigi client manages requests to the Teigi service
type Teigi struct {
	URL *url.URL
}

// NewTeigiClient constructs a new client configuration
func NewTeigiClient(endpoint string) (*Teigi, error) {
	url, err := url.Parse(endpoint)

	if err != nil {
		return nil, err
	}

	url, err = resolveURL(url)
	if err != nil {
		return nil, err
	}

	return &Teigi{
		URL: url,
	}, nil
}

// Get request
func (t Teigi) Get(hostgroup string, key string) (*Secret, string, error) {
	client := http.Client{}

	url := fmt.Sprintf("%s/tbag/v2/hostgroup/%s/secret/%s/", t.URL, strings.ReplaceAll(hostgroup, "/", "-"), key)
	log.Printf("[DEBUG] Request url constructed as follows: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "Failed to create an http request", err
	}

	// Decorate the request with Kerberos credentials
	err = negotiate.Negotiate(req)
	if err != nil {
		return nil, "Teigi request requires authorization", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "Failed to do an http Get", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "Failed to read the body of the response", err
	}

	secret := new(Secret)
	err = json.Unmarshal(body, &secret)
	if err != nil {
		return nil, "Failed to unmarshal response body. Check that the Teigi key is valid", err
	}

	return secret, "Success", nil
}
