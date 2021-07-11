package main

import (
	"fmt"
	"mini/lib"
	"net"
)

func main() {
	fmt.Println("mini")
	//ips := lib.GatherValidIPs()
	//for _, ip := range ips {
	//	lib.Connect(ip)
	//}
	//104.238.183.241
	//35.232.92.5
	lib.Connect(net.ParseIP("64.98.145.30"))
}
