package rando

import (
    "errors"
    "fmt"
    "github.com/arch4ngel/go-commoners/rando/adjectives"
    "github.com/arch4ngel/go-commoners/rando/nouns"
    "github.com/arch4ngel/go-commoners/rando/shared"
    "github.com/arch4ngel/go-commoners/rando/wordlists"
    "github.com/google/uuid"
    "golang.org/x/exp/rand"
    "net"
    "strings"
    "time"
)

var (
    asciiUpper, asciiLower, asciiNumber []string
    AsciiUpper                          []string
    AsciiLower                          []string
    AsciiNumber                         []string
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
    //var dest, expDest []string
    var dest, expDest *[]string
    for destOff, tup := range [3][2]uint8{{0x30, 0x39}, {0x41, 0x5a}, {0x61, 0x7a}} {

        // Derive the destination array
        switch destOff {
        case 0:
            dest = &asciiNumber
            expDest = &AsciiNumber
        case 1:
            dest = &asciiUpper
            expDest = &AsciiUpper
        case 2:
            dest = &asciiLower
            expDest = &AsciiLower
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

    if err == nil {
        var n *string
        if n, err = nouns.Get(nT); err == nil {
            v = fmt.Sprintf("%v-%v", *adjectives.Get(), *n)
        }
    }

    return v, err

}

// Uuid returns a Uuid string.
func Uuid() string {
    return uuid.New().String()
}

// DnsNameMustNotResolve perpetually generates FQDN values and attempts
// name resolution until one fails, resulting in a unique FQDN that is
// not associated with an active DNS A record.
func DnsNameMustNotResolve[S shared.StringType](nT S, dnsApex string) (name string, err error) {
    resolved := true
    for resolved && err == nil {
        name, resolved, _, err = DnsName(nT, dnsApex, true)
    }
    return name, err
}

// DnsName returns a randomized DNS FQDN.
//
// When doResolve is true, DNS resolution will occur.
//
// nT indicates the type of noun to use, i.e., common or proper.
func DnsName[S shared.StringType](nT S, dnsApex string, doResolve bool) (name string, resolved bool, ips []net.IP, err error) {

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

type MaxDnsTries string

func (m MaxDnsTries) Error() string {
    return string(m)
}

func DnsUntil[S shared.StringType](nT S, dnsNameApex string, dnsNameUntilNoResolve bool) (fqdn string, resolved bool, ips []net.IP, err error) {

    tries := 20
    for ; tries > 0; tries-- {

        //==============================
        // GENERATE AND RESOLVE THE NAME
        //==============================

        fqdn, resolved, ips, err = DnsName(nT, dnsNameApex, true)

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
        err = MaxDnsTries("max dns tries of 20 reached")
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

// AnyAsciiSlice will generate any alphanumeric value, including upper case,
// in a slice. Suppress numbers by setting allowNumbers to false.
func AnyAsciiSlice(minLen uint32, allowNumbers bool) (out []string) {
    var src []string
    out = []string{}
    for len(out) < int(minLen) {
        switch rand.Intn(3) {
        case 0:
            src = AsciiUpper
        case 1:
            src = AsciiLower
        case 2:
            if !allowNumbers {
                continue
            }
            src = AsciiNumber
        }
        out = append(out, src[rand.Intn(len(src))])
    }
    return out
}

// AnyAsciiString uses AnyAsciiSlice to generate a string value that
// was joined on del.
func AnyAsciiString(minLen uint32, allowNumbers bool, del string) string {
    return strings.Join(AnyAsciiSlice(minLen, allowNumbers), del)
}

// AnyString returns a string of random values joined by
// a delimiter.
//
// minLen determines the maximum length of the values joined on an
// empty string.
func AnyString(minLen uint32, del string) string {
    return strings.Join(AnySlice(minLen), del)
}

// AnySlice returns a slice of random values.
//
// minLen determines the maximum length of the values joined on an
// empty string.
func AnySlice(minLen uint32) (out []string) {

    // buffer to hold a pointer to the current
    // string
    var bv *string

    //=========================
    // CONSTRUCT THE PASSPHRASE
    //=========================

    for len(strings.Join(out, "")) < int(minLen) {

        if shared.Rnd.Intn(2) == 1 {

            // Get a noun of random type
            bv, _ = nouns.Get(NounType())

        } else {

            // Get an adjective
            bv = adjectives.Get()

        }

        out = append(out, *bv)

    }

    return out
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
    for ; err == nil && out == 0; {
        out = uint(rand.Intn(int(max)))
    }
    return out, err
}

// AnyAsciiRandfix concatenates a random ASCII string of 1 to min length to the
// beginning and end of in. This aims to enhance randomization of values.
func AnyAsciiRandfix(in, del string, min uint) (_ string, err error) {
    if err = gtz(min); err == nil {
        r, _ := UintGTZ(min)
        in = AnyAsciiString(uint32(r), true, "") + del + in
        r, _ = UintGTZ(min)
        in += del + AnyAsciiString(uint32(r), true, "")
    }
    return in, err
}
