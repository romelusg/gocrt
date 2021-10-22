package main

import (
    "testing"
    "reflect"
    "os"
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
        {"mail@example.org", "mail%40example.org"},
        {"example-23894723\"!§&%§=()=$.org", "example-23894723%22%21%C2%A7%26%25%C2%A7%3D%28%29%3D%24.org"},
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
/*func TestGetCrtShJson(t *testing.T) {
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
}*/

// TEST extractSubodmainsFromJson()
func TestExtractSubodmainsFromJson(t *testing.T) {
    data := `[{"issuer_ca_id":185756,"issuer_name":"C=US, O=DigiCert Inc, CN=DigiCert TLS RSA SHA256 2020 CA1","common_name":"dev.example.com\nexample.com\nproducts.example.com\nsupport.example.com\nwww.example.com","name_value":"mail.example.com\ndev.example.com","id":3704614715,"entry_timestamp":"2020-11-27T13:49:06.706","not_before":"2020-11-24T00:00:00","not_after":"2021-12-25T23:59:59","serial_number":"0fbe08b0854d05738ab0cce1c9afeec9"},{"issuer_ca_id":185756,"issuer_name":"C=US, O=DigiCert Inc, CN=DigiCert TLS RSA SHA256 2020 CA1","common_name":"www.example.org","name_value":"test.example.com","id":3704614715,"entry_timestamp":"2020-11-27T13:49:06.706","not_before":"2020-11-24T00:00:00","not_after":"2021-12-25T23:59:59","serial_number":"0fbe08b0854d05738ab0cce1c9afeec9"}]`
    want := []string{
        "dev.example.com", "example.com", "products.example.com",
        "support.example.com", "www.example.com", "mail.example.com",
        "www.example.org", "test.example.com",
    }
    have := extractSubdomainsFromJson(data)

    if ! reflect.DeepEqual(have, want) {
        t.Errorf("Want subdomain list of '%s'; have '%s'",
            want, have)
    }
}

// TEST saveSubdomains()
func TestSaveSubdomains(t *testing.T) {
    data := []string {
        "test.example.com",
        "dev.example.com",
        "hello.example.com",
    }
    fileFlags := os.O_CREATE|os.O_RDWR|os.O_TRUNC
    have := saveSubdomains("testfiles-gocrt", "example.com", data, fileFlags)

    if ! have {
        t.Errorf("Could not write subdomains to file")
    }
}

// TEST filterInvalidDomains()
func TestFilterInvalidDomains(t *testing.T) {
    data := []string {
        "test.example.com",
        "dev.example.com",
        "fail/test.org",
        "invalid@example.com",
        "hello.example.com",
    }
    want := []string {
        "test.example.com",
        "dev.example.com",
        "hello.example.com",
    }
    have := filterInvalidDomains(data)

    if ! reflect.DeepEqual(have, want) {
        t.Errorf("Invalid domain in list %s", have)
    }
}
