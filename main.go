package main

import (
	"fmt"
	"net"
	"time"

	"github.com/andrewarrow/mini/lib"
)

func main() {
	fmt.Println("mini")

	go func() {
		for mp := range lib.MiniPostChan {
			fmt.Println(mp.Body)
			fmt.Println(mp.ImageURLs)
			fmt.Println("")
			fmt.Println(time.Unix(mp.Timestamp, 0))
			fmt.Println("")
			fmt.Println("https://bitclout.com/posts/" + mp.PostHashHex)
			fmt.Println("Poster Public Key", mp.PosterPub58)
			fmt.Println("")
			fmt.Println("")
		}
	}()

	lib.Connect("peer1", net.ParseIP("35.232.92.5"))
	//lib.Connect("peer2", net.ParseIP("46.4.89.216"))
	//lib.Connect("peer3", net.ParseIP("78.46.99.243"))
}

func main2() {
	for _, ip := range ips {
		lib.TestConnect(ip, net.ParseIP(ip))
	}
}

var ips = []string{
	"46.4.89.216",
	"78.46.99.243",
	"107.161.127.130",
	"54.36.166.13",
	"144.217.207.70",
	"45.77.52.104",
	"52.59.0.100"}
