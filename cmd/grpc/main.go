package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const help = `Usage: grpc gen protos...
`

var o = flag.String("o", ".", "output directory")

func main() {
	flag.Parse()

	if *o != "." {
		if err := os.MkdirAll(*o, 0755); err != nil {
			log.Fatalf("Cannot create output dir: %v", err)
		}
	}

	workspace, err := ioutil.TempDir("", "grpc")
	if err != nil {
		log.Fatalf("Cannot create temp workspace: %v", err)
	}
	defer os.RemoveAll(workspace)

	var protos []string
	for _, p := range flag.Args() {
		path, err := prepare(workspace, p)
		if err != nil {
			log.Fatalf("Cannot download %s: %v", p, err)
		}
		protos = append(protos, path)
	}
	if err := gen(workspace, protos); err != nil {
		log.Fatal(err)
	}
}

func prepare(workspace, url string) (string, error) {
	var r io.Reader
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		res, err := http.Get(url)
		if err != nil {
			return "", nil
		}
		defer res.Body.Close()
		r = res.Body
	} else {
		f, err := os.Open(url)
		if err != nil {
			return "", err
		}
		r = f
	}
	f, err := ioutil.TempFile(workspace, "protobuf")
	if err != nil {
		return "", nil
	}
	if _, err = io.Copy(f, r); err != nil {
		return "", err
	}
	return f.Name(), f.Close()
}
