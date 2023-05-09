package main

import (
    "github.com/spf13/cobra"
)

var (
    root = cobra.Command{
        Short: "Run rando.",
        Long:  "Run rando to generate random values.",
    }
    nounType      string
    nounValueHelp = "Type of noun to use. Choices: common, proper"
)

func main() {
    root.Execute()
}
