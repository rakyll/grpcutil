// +build !darwin

package main

import "errors"

func install() error {
	return errors.New("not supported on this platform")
}

func gen(lang string) error {
	return errors.New("not supported on this platform")
}
