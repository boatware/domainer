package domainer

import (
	"golang.org/x/net/publicsuffix"
	"net"
	"strconv"
	"strings"
)

// Query is a key-value pair used in a URL query string.
type Query struct {
	// Key is the key of the query.
	// Example: "q" in "https://example.com/search?q=hello+world"
	Key string `json:"key"`

	// Value is the value of the query.
	// Example: "hello+world" in "https://example.com/search?q=hello+world"
	Value string `json:"value"`
}

// URL is a split of a given domain name.
type URL struct {
	// FullURL represents the full domain name this struct has been created with.
	// Example: "https://www.example.com:443/search?q=hello+world#test"
	FullURL string `json:"full_url"`

	// Protocol represents the protocol used to access the domain.
	// Example: "https" in "https://www.example.com:443/search?q=hello+world#test"
	Protocol string `json:"protocol"`

	// Subdomain represents the subdomain of the domain.
	// Example: "www" in "https://www.example.com:443/search?q=hello+world#test"
	Subdomain string `json:"subdomain"`

	// Hostname represents the hostname of the domain.
	// Example: "example.com" in "https://www.example.com:443/search?q=hello+world#test"
	Hostname string `json:"hostname"`

	// Domain represents the domain name (or second level domain).
	// Example: "example" in "https://www.example.com:443/search?q=hello+world#test"
	Domain string `json:"domain"`

	// TLD represents the top level domain.
	// Example: "com" in "https://www.example.com:443/search?q=hello+world#test"
	TLD string `json:"tld"`

	// Port represents the port used to access the domain.
	// Example: 443 in "https://www.example.com:443/search?q=hello+world#test"
	Port int `json:"port"`

	// Path represents the path used to access the domain.
	// Example: "/search" in "https://www.example.com:443/search?q=hello+world#test"
	Path string `json:"path"`

	// Query represents the query used to access the domain.
	// Example: []Query{{"q", "hello+world"}} in "https://www.example.com:443/search?q=hello+world#test"
	Query []Query `json:"query"`

	// Fragment represents the fragment used to access the domain.
	// Example: "test" in "https://www.example.com:443/search?q=hello+world#test"
	Fragment string `json:"fragment"`

	// Username represents the username used to access the domain.
	// Example: "user" in "https://user:pass@example.com:443/search?q=hello+world#test"
	Username string `json:"username"`

	// Password represents the password used to access the domain.
	// Example: "pass" in "https://user:pass@example.com:443/search?q=hello+world#test"
	Password string `json:"password"`

	// IPAddress represents the IP address the domain resolves to.
	// Example: "127.0.0.1" (obviously not a real server IP address)
	IPAddress string `json:"ip_address"`
}

// FromString parses a given domain name and returns a URL struct.
//
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

	u.Hostname = tldPlusOne

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

	// Get the IP address
	ip, err := net.LookupIP(u.Hostname)
	if err != nil {
		return nil, err
	}
	u.IPAddress = ip[0].String()

	return u, nil
}
