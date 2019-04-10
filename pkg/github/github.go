package github

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func SplitModule(mod string) (string, string, error) {
	spl := strings.Split(mod, "/")
	if len(spl) < 3 {
		return "", "", errors.WithStack(fmt.Errorf("module malformed"))
	}
	return spl[1], spl[2], nil
}
