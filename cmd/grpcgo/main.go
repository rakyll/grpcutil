/*
 *
 * Copyright 2017, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

// +build darwin, linux

package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var o = flag.String("o", ".", "output directory")
var i = flag.String("I", "", "include directory")

func main() {
	flag.Parse()

	if *o != "." {
		if err := os.MkdirAll(*o, 0755); err != nil {
			log.Fatalf("Cannot create output dir: %v", err)
		}
	}

	var includes []string
	if *i != "" {
		includes = strings.Split(*i, " ")
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
	if err := gen(workspace, includes, protos); err != nil {
		log.Fatal(err)
	}
}

func prepare(workspace, url string) (string, error) {
	var rc io.ReadCloser
	var name string
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		res, err := http.Get(url)
		if err != nil {
			return "", nil
		}
		rc = res.Body
		name = path.Base(url)
	} else {
		f, err := os.Open(url)
		if err != nil {
			return "", err
		}
		rc = f
		name = filepath.Base(url)
	}
	defer rc.Close()

	file := filepath.Join(workspace, name)
	f, err := os.Create(file)
	if err != nil {
		return "", nil
	}
	if _, err = io.Copy(f, rc); err != nil {
		return "", err
	}
	return f.Name(), f.Close()
}
