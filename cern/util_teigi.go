package cern

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/dpotapov/go-spnego"
)

// SecretRequest defines the Teigi request structure
type SecretRequest struct {
	Secret   string `json:"secret"`
	Encoding string `json:"encoding,omitempty"`
}

// SecretResponse defines the Teigi response structure
type SecretResponse struct {
	SecretRequest
	Hostgroup     string `json:"hostgroup"`
	Skey          string `json:"skey"`
	UpdateTime    string `json:"update_time"`
	UpdateTimeStr string `json:"update_time_str"`
	UpdatedBy     string `json:"updated_by"`
}

func SerializeHostgroup(hostgroup string) string {
	hostgroup = strings.Trim(hostgroup, "/")
	return strings.ReplaceAll(hostgroup, "/", "-")
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

func fqndify(host string) (string, error) {
	fqdn := host
	if strings.ContainsAny(fqdn, ".") {
		fqdn = host + ".cern.ch"
	}

	addrs, err := net.LookupHost(fqdn)
	if err != nil {
		return "", err
	}

	if len(addrs) == 0 {
		return "", fmt.Errorf("fqdn '%s' does not resolve", fqdn)
	} else if len(addrs) > 2 {
		return "", fmt.Errorf("fqdn '%s' may be an alias", fqdn)
	}

	return fqdn, nil
}

func (t Teigi) Create(ctx context.Context, scope string, entity, key string, secretRequest SecretRequest) error {
	_, err := t.do(ctx, "POST", scope, entity, key, secretRequest)
	return err
}

func (t Teigi) Delete(ctx context.Context, scope string, entity, key string) error {
	_, err := t.do(ctx, "DELETE", scope, entity, key, SecretRequest{})
	return err
}

func (t Teigi) Get(ctx context.Context, scope string, entity, key string) (*SecretResponse, error) {
	return t.do(ctx, "GET", scope, entity, key, SecretRequest{})
}

// Get request
func (t Teigi) do(ctx context.Context, method string, scope string, entity string, key string, secretRequest SecretRequest) (*SecretResponse, error) {
	var err error
	client := http.Client{
		Transport: &spnego.Transport{},
	}

	if scope == "host" {
		entity, err = fqndify(entity)
		if err != nil {
			return nil, err
		}
	} else if scope == "hostgroup" {
		entity = SerializeHostgroup(entity)
	}

	url := fmt.Sprintf("%s/tbag/v2/%s/%s/secret/%s/", t.URL, scope, entity, key)
	log.Printf("[DEBUG] Request url constructed as follows: %s", url)

	var requestData []byte
	if method == "POST" {
		requestData, _ = json.Marshal(secretRequest)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "deflate")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non 2xx error code (%d)", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var secretResponse SecretResponse
	if method == "GET" {
		err = json.Unmarshal(body, &secretResponse)
		if err != nil {
			return nil, err
		}
	}

	return &secretResponse, nil
}
