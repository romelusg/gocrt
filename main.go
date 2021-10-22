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
    "regexp"
)

// gocrt version
var gocrtVersion = "0.0.1-dev"

// Wrapper for printing messages
// - used to print only subdomains to STDOUT 
func printMessage(silent bool, message string, arguments ...interface{}) {
    if ! silent {
        fmt.Printf(message, arguments...)
    }
}

// filter invalid domains/subdomains
func filterInvalidDomains(domains []string) ([]string) {
    var filtered []string
    regex, _ := regexp.Compile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`)
    for _, domain := range domains {
        if regex.Match([]byte(domain)) {
            filtered = append(filtered, domain)
        }
    }

    return filtered
}

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
    parsedUrl, err := url.Parse(link)

    if err == nil && len(parsedUrl.Hostname()) != 0 {
        link = parsedUrl.Hostname()
    }

    link = strings.TrimSpace(link)
    return url.QueryEscape(
        strings.ToLower(link))
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

    subdomains = filterInvalidDomains(subdomains)
    return unique(subdomains)
}

// Write subdomains into file
func saveSubdomains(dir string, domain string,
    subdomains []string, fileFlags int) (bool){
    outputDir := "./" + dir
    os.Mkdir(outputDir, 0755)

    filePath := outputDir + "/" + domain
    file, err := os.OpenFile(filePath, fileFlags, 0644)

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
        h := "gocrt is a command line client for crt.sh written in golang.\n\n"

        h += "Usage:\n"
        h += "  gocrt [OPTIONS] [FILE|URL|-]\n\n"

        h += "Options:\n"
        h += "  -h, --help       Print help/usage informations\n"
        h += "  -o, --output     Custom output directory for all found subdomains of given domains, DEFAULT: 'subdomains'\n"
        h += "  -c, --combine    Additionally combine output for all found subdomains of given domains in one file\n"
        h += "  -s, --stdout     Print only subdomains to STDOUT so they can be piped directly to other tools, they will not be saved into files\n"
        h += "      --version    Print version information\n"
        h += "\n"

        h += "Examples:\n"
        h += "  cat domains.txt | gocrt -o domains-crt\n"
        h += "  gocrt -o domains-crt example.com \n"
        h += "  gocrt < domains.txt\n"
        h += "  gocrt -s < domains.txt | tee combined.txt | httprobe\n"
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

    // "combine" command line argument
    var combine bool
    flag.BoolVar(&combine, "combine", false, "")
    flag.BoolVar(&combine, "c", false, "")

    // "stdout" command line argument
    var stdout bool
    flag.BoolVar(&stdout, "stdout", false, "")
    flag.BoolVar(&stdout, "s", false, "")

    flag.Parse()

    // Print version
    if version {
        printMessage(stdout, "gocrt version: %s\n", gocrtVersion)
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

            printMessage(stdout, "Get subdomains from: %s\n", domain)
            jsonData := getCrtShJson(domain)
            subdomains := extractSubdomainsFromJson(jsonData)

            if stdout { // Print subdomains to STDOUT
                for _, subdomain := range subdomains {
                    printMessage(!stdout, "%s\n", subdomain)
                }
            } else { // Print messages to STDOUT and save subdomains to files
                printMessage(stdout, "Save subdomains from: %s", domain)
                saved := saveSubdomains(output, domain,
                    subdomains, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
                if saved {
                    printMessage(stdout, " -> saved\n")
                }

                if combine { // additionally combine all domains
                    saveSubdomains(output, "combined.txt", subdomains,
                        os.O_CREATE|os.O_RDWR|os.O_APPEND)
                }
            }
        }(domain)
    }
    worker.Wait()

    if combine {
        printMessage(stdout, "Additionally saved subdomains combined in one file\n")
    }
    printMessage(stdout, "[\u2713] Done\n")
}
