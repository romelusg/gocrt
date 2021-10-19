# gocrt
A crt.sh command line client written in golang.

```bash
$ cat domains.txt | gocrt -o domains-crt 

$ gocrt -o domains-crt example.com

$ gocrt example.com

$ gocrt < domains.txt 
```

```bash
# To update version while "go build"
go build -ldflags "-X main.gocrtVersion=<VERSION>"
```
