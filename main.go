package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

var gocrtVersion = "dev"

func init() {
	flag.Usage = func() {
		h := "A crt.sh command line client written in golang.\n\n"

		h += "Usage:\n"
		h += "  gocrt [OPTIONS] [FILE|URL|-]\n\n"

		h += "Options:\n"
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

func main() {	
	// "version" command line argument
	var version bool
	flag.BoolVar(&version, "version" false, "")
	
	// "output" command line argument
	var output string
	flag.StringVar(&output, "output", "gocrt", "")
	flag.StringVar(&output, "o", "gocrt", "")

	flag.Parse()
	
	// Print version
	if version {
		fmt.Printf("gocrt version %s\n", gocrtVersion)
		os.Exit(0)
	}
	
	// TESTING
	fmt.Printf("-o or --output argument: $s\n", output)
}
