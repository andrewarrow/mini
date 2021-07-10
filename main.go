package main

import (
	"fmt"
	"mini/lib"
)

func main() {
	fmt.Println("mini")
	lib.GatherValidIPs()
	lib.Connect("")
}
