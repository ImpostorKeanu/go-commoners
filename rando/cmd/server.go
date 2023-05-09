package main

import (
    "github.com/arch4ngel/go-commoners/net"
    "github.com/arch4ngel/go-commoners/rando/server"
    "github.com/spf13/cobra"
    "strings"
)

var (
    serverAddress string
    serverCmdDesc = "Start an API server."
    serverCmd     = cobra.Command{
        Use:     "server",
        Aliases: []string{"s"},
        Short:   serverCmdDesc,
        Long:    serverCmdDesc,
        Run: func(cmd *cobra.Command, args []string) {
            if _, err := net.FindInterface(strings.Split(serverAddress, ":")[0]); err != nil {
                ERR.Fatalf(err.Error())
            }
            server.RandoServer(serverAddress).Start()
        },
    }
)

func init() {
    root.AddCommand(&serverCmd)
    serverCmd.Flags().StringVarP(&serverAddress, "address", "a",
        serverAddress, "In socket format, the address the server will listen on.")
    serverCmd.MarkFlagRequired("address")
}
