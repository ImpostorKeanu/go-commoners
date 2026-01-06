package main

import (
	"fmt"

	go_rando "github.com/impostorkeanu/go-commoners/rando"
	"github.com/spf13/cobra"
)

var uuidCmd = cobra.Command{
	Use:     "uuid",
	Aliases: []string{"u"},
	Short:   "Generate a UUID.",
	Long:    "Generate a UUID.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(go_rando.UUID())
	},
}

func init() {
	root.AddCommand(&uuidCmd)
}
