package shared

import (
    "math/rand"
    "time"
)

var (
    Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandVal returns a pointer a random value from s.
func RandVal(s *[]string) *string {
    return &(*s)[Rnd.Intn(len(*s))]
}
