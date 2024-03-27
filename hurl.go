package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("hello world")
	host := "example.com"

	args := os.Args[1:]

	addrs, _ := net.LookupHost(host)
	fmt.Println(args, addrs, len(addrs))
}
