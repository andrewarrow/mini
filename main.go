package main

import (
	"fmt"
	"mini/lib"
	"net"
	"time"
)

func main() {
	fmt.Println("mini")
	/*
		ips := lib.GatherValidIPs()
		for _, ip := range ips {
			go lib.Connect(ip.String(), ip)
		}*/
	//104.238.183.241
	//lib.Connect(net.ParseIP("64.98.145.30"))
	//lib.Connect(net.ParseIP("104.238.183.241"))
	lib.Connect("35.232.92.5", net.ParseIP("35.232.92.5"))
	for {
		time.Sleep(time.Second * 1)
	}
}
