package main

import (
    "flag"
    "strings"
    "fmt"
    "os"
    "log"
    "bufio"
    "sync"
    "io/ioutil"
    "net/url"
    "net/http"
    "encoding/json"
)

// gocrt version
var gocrtVersion = "0.0.1-dev"

// make given list unique
func unique(list []string) ([]string) {
    allKeys := make(map[string]bool)
    uniqueList := []string{}
    for _, item := range list{
        if _, value := allKeys[item]; !value {
            allKeys[item] = true
            uniqueList = append(uniqueList, item)
        }
    }
    return uniqueList
}

// Extract domain from possible link
func extractDomain(link string) (string) {
    url, err := url.Parse(link)

    if len(url.Hostname()) != 0 && err == nil {
        link = url.Hostname()
    }

    link = strings.TrimSpace(link)
    return strings.ToLower(link)
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

    return unique(domains), err
}

// Get JSON-data from crt.sh
func getCrtShJson(domain string) (string) {
    crtUrl := fmt.Sprintf("https://crt.sh?q=%s&output=json", domain)

    response, err := http.Get(crtUrl)
    if err != nil {
        log.Fatal(err)
    }

    defer response.Body.Close()
    data, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    return string(data)
}

// Extract subdomains from JSON-data
func extractSubdomainsFromJson(jsonData string) ([]string){
    var subdomains []string
    var entries []map[string]interface{}

    json.Unmarshal([]byte(jsonData), &entries)
    for _, entry := range entries {
        commonName := entry["common_name"].(string)
        commonNameList := strings.Split(commonName, "\n")
        subdomains = append(subdomains, commonNameList...)

        nameValue := entry["name_value"].(string)
        nameValueList := strings.Split(nameValue, "\n")
        subdomains = append(subdomains, nameValueList...)
    }

    return unique(subdomains)
}

// Write subdomains into file
func saveSubdomains(dir string, domain string, subdomains []string) (bool){
    outputDir := "./" + dir
    os.Mkdir(outputDir, 0755)

    filePath := outputDir + "/" + domain
    file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

    if err != nil {
        log.Fatalf("Could not create file to save subdomains: %s", err)
        return false
    }

    writer := bufio.NewWriter(file)
    for _, subdomain := range subdomains {
        writer.WriteString(subdomain + "\n")
    }
    writer.Flush()
    file.Close()

    return true
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
    flag.StringVar(&output, "output", "subdomains", "")
    flag.StringVar(&output, "o", "subdomains", "")

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

    // Get JSON-data from crt.sh
    var worker sync.WaitGroup
    for _, domain := range domains {
        worker.Add(1)
        go func(domain string) {
            defer worker.Done()

            fmt.Printf("Get subdomains from: %s\n", domain)
            jsonData := getCrtShJson(domain)
            subdomains := extractSubdomainsFromJson(jsonData)

            fmt.Printf("Save subdomains from: %s", domain)
            saved := saveSubdomains(output, domain, subdomains)
            if saved {
                fmt.Printf(" -> saved\n")
            }
        }(domain)
    }
    worker.Wait()
    fmt.Printf("[\u2713] Done\n")
}
