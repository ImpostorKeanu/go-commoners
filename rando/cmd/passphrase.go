package main

import (
    "fmt"
    "github.com/impostorkeanu/go-commoners/rando"
    "github.com/spf13/cobra"
    "strings"
)

var (
    passphraseCmdDesc = "Generate a random 'adjective-noun' passphrase."
    ppLen             uint32
    ppDel             string
    ppRandFix         bool
    passphraseCmd     = cobra.Command{
        Use:     "passphrase",
        Aliases: []string{"pp"},
        Short:   passphraseCmdDesc,
        Long:    passphraseCmdDesc,
        Run: func(cmd *cobra.Command, args []string) {
            pp := rando.AnyString(ppLen, ppDel)
            if ppRandFix {
                pp, _ = rando.AnyAsciiRandfix(pp, ppDel, 15)
            }
            INFO.Printf("Value: \"%v\"", pp)
            INFO.Printf("Delimiter: \"%v\"", ppDel)
            INFO.Printf("Length:\n  - Requested: %v\n  - Final: %v", ppLen, len(pp))
            INFO.Printf("Word Count: %v", len(strings.Split(pp, ppDel)))
            fmt.Println(pp)
        },
    }
)

func init() {
    passphraseCmd.Flags().Uint32VarP(&ppLen, "min-length",
        "m", 25,
        "Minimum passphrase length.")
    passphraseCmd.Flags().StringVarP(&ppDel, "delimiter",
        "d", " ",
        "Word delimiter.")
    passphraseCmd.Flags().BoolVarP(&ppRandFix, "random-fix",
        "r", false, "For enhanced randomization, prefix and suffix the final output with up to 15 ASCII characters")
    root.AddCommand(&passphraseCmd)
}
