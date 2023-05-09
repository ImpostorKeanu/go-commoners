package adjectives

import (
    "github.com/arch4ngel/go-commoners/rando/shared"
)

var (
    Values = values{}
)

// adjectives is a struct pointing to common adjective values.
type values struct {
    Common []string
}

// Get returns a pointer to a random adjective string.
func Get() *string {
    return shared.RandVal(&Values.Common)
}
