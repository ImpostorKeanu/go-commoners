package net

import (
    "errors"
    "fmt"
    "net"
)

// FindInterface iterates over all network interfaces and attempts to find one
// that matches either the interface's name or IP address.
//
// An error is returned when an interface is requested by name, but does not
// have an IP address associated with it.
//
// # Notes
//
// - This function is mostly useful in simple environments, such as Docker
//   containers.
// - It's overall naive and does not attempt to anticipate all possible
//   deployment scenarios.
func FindInterface(target string) (match string, err error) {

    var ifaces []net.Interface
    if ifaces, err = net.Interfaces(); err != nil {
        return match, err
    } else if len(ifaces) == 0 {
        return match, errors.New("no network interfaces detected")
    }

    //==============================
    // SEARCH ALL NETWORK INTERFACES
    //==============================

    var addrs []net.Addr
    for _, iface := range ifaces {

        if addrs, err = iface.Addrs(); err != nil {
            return target, err
        } else if len(addrs) < 1 && iface.Name == target {
            return target, errors.New("failed to get address from interface name")
        } else if iface.Name == target {

            //====================================
            // PULL IP ADDRESS FROM INTERFACE NAME
            //====================================

            match = addrs[0].(*net.IPNet).IP.String()
            break

        } else {

            //===============================
            // SEARCH FOR MATCHING IP ADDRESS
            //===============================

            for _, iA := range addrs {

                if target == iA.(*net.IPNet).IP.String() {
                    match = target
                    break
                }
            }
        }
    }

    if match == "" {
        err = errors.New(fmt.Sprintf("failed to find requested interface %s", match))
    }

    return match, err
}
