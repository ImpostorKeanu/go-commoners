package main

import (
	"github.com/spf13/cobra"
)

var (
	root = cobra.Command{
		Short: "Run rando.",
		Long:  "Run rando to generate random values.",
	}
	nounType       string
	checkProfanity bool
)

func init() {
	root.PersistentFlags().
		StringVarP(&nounType, "noun-type", "n",
			"common", "Type of noun to use. Choices: common, proper")
	root.PersistentFlags().
		BoolVarP(&checkProfanity, "clean", "c",
			false, "Apply a profanity filter.")
}

func main() {
	root.Execute()
}
