package main

import (
    "flag"
    "fmt"
    "os"
    "log"
    "bufio"
    "net/url"
)

// gocrt version
var gocrtVersion = "dev"

// Extract domain from possible link
func extractDomain(link string) (string) {
    url, err := url.Parse(link)

    if len(url.Hostname()) != 0 && err == nil {
        link = url.Hostname()
    }

    return link
}

// Get domains from stdin/pipe/command line argument
func getDomains() ([]string, error) {
    var domain string
    var domains []string
    var err error

    if flag.NArg() == 0 { // read from stdin/pipe

        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            domains = append(domains, extractDomain(scanner.Text()))
        }

        err = scanner.Err()

    } else { // read from argument

        domain = extractDomain(os.Args[len(os.Args) - 1])
        domains = append(domains, domain)

    }

    // TODO: remove duplicates from "domains"
    return domains, err
}

// init, get called automatic before main()
func init() {
    flag.Usage = func() {
        h := "A crt.sh command line client written in golang.\n\n"

        h += "Usage:\n"
        h += "  gocrt [OPTIONS] [FILE|URL|-]\n\n"

        h += "Options:\n"
        h += "  -h, --help       Print usage informations\n"
        h += "  -o, --output     Output directory for all found subdomains of given domains\n"
        h += "      --version    Print version information\n"
        h += "\n"

        h += "Examples:\n"
        h += "  cat domains.txt | gocrt -o domains-crt\n"
        h += "  gocrt -o domains-crt example.com \n"
        h += "  gocrt < domains.txt\n"
        h += "  gocrt example.com\n"

        fmt.Fprintf(os.Stderr, h)
    }
}

// main, magic happens here
func main() {
    // "version" command line argument
    var version bool
    flag.BoolVar(&version, "version", false, "")

    // "output" command line argument
    var output string
    flag.StringVar(&output, "output", "gocrt", "")
    flag.StringVar(&output, "o", "gocrt", "")

    flag.Parse()

    // Print version
    if version {
        fmt.Printf("gocrt version: %s\n", gocrtVersion)
        os.Exit(0)
    }

    // Get domains to request
    domains, err := getDomains()
    if err != nil {
        log.Fatal(err)
        os.Exit(3)
    }
    fmt.Println(domains)
}
