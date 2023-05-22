package wordlists

import (
    "bufio"
    "bytes"
    "embed"
    "fmt"
    "github.com/impostorkeanu/go-commoners/rando/adjectives"
    "github.com/impostorkeanu/go-commoners/rando/nouns"
    "io/fs"
    "strings"
)

//===================
// EMBED DATA SOURCES
//===================

// Source: https://github.com/janester/mad_libs
//
//go:embed "mad_libs/List of Proper Nouns.txt"
var eProperNouns embed.FS

// Sources:
// - https://github.com/taikuukaits/SimpleWordlists
// - https://github.com/sroberts/wordlists/blob/master/animals.txt
//go:embed SimpleWordlists/Wordlist-Nouns-Common-Audited-Len-3-6.txt sroberts/animals.txt
var eCommonNouns embed.FS

// Source: https://github.com/taikuukaits/SimpleWordlists
//
//go:embed SimpleWordlists/Wordlist-Adjectives-Common-Audited-Len-3-6.txt
var eAdjectives embed.FS

var (
    // AllNouns is a slice containing all known noun values.
    AllNouns *[]string
)

func LoadValues() {
    loadValues(eAdjectives, &adjectives.Values.Common)
    loadValues(eCommonNouns, &nouns.Values.Common)
    loadValues(eProperNouns, &nouns.Values.Proper)
    AllNouns = nouns.All()
}

// loadValues will iterate over each file in embedded and
// read each line into dest.
func loadValues(embedded embed.FS, dest *[]string) {

    // Find all embedded files
    matches, err := fs.Glob(embedded, "**/*.*")
    if err != nil {
        panic(fmt.Sprintf("Failed to read embedded file: %v", err))
    }

    // Read in input files
    for _, m := range matches {

        // Read the file in
        f, _ := embedded.Open(m)
        s, _ := f.Stat()
        b := make([]byte, s.Size())
        f.Read(b)

        // Replace CRLFs with LFs
        bytes.Replace(b, []byte("\x0d\x0a"), []byte("\n"), -1)

        // Pass the content to a reader and scan it for newline
        // values.
        reader := bytes.NewReader(b)
        scanner := bufio.NewScanner(reader)
        for scanner.Scan() {
            *dest = append(*dest, strings.ToLower(scanner.Text()))
        }
        if err = scanner.Err(); err != nil {
            fmt.Sprintf("failed to parse wordlist: %v")
            panic(err)
        }
    }
}
