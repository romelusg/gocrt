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

func help() {

}

func main() {
	var help bool
	flag.BoolVar(&help, "h", false, "Prints the help page") 

	flag.Parse()

	if help {
	    help()
	    exit
	}
}
