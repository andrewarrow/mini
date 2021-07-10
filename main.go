package main

import (
	"fmt"
	"mini/lib"
)

func main() {
	fmt.Println("mini")
	for _, seed := range lib.DNSSeeds {
		lib.IPsForHost(seed)
	}
	lib.Connect("")
}
