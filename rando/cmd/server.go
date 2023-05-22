package main

import (
    "fmt"
    "github.com/impostorkeanu/go-commoners/net"
    "github.com/impostorkeanu/go-commoners/rando/server"
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
            var a, p string
            s := strings.Split(serverAddress, ":")
            if len(s) < 2 {
                ERR.Fatalf("Socket input is expected")
            }
            a = s[0]
            p = s[1]
            if a, err := net.FindInterface(a); err != nil {
                ERR.Fatalf(err.Error())
            } else {
                server.RandoServer(fmt.Sprintf("%s:%s", a, p)).Start()
            }
        },
    }
)

func init() {
    root.AddCommand(&serverCmd)
    serverCmd.Flags().StringVarP(&serverAddress, "address", "a",
        serverAddress, "In socket format, the address the server will listen on.")
    serverCmd.MarkFlagRequired("address")
}
