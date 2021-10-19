# gocrt
`gocrt` is a command line client for [crt.sh](https://crt.sh/) written in golang.
```text
$ gocrt -o github-demo < domains.txt
Get subdomains from: test.de
Get subdomains from: example.com
Save subdomains from: test.de -> saved
Save subdomains from: example.com -> saved
[✓] Done

$ ls -ls github-demo
drwxr-xr-x 2 tom tom 4096 16. Okt 14:29 .
drwxr-xr-x 5 tom tom 4096 16. Okt 14:29 ..
-rw-r--r-- 1 tom tom  133 16. Okt 14:29 example.com
-rw-r--r-- 1 tom tom  473 16. Okt 14:29 test.de
```

## Installation
If you've got Go installed and configured you can install `gocrt` with:
```bash
$ go get -u github.com/tomschwarz/gocrt
```

## Usage 
Get domain from `command line`:
```bash
$ gocrt example.com
```

Get domain from `stdin`:
```bash
$ cat domains.txt | gocrt
# OR
$ gocrt < domains.txt 
```

Pipe found subdomains to other tools:
```bash
$ gocrt -s < domains.txt | httprobe
# OR
$ cat domains.txt | gocrt -s | httprobe
# OR
$ gocrt --stdout example.com | httprobe
```

Store subdomains to custom directory:
```bash
$ cat domains.txt | gocrt -o my-custom-dir 
# OR
$ gocrt --output my-custom-dir < domains.txt
```

## Get Help
```text
$ gocrt --help
gocrt is a command line client for crt.sh written in golang.

Usage:
  gocrt [OPTIONS] [FILE|URL|-]

Options:
  -h, --help       Print usage informations
  -o, --output     Output directory for all found subdomains of given domains
  -c, --combine    Additionally combine output for all found subdomains of given domains in one file
  -s, --stdout     Print only subdomains to STDOUT so they can be piped to other tools, they also get saved into files
      --version    Print version information

Examples:
  cat domains.txt | gocrt -o domains-crt
  gocrt -o domains-crt example.com 
  gocrt < domains.txt
  gocrt example.com
```

## TODO - Release script
```bash
# To update version while "go build"
go build -ldflags "-X main.gocrtVersion=<VERSION>"
```
