package rando

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/impostorkeanu/go-commoners/rando/adjectives"
	"github.com/impostorkeanu/go-commoners/rando/nouns"
	"github.com/impostorkeanu/go-commoners/rando/shared"
	"github.com/impostorkeanu/go-commoners/rando/wordlists"
	"golang.org/x/exp/rand"
)

var (
	asciiUpper, asciiLower, asciiNumber []string
	ASCIIUpper                          []string
	ASCIILower                          []string
	ASCIINumber                         []string
)

func init() {
	wordlists.LoadValues()

	// Seed the random number generator
	rand.Seed(uint64(time.Now().UnixMicro()))

	//==================================
	// INITIALIZE ASCII CHARACTER RANGES
	//==================================

	// ASCII ranges to consider
	// 0-9 > 0x30-0x39
	// A-Z > 0x41-0x5a
	// a-z > 0x61-0x7a

	// dest will be the target tuple to read from
	// var dest, expDest []string
	var dest, expDest *[]string
	for destOff, tup := range [3][2]uint8{{0x30, 0x39}, {0x41, 0x5a}, {0x61, 0x7a}} {

		// Derive the destination array
		switch destOff {
		case 0:
			dest = &asciiNumber
			expDest = &ASCIINumber
		case 1:
			dest = &asciiUpper
			expDest = &ASCIIUpper
		case 2:
			dest = &asciiLower
			expDest = &ASCIILower
		}

		for i := tup[0]; i <= tup[1]; i++ {
			*dest = append(*dest, string(i))
			*expDest = append(*expDest, string(i))
		}
	}
}

// Hostname returns a randomized "adjective-noun" hostname.
//
// nT indicates the type of noun to use, i.e., common or proper.
func Hostname[S shared.StringType](nT S) (v string, err error) {
	return UntilCleanString(true, func() (string, error) {
		var n *string
		if n, err = nouns.Get(nT); err == nil {
			v = fmt.Sprintf("%v-%v", *adjectives.Get(), *n)
		}
		return v, err
	})
}

// UUID returns a UUID string.
func UUID() string {
	return uuid.New().String()
}

// DNSNameMustNotResolve perpetually generates FQDN values and attempts
// name resolution until one fails, resulting in a unique FQDN that is
// not associated with an active DNS A record.
func DNSNameMustNotResolve[S shared.StringType](nT S, dnsApex string) (name string, err error) {
	resolved := true
	for resolved && err == nil {
		name, resolved, _, err = DNSName(nT, dnsApex, true)
	}
	return name, err
}

// DNSName returns a randomized DNS FQDN.
//
// When doResolve is true, DNS resolution will occur.
//
// nT indicates the type of noun to use, i.e., common or proper.
func DNSName[S shared.StringType](nT S, dnsApex string, doResolve bool) (name string, resolved bool, ips []net.IP, err error) {
	//==================
	// GENERATE THE FQDN
	//==================

	hn, _ := Hostname(nT)
	name = fmt.Sprintf("%v.%v", hn, dnsApex)

	if doResolve {
		//======================================
		// ATTEMPT TO RESOLVE THE GENERATED FQDN
		//======================================

		if ips, err = net.LookupIP(name); err == nil && len(ips) > 0 {
			//===================
			// COLLISION OCCURRED
			//===================

			resolved = true
		} else if err != nil {
			if e, ok := err.(*net.DNSError); ok && e.Err == "no such host" {
				//=================
				// RESET ERR TO NIL
				//=================
				// Failure is desired because we want a unique DNS record!

				err = nil
			}
		}
	}

	if ips == nil {
		ips = []net.IP{}
	}

	return name, resolved, ips, err
}

type MaxDNSTries string

func (m MaxDNSTries) Error() string {
	return string(m)
}

func DNSUntil[S shared.StringType](nT S, dnsNameApex string, dnsNameUntilNoResolve bool) (fqdn string, resolved bool, ips []net.IP, err error) {
	tries := 20
	for ; tries > 0; tries-- {

		//==============================
		// GENERATE AND RESOLVE THE NAME
		//==============================

		fqdn, resolved, ips, err = DNSName(nT, dnsNameApex, true)

		if err != nil || resolved && !dnsNameUntilNoResolve {
			break
		} else if resolved {
			//=========================
			// RESOLVED BUT KEEP TRYING
			//=========================

			continue
		} else {
			//======================
			// UNIQUE NAME GENERATED
			//======================

			break
		}
	}

	if tries == 0 {
		err = MaxDNSTries("max dns tries of 20 reached")
	}

	return fqdn, resolved, ips, err
}

// NounType returns one of the two wordlists.NounType
// values.
func NounType() nouns.NounType {
	if shared.Rnd.Intn(2) == 1 {
		return nouns.ProperNounType
	} else {
		return nouns.CommonNounType
	}
}

// AnyASCIISlice will generate any alphanumeric value, including upper case,
// in a slice. Suppress numbers by setting allowNumbers to false.
func AnyASCIISlice(minLen uint32, allowNumbers bool) []string {
	v, _ := UntilCleanSlice(true, func() ([]string, error) {
		var src []string
		out := []string{}
		for len(out) < int(minLen) {
			switch rand.Intn(3) {
			case 0:
				src = ASCIIUpper
			case 1:
				src = ASCIILower
			case 2:
				if !allowNumbers {
					continue
				}
				src = ASCIINumber
			}
			out = append(out, src[rand.Intn(len(src))])
		}
		return out, nil
	})
	return v
}

// AnyASCIIString uses AnyAsciiSlice to generate a string value that
// was joined on del.
func AnyASCIIString(minLen uint32, allowNumbers bool, del string) string {
	v, _ := UntilCleanString(true, func() (string, error) {
		return strings.Join(AnyASCIISlice(minLen, allowNumbers), del), nil
	})
	return v
}

// AnyString returns a string of random values joined by
// a delimiter.
//
// minLen determines the maximum length of the values joined on an
// empty string.
func AnyString(minLen uint32, del string) string {
	v, _ := UntilCleanString(true, func() (string, error) {
		return strings.Join(AnySlice(minLen), del), nil
	})
	return v
}

// AnySlice returns a slice of random values.
//
// minLen determines the maximum length of the values joined on an
// empty string.
func AnySlice(minLen uint32) (out []string) {
	v, _ := UntilCleanSlice(true, func() ([]string, error) {
		var bv *string
		for len(strings.Join(out, "")) < int(minLen) {
			if shared.Rnd.Intn(2) == 1 {
				// get a noun
				bv, _ = nouns.Get(NounType())
			} else {
				// get an adjective
				bv = adjectives.Get()
			}
			out = append(out, *bv)
		}
		return out, nil
	})
	return v
}

// AnyASCIIRandfix concatenates a random ASCII string of 1 to min length to the
// beginning and end of in. This aims to enhance randomization of values.
func AnyASCIIRandfix(in, del string, min uint) (_ string, err error) {
	if err = gtz(min); err == nil {
		r, _ := UintGTZ(min)
		in = AnyASCIIString(uint32(r), true, "") + del + in
		r, _ = UintGTZ(min)
		in += del + AnyASCIIString(uint32(r), true, "")
	}
	return in, err
}
