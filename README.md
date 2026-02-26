# go-commoners

A collection of common functions and objects used to support my automation efforts.

# Rando Template Functions

Use `rando.TemplateFuncs` to get a function map for Go templates.

Outputs contain only numbers and letters.

## Usage

See `rando/util_test.go` for an example.

## Function Signatures

```go
// randoAscii returns a random ascii value.
func randoAscii(minLen uint32, allowNumbers bool, delimiter string) string {}

// randoAnyString returns a series of adjective/noun
// values.
func randoAnyString(minLen uint32, delimiter string) string {}

// randoAnyString returns an adjective-noun hostname.
func randoAnyHostname(nounType string) (string, error) {}

// randoUniqueDNSName generates adjective-noun records until
// one that _does not resolve_ is created.
//
// Note that this function performs a DNS query for each
// value generated.
func randoUniqueDNSName(nounType string, dnsApex string) (string, error) {}

// randoIsProfane is used to see if a value contains potential
// profanity.
func randoIsProfane(v string) bool {}
```
