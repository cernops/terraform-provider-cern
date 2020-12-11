package cern

import (
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strings"
	"time"
)

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
