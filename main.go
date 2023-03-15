package main

import (
	"golang.org/x/net/publicsuffix"
	"strconv"
	"strings"
)

type Query struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// URL is a split of a given domain name.
type URL struct {
	// FullURL represents the full domain name this struct has been created with.
	FullURL string `json:"full_url"`

	// Protocol represents the protocol used to access the domain.
	Protocol string `json:"protocol"`

	// Subdomain represents the subdomain of the domain.
	Subdomain string `json:"subdomain"`

	// Domain represents the domain name.
	Domain string `json:"domain"`

	// TLD represents the top level domain.
	TLD string `json:"tld"`

	// Port represents the port used to access the domain.
	Port int `json:"port"`

	// Path represents the path used to access the domain.
	Path string `json:"path"`

	// Query represents the query used to access the domain.
	Query []Query `json:"query"`

	// Fragment represents the fragment used to access the domain.
	Fragment string `json:"fragment"`

	// Username represents the username used to access the domain.
	Username string `json:"username"`

	// Password represents the password used to access the domain.
	Password string `json:"password"`
}

//goland:noinspection HttpUrlsUsage
func FromString(url string) (*URL, error) {
	u := &URL{}

	// Set the full url, so we can work with the original value
	u.FullURL = url

	// Get the protocol
	// If the protocol is not set, we assume it's http
	if strings.HasPrefix(url, "http://") {
		u.Protocol = "http"
		url = strings.TrimPrefix(url, "http://")
	}

	if strings.HasPrefix(url, "https://") {
		u.Protocol = "https"
		url = strings.TrimPrefix(url, "https://")
	}

	// Find the first occurrence of a slash, which indicates the end of the url and the start of the path
	// If no slash is found, we assume the url is the full url
	slashIndex := strings.Index(url, "/")
	if slashIndex == -1 {
		slashIndex = len(url)
	}

	// Cut the url at the slash
	path := url[slashIndex:]
	url = url[:slashIndex]

	// Find the first occurence of an @, which indicates the end of the username and password and the start of the domain
	// If no @ is found, we assume the url is the full url
	atIndex := strings.Index(url, "@")
	if atIndex > -1 {
		// Cut the url at the @
		credentials := url[:atIndex]
		url = url[atIndex+1:]

		// Find the first occurence of a :, which indicates the end of the username and the start of the password
		// If no : is found, we assume the password is empty
		colonIndex := strings.Index(credentials, ":")
		if colonIndex == -1 {
			u.Username = credentials
		} else {
			u.Username = credentials[:colonIndex]
			u.Password = credentials[colonIndex+1:]
		}
	}

	// Find the first occurrence of a colon, which indicates the end of the url and the start of the port
	// If no colon is found, we assume the port is the default port for the protocol
	colonIndex := strings.Index(url, ":")
	if colonIndex == -1 {
		colonIndex = len(url)
	}

	// Cut the url at the colon
	port := url[colonIndex:]
	url = url[:colonIndex]

	// If the port is not empty, we convert it to an integer
	if port != "" {
		// Remove the colon
		port = strings.TrimPrefix(port, ":")
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}
		u.Port = p
	}

	// Find the first occurrence of a question mark, which indicates the end of the path and the start of the query
	// If no question mark is found, we assume the query is empty
	questionMarkIndex := strings.Index(path, "?")
	if questionMarkIndex == -1 {
		questionMarkIndex = len(path)
	}

	// Cut the path at the question mark
	query := path[questionMarkIndex:]
	path = path[:questionMarkIndex]

	// Before we go on, we can add the path to the url
	u.Path = path

	// Before we go on, we need to check if there's a fragment
	// Find the first occurrence of a hash, which indicates the end of the query and the start of the fragment
	// If no hash is found, we assume the fragment is empty
	hashIndex := strings.Index(query, "#")
	if hashIndex == -1 {
		hashIndex = len(query)
	}

	// Cut the query at the hash
	fragment := query[hashIndex:]
	query = query[:hashIndex]

	// Remove the hash
	fragment = strings.TrimPrefix(fragment, "#")

	// Before we go on, we can add the fragment to the url
	u.Fragment = fragment

	// If the query is not empty, we split it into key-value pairs
	if query != "" {
		// Remove the question mark
		query = strings.TrimPrefix(query, "?")

		// Split the query into key-value pairs
		queryParts := strings.Split(query, "&")

		// Iterate over the key-value pairs
		for _, queryPart := range queryParts {
			// Split the key-value pair into key and value
			queryPartParts := strings.Split(queryPart, "=")

			// If the query part contains a key and a value, we add it to the query
			if len(queryPartParts) == 2 {
				u.Query = append(u.Query, Query{
					Key:   queryPartParts[0],
					Value: queryPartParts[1],
				})
			}
		}
	}

	tldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(url)
	if err != nil {
		return nil, err
	}

	// Split the tldPlusOne into url and tld
	tldPlusOneParts := strings.Split(tldPlusOne, ".")
	tld := strings.Join(tldPlusOneParts[1:], ".")

	if tld != "" {
		u.TLD = tld
	}

	// Remove the tld from the url
	url = strings.TrimSuffix(url, "."+tld)

	// Now we can split the url into subdomain and url
	domainParts := strings.Split(url, ".")

	// The last part of the url is the url itself
	u.Domain = domainParts[len(domainParts)-1]

	// The rest of the url is the subdomain
	u.Subdomain = strings.Join(domainParts[:len(domainParts)-1], ".")

	return u, nil
}
