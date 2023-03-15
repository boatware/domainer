package domainer

import "testing"

var tests = []struct {
	name     string
	domain   string
	expected URL
}{
	{
		"Parse full domain with every part", "https://www.example.com:443/search?q=hello+world#test", URL{
			FullURL:   "https://www.example.com:443/search?q=hello+world#test",
			Protocol:  "https",
			Subdomain: "www",
			Domain:    "example",
			TLD:       "com",
			Port:      443,
			Path:      "/search",
			Query: []Query{
				{
					Key:   "q",
					Value: "hello+world",
				},
			},
			Fragment: "test",
		},
	},
	{
		"Parse full URL with multipart TLD", "https://www.example.co.uk:443/search?q=hello+world#test", URL{
			FullURL:   "https://www.example.co.uk:443/search?q=hello+world#test",
			Protocol:  "https",
			Subdomain: "www",
			Domain:    "example",
			TLD:       "co.uk",
			Port:      443,
			Path:      "/search",
			Query: []Query{
				{
					Key:   "q",
					Value: "hello+world",
				},
			},
			Fragment: "test",
		},
	},
	{
		"Parse full URL with no subdomain", "https://example.com:443/search?q=hello+world#test", URL{
			FullURL:  "https://example.com:443/search?q=hello+world#test",
			Protocol: "https",
			Domain:   "example",
			TLD:      "com",
			Port:     443,
			Path:     "/search",
			Query: []Query{
				{
					Key:   "q",
					Value: "hello+world",
				},
			},
			Fragment: "test",
		},
	},
	{
		"Parse simple URL", "https://example.com", URL{
			FullURL:  "https://example.com",
			Protocol: "https",
			Domain:   "example",
			TLD:      "com",
		},
	},
	{
		"Parse URL with no protocol given", "example.com", URL{
			FullURL: "example.com",
			Domain:  "example",
			TLD:     "com",
		},
	},
	{
		"Parse URL with username and password", "user:pass@example.com", URL{
			FullURL:  "user:pass@example.com",
			Domain:   "example",
			TLD:      "com",
			Username: "user",
			Password: "pass",
		},
	},
	{
		"Parse URL with only username", "user@example.com", URL{
			FullURL:  "user@example.com",
			Domain:   "example",
			TLD:      "com",
			Username: "user",
		},
	},
	{
		"Parse URL with username and port", "user@example.com:80", URL{
			FullURL:  "user@example.com:80",
			Domain:   "example",
			TLD:      "com",
			Username: "user",
			Port:     80,
		},
	},
}

func TestFromString(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := FromString(tt.domain)
			if err != nil {
				t.Error(err)
			}
			if d.FullURL != tt.expected.FullURL {
				t.Errorf("FullURL: Expected '%s', got '%s'", tt.expected.FullURL, d.FullURL)
			}
			if d.Protocol != tt.expected.Protocol {
				t.Errorf("Protocol: Expected '%s', got '%s'", tt.expected.Protocol, d.Protocol)
			}
			if d.Subdomain != tt.expected.Subdomain {
				t.Errorf("Subdomain: Expected '%s', got '%s'", tt.expected.Subdomain, d.Subdomain)
			}
			if d.Domain != tt.expected.Domain {
				t.Errorf("URL: Expected '%s', got '%s'", tt.expected.Domain, d.Domain)
			}
			if d.TLD != tt.expected.TLD {
				t.Errorf("TLD: Expected '%s', got '%s'", tt.expected.TLD, d.TLD)
			}
			if d.Port != tt.expected.Port {
				t.Errorf("Port: Expected %d, got %d", tt.expected.Port, d.Port)
			}
			if d.Path != tt.expected.Path {
				t.Errorf("Path: Expected '%s', got '%s'", tt.expected.Path, d.Path)
			}
			if len(d.Query) != len(tt.expected.Query) {
				t.Errorf("Query Length: Expected %d, got %d", len(tt.expected.Query), len(d.Query))
			}
			for i, q := range d.Query {
				if q.Key != tt.expected.Query[i].Key {
					t.Errorf("\tQuery #%d: Expected '%s', got '%s'", i, tt.expected.Query[i].Key, q.Key)
				}
				if q.Value != tt.expected.Query[i].Value {
					t.Errorf("\tQuery #%d: Expected '%s', got '%s'", i, tt.expected.Query[i].Value, q.Value)
				}
			}
			if d.Fragment != tt.expected.Fragment {
				t.Errorf("Fragment: Expected '%s', got '%s'", tt.expected.Fragment, d.Fragment)
			}
			if d.Username != tt.expected.Username {
				t.Errorf("Username: Expected '%s', got '%s'", tt.expected.Username, d.Username)
			}
			if d.Password != tt.expected.Password {
				t.Errorf("Password: Expected '%s', got '%s'", tt.expected.Password, d.Password)
			}
		})
	}
}
