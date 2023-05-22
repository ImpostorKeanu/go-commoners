package main

import (
    "fmt"
    "github.com/impostorkeanu/go-commoners/rando"
    "github.com/spf13/cobra"
)

var (
    asciiCmdDesc  = "Generate a random ascii string."
    asciiLen      uint32
    asciiDel      string
    asciiAllownum bool
    asciiCmd      = cobra.Command{
        Use:     "ascii",
        Aliases: []string{"a"},
        Short:   passphraseCmdDesc,
        Long:    passphraseCmdDesc,
        Run: func(cmd *cobra.Command, args []string) {
            a := rando.AnyAsciiString(asciiLen, asciiAllownum, asciiDel)
            fmt.Println(a)
        },
    }
)

func init() {
    asciiCmd.Flags().Uint32VarP(&asciiLen, "min-length",
        "m", 25,
        "Minimum value length.")
    asciiCmd.Flags().StringVarP(&asciiDel, "delimiter",
        "d", "",
        "Value delimiter.")
    asciiCmd.Flags().BoolVarP(&asciiAllownum, "allow-numbers", "a", true,
        "Determines if numbers will be included.")
    root.AddCommand(&asciiCmd)
}
