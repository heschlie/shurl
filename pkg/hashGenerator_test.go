package shurl

import (
	"regexp"
	"testing"
)

func TestRandStringRunes(t *testing.T) {
	r := regexp.MustCompile("^[a-zA-Z]{8}$")
	got := RandStringRunes(8)
	if !r.MatchString(got) {
		t.Errorf("Result of RandRuneStringRunes did not match regex: %s", got)
	}
}
