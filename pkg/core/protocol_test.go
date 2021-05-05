package core

import (
	"strings"
	"testing"
)

func TestGetSize(t *testing.T) {
	res, err := getSize(strings.NewReader("*1\r\n"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if res != 1 {
		t.Errorf("should be 1 not %d", res)
	}
}
