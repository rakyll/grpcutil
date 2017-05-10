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

func gen(protos []string) error {
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
		return errors.New("language plugin is not available")
	}

	args := append(protos, "--go_out=plugins=grpc:.")
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
