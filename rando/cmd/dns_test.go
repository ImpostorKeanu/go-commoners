package main

import (
    "testing"
)

func init() {
    nounType = "common"
    dnsNameResolve = true
    dnsNameApex = "autonode.net"
    dnsNameUntilNoResolve = true
}

func TestDnsName(t *testing.T) {
    args := []string{}
    dnsName(&dnsNameCmd, args)
}
