package rando

import (
	"regexp"
	"testing"
)

func TestAnyAsciiString(t *testing.T) {
	t.Logf("%v", AnyASCIIString(uint32(10), true, ""))
	if matched, _ := regexp.MatchString("^[0-9]?.+[0-9]", AnyASCIIString(uint32(100), false, "")); matched {
		t.Fatal("Number returned in string!\n")
	}
}
