package helper

import (
	"net"
	urlnet "net/url"

	d "github.com/bobesa/go-domain-util/domainutil"
)

// Domain gets the top level domain from url
func Domain(url string) string {

	// doesn't work with port so we have to do this...
	host, _ := hostFromURL(url)
	return d.Domain(host)
}

// DomainSuffix gets the domain suffix from url
func DomainSuffix(url string) string {
	// doesn't work with port so we have to do this...
	host, _ := hostFromURL(url)
	return d.DomainSuffix(host)
}

// HasSubdomain checks if url has subdomain
func HasSubdomain(url string) bool {
	// doesn't work with port so we have to do this...
	host, _ := hostFromURL(url)
	return d.HasSubdomain(host)
}

// Subdomain gets the subdomain from url
func Subdomain(url string) string {
	// doesn't work with port so we have to do this...
	host, _ := hostFromURL(url)
	return d.Subdomain(host)
}

// Protocol gets the protocol from url
func Protocol(url string) string {
	return d.Protocol(url)
}

// Username gets the username from credentials of url
func Username(url string) string {
	return d.Username(url)
}

// Password get the password from credentials of url
func Password(url string) string {
	return d.Password(url)
}

func hostFromURL(url string) (string, error) {
	u, err := urlnet.Parse(url)
	if err != nil {
		return url, err
	}

	host, _, err := net.SplitHostPort(u.Host)
	if err != nil || host == "" {
		return url, err
	}

	return host, nil
}