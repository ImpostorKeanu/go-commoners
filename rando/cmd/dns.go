package main

import (
	"fmt"
	"os"

	goaway "github.com/TwiN/go-away"
	"github.com/impostorkeanu/go-commoners/rando"
	"github.com/impostorkeanu/go-commoners/rando/nouns"
	"github.com/spf13/cobra"
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
		Run:     dnsName,
	}
)

func init() {
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
		fmt.Printf("nvalid noun type: %v", nounType)
		fmt.Println("Valid types: proper, common")
		os.Exit(1)
	}

outer:
	for {

		//==============================
		// GENERATE AND RESOLVE THE NAME
		//==============================

		name, resolved, _, err := rando.DNSName(nT, dnsNameApex, dnsNameResolve)

		switch {
		case err != nil:
			ERR.Fatalf("failed to generate dns name: %v", err)
		case checkProfanity && goaway.IsProfane(name):
			continue
		case !dnsNameResolve:
			WARN.Println("dns name resolution was disabled")
			fmt.Println(name)
			break outer
		case resolved && !dnsNameUntilNoResolve:
			WARN.Fatalf("known dns name resolved: %s", name)
		case resolved:
			continue
		}

		INFO.Printf("generated unique dns name: %s", name)
		fmt.Println(name)
		break

	}
}
