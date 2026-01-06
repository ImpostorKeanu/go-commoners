package main

import (
	"fmt"
	"strings"

	"github.com/impostorkeanu/go-commoners/rando"
	"github.com/spf13/cobra"
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
			pp, err := rando.UntilCleanString(checkProfanity, func() (v string, err error) {
				v = rando.AnyString(ppLen, ppDel)
				if ppRandFix {
					v, err = rando.AnyASCIIRandfix(v, ppDel, 15)
				}
				return
			})
			if err != nil {
				ERR.Fatalf("error while generating passphrase: %v", err)
			}
			INFO.Printf("value: \"%v\"", pp)
			INFO.Printf("delimiter: \"%v\"", ppDel)
			INFO.Printf("length:\n  - requested: %v\n  - final: %v", ppLen, len(pp))
			INFO.Printf("word count: %v", len(strings.Split(pp, ppDel)))
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
