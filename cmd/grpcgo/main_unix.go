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

// +build darwin linux

package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func install(force bool) error {
	_, ok := whichProtoc()
	if ok || !force {
		return nil
	}
	// Install from the source.
	return nil
}

func gen(workspace string, includes []string, protos []string) error {
	if err := install(false); err != nil {
		log.Fatalf("Cannot install protoc: %v", err)
	}
	_, ok := whichProtocGen()
	if !ok {
		// Install langauge plugin.
		if err := installProtocGen(); err != nil {
			return err
		}
	}
	_, ok = whichProtocGen()
	if !ok {
		return errors.New("language plugin is not available; make sure GOPATH/bin is in your PATH")
	}

	args := []string{"--proto_path"}
	args = append(args, workspace)
	args = append(args, "-I")
	args = append(args, includes...)
	args = append(args, workspace)
	args = append(args, protos...)
	args = append(args, "--go_out=plugins=grpc:"+*o)

	cmd := exec.Command("protoc", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot generate: %s", out)
	}
	return nil
}

func installProtocGen() error {
	// TODO(jbd): Support other languages.
	cmd := exec.Command("go", "get", "-u", "github.com/golang/protobuf/protoc-gen-go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot go get protoc plugin: %s", out)
	}
	return nil
}

func whichProtocGen() (protocGen string, ok bool) {
	// TODO(jbd): Support other languages.
	return whichBinary("protoc-gen-go")
}

func whichProtoc() (protoc string, ok bool) {
	return whichBinary("protoc")
}

func whichBinary(name string) (path string, ok bool) {
	cmd := exec.Command("which", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(out)), len(out) > 0
}
