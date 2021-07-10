package lib

import (
	"fmt"
	"net"
)

func Connect(ip net.IP) {
	fmt.Println("connecting to peer", ip)
}
