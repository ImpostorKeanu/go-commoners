package main

import (
    "fmt"
    "github.com/ImpostorKeanu/go-commoners/rando"
    "github.com/ImpostorKeanu/go-commoners/rando/nouns"
    "github.com/spf13/cobra"
    "os"
)

var (
    dnsnameCmdDesc        = "Generate a random 'adjective-noun' DNS name."
    dnsNameResolve        bool
    dnsNameApex           string
    dnsNameUntilNoResolve bool
    dnsNameCmd            = cobra.Command{
        Use:     "dns",
        Aliases: []string{"dn"},
        Short:   dnsnameCmdDesc,
        Long:    dnsnameCmdDesc,
        Run:     dnsName}
)

func init() {
    dnsNameCmd.Flags().StringVarP(&nounType,
        "noun-type", "n", "common", nounValueHelp)
    dnsNameCmd.Flags().BoolVarP(&dnsNameResolve,
        "resolve", "r", true,
        "Determines if name resolution should occur.")
    dnsNameCmd.Flags().StringVarP(&dnsNameApex,
        "apex-domain", "a", "autonode.net",
        "Apex domain that will be used for name resolution.")
    dnsNameCmd.Flags().BoolVarP(&dnsNameUntilNoResolve,
        "until-unique", "u", true,
        "Continue generate values until name resolution fails.")
    root.AddCommand(&dnsNameCmd)
}

func dnsName(cmd *cobra.Command, args []string) {

    nT := nouns.NounType(nounType)

    if !nouns.ValidType(nT) {
        fmt.Printf("Invalid noun type: %v\n", nounType)
        fmt.Println("Valid types: proper, common")
        os.Exit(1)
    }

    for {

        //==============================
        // GENERATE AND RESOLVE THE NAME
        //==============================

        name, resolved, _, err := rando.DnsName(nT, dnsNameApex, dnsNameResolve)

        if err != nil {

            ERR.Fatalf("Failed to generate DNS name: %v\n", err)

        } else if !dnsNameResolve {

            //========================
            // DISABLED DNS RESOLUTION
            //========================

            WARN.Println("DNS name resolution was disabled")
            fmt.Println(name)
            break

        } else if resolved && !dnsNameUntilNoResolve {

            //============================
            // RESOLVED BUT NO "UNTIL NEW"
            //============================

            WARN.Fatalf("DNS name DID resolved: %s\n", name)

        } else if resolved {

            //=========================
            // RESOLVED BUT KEEP TRYING
            //=========================

            INFO.Printf("DNS name resolved; retrying: %s\n", name)
            continue

        } else {

            //======================
            // UNIQUE NAME GENERATED
            //======================

            INFO.Printf("Generated unique DNS name: %s\n", name)
            fmt.Println(name)
            break
        }

    }

}
