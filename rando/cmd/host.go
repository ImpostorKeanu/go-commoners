package main

import (
    "fmt"
    "github.com/impostorkeanu/go-commoners/rando"
    "github.com/spf13/cobra"
)

var (
    hostnameCmdDesc = "Generate a random 'adjective-noun' hostname."
    hostnameCmd     = cobra.Command{
        Use:     "host",
        Aliases: []string{"hn"},
        Short:   hostnameCmdDesc,
        Long:    hostnameCmdDesc,
        Run: func(cmd *cobra.Command, args []string) {
            hn, err := rando.Hostname(nounType)
            if err != nil {
                panic(err)
            }
            fmt.Println(hn)
        },
    }
)

func init() {
    root.AddCommand(&hostnameCmd)
    hostnameCmd.Flags().StringVarP(&nounType, "noun-type", "n",
        "common", nounValueHelp)
}
