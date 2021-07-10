package main

import (
	"fmt"
	"mini/lib"
)

func main() {
	fmt.Println("mini")
	ips := lib.GatherValidIPs()
	for _, ip := range ips {
		lib.Connect(ip)
	}
}
