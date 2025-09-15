package storage

import (
	"testing"

	"github.com/pkg/errors"
)

func TestError(t *testing.T) {
	err := &NoVideoFileError{Path: "/some/path"}
	if errors.Is(err, &NoVideoFileError{}) {
		t.Log("is NoVideoFileError")
	} else {
		t.Error("not match")
	}
}
