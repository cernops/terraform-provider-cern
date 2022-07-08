package cern

import (
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HTTPError struct {
	URL        string
	StatusCode int
	RespBody   string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf(
		"HTTP Error:{\n"+
			"  url:        [%s]\n"+
			"  statusCode: [%d]\n"+
			"  respBody:   [%s]\n"+
			"}",
		e.URL,
		e.StatusCode,
		e.RespBody,
	)
}

// Taken from terraform-openstack-provider
// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, prefix string, err error) error {
	if httpError, ok := err.(HTTPError); ok && httpError.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("%s: %s", prefix, err)
}
func resolveURL(u *url.URL) (*url.URL, error) {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	// Take a url, extract the hostname, resolve to canonical
	addrs, err := net.LookupHost(u.Hostname())
	if err != nil {
		return nil, err
	}
	port := u.Port()
	randAddr := addrs[rand.Intn(len(addrs))]
	hosts, err := net.LookupAddr(randAddr)
	if err != nil {
		return nil, err
	}
	newHost := strings.TrimSuffix(hosts[rand.Intn(len(hosts))], ".")
	if port != "" {
		newHost = fmt.Sprintf("%s:%s", newHost, port)
	}
	u.Host = newHost
	return u, nil
}
