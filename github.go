package main

import (
	"fmt"
	"strings"
)

func splitGHModule(mod string) (string, string, error) {
	spl := strings.Split(mod, "/")
	if len(spl) < 3 {
		return "", "", fmt.Errorf("module malformed")
	}
	return spl[1], spl[2], nil
}
