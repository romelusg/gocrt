package main

import (
    "testing"
)

// TEST unique()
func TestUnique(t *testing.T) {
    list := []string{
        "example.com",
        "example.org",
        "test.de",
        "example.com",
    }
    want := len(list) - 1
    have := unique(list)
    if len(have) != want {
        t.Errorf("Want an count of '%d' for unique(); have '%d'",
            want, len(have))
    }
}

// TEST extractDomain()
func TestExtractDomain(t *testing.T) {
    tests := []struct {
        url string
        want string 
    }{
        {"http://example.com", "example.com"},
        {"http://example.com", "example.com"},
        {"https://example.org:1234", "example.org"},
        {"mail@example.org", "mail@example.org"},
        {"https://example.org:1234/something/here", "example.org"},
        {"https://example.org/test/there", "example.org"},
    }

    for _, test := range tests {
        have := extractDomain(test.url)
        if have != test.want {
            t.Errorf("Want '%s' for extractDomain('%s'); have '%s'",
                test.want, test.url, have)
        }
    }
}

// TEST getDomains()
func TestGetDomains(t *testing.T) {
    // TODO: improve test!
    want := 0

    have, err := getDomains()
    if len(have) != want || err != nil{
        t.Errorf("Want an count of '%d' for getDomains(); have '%d'",
            want, len(have))
    }
}

// TEST getCrtShJson()
func TestGetCrtShJson(t *testing.T) {
    domains := []struct{
        url string
        want bool
    }{
        {"example.com", true},
        {"doesnotexist1234heheheheexample.com", false},
    }

    for _, domain := range domains {
        have := getCrtShJson(domain.url)
        if have == "[]" && domain.want {
            t.Errorf("Want JSON-Data from domain '%s'; have '%s'",
                domain.url, have)
        }
    }
}
