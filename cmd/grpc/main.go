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
		path, err := downloadProto(workspace, p)
		if err != nil {
			log.Fatalf("Cannot download %s: %v", p, err)
		}
		protos = append(protos, path)
	}
	if err := gen(workspace, protos); err != nil {
		log.Fatal(err)
	}
}

func downloadProto(workspace, url string) (string, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		// local file
		return url, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return "", nil
	}
	defer res.Body.Close()

	f, err := ioutil.TempFile(workspace, "protobuf")
	if err != nil {
		return "", nil
	}
	if _, err = io.Copy(f, res.Body); err != nil {
		return "", err
	}
	return f.Name(), f.Close()
}
