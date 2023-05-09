package net

import (
    "net"
    "testing"
)

func TestFindInterface(t *testing.T) {
    ifaces, err := net.Interfaces()
    if err != nil {
        t.Fatalf("failed to list interfaces: %v", err)
    }
    checked := 0
    for _, i := range ifaces {
        var a []net.Addr
        if a, _ = i.Addrs(); len(a) == 0 {
            t.Logf("skipping interface with no addresses: %v", i.Name)
            continue
        }
        var ip string
        if ip, err = FindInterface(i.Name); err != nil {
            t.Fatalf("failed to find interface by name: %v", err)
            checked++
        } else {
            t.Logf("found %s (%s) interface by name", ip, i.Name)
            checked++
        }
        if _, err = FindInterface(ip); err != nil {
            t.Fatalf("failed to find interface by ip: %v", ip)
            checked++
        } else {
            t.Logf("found interface by ip: %s", ip)
        }
    }
    if checked == 0 {
        t.Fatalf("no valid network interfaces found to check")
    }
}
