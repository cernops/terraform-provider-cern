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

// Certmr client manages requests to the Certmgr service
type CertMgr struct {
	URL *url.URL
}

type CertMgrRequest struct {
	Hostname string `json:"hostname"`
}
type CertMgrResponse struct {
	Hostname  string `json:"hostname"`
	Id        int    `json:"id"`
	Requestor string `json:"requestor"`
	Start     string `json:"start"`
	End       string `json:"End"`
}

// NewCertMgrClient constructs a new client configuration
func NewCertMgrClient(endpoint string) (*CertMgr, error) {
	url, err := url.Parse(endpoint)

	if err != nil {
		return nil, err
	}

	url, err = resolveURL(url)
	if err != nil {
		return nil, err
	}

	return &CertMgr{
		URL: url,
	}, nil
}

// Get a new certficiate for a given host
func (c CertMgr) Do(ctx context.Context, hostname string) (*CertMgrResponse, error) {
	client := http.Client{
		Transport: &spnego.Transport{},
	}

	url := fmt.Sprintf("%s/krb/certmgr/staged/", c.URL)
	log.Printf("[DEBUG] Request url constructed as follows: %s", url)

	requestData, _ := json.Marshal(CertMgrRequest{
		Hostname: hostname,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non 2xx error code (%d)", resp.StatusCode)
	}

	var certMgrResponse CertMgrResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &certMgrResponse)
	if err != nil {
		return nil, err
	}

	return &certMgrResponse, nil
}
