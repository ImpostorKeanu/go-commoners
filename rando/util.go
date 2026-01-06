package rando

import (
	"errors"
	"strings"

	goaway "github.com/TwiN/go-away"
	"golang.org/x/exp/rand"
)

// UntilCleanString runs the get function and checks the output
// against profanity filters.
//
// Setting check to false bypasses the profanity filter.
func UntilCleanString(check bool, get func() (string, error)) (v string, err error) {
	v, err = get()
	for err == nil && check && goaway.IsProfane(v) {
		v, err = get()
	}
	return
}

// UntilCleanSlice converts the value returned by the get function to
// a string and checks it against profanity filters.
//
// Setting check to false bypasses the profanity filter.
func UntilCleanSlice(check bool, get func() ([]string, error)) (v []string, err error) {
	v, err = get()
	for err == nil && check && goaway.IsProfane(strings.Join(v, "")) {
		v, err = get()
	}
	return
}

func gtz(v uint) error {
	if v == 0 {
		return errors.New("max must be > 0")
	}
	return nil
}

// UintGTZ returns an uint that is greater-than-zero.
func UintGTZ(max uint) (out uint, err error) {
	err = gtz(max)
	for err == nil && out == 0 {
		out = uint(rand.Intn(int(max)))
	}
	return out, err
}
