package main

import (
	"fmt"
	"net"
	"os"

	"github.com/neil-and-void/hurl/src/parser"
)

func main() {
	// // parse command line with flags
	// args := os.Args[1:]
	//
	// // parse file
	// // build request
	// content, err := os.ReadFile(args[0])
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// crlf := []byte{13, 10} // [c]arriage [r]eturn, [l]ine [f]eed
	//
	// formattedContent := []byte{}
	// for _, c := range content {
	// 	if c == 10 {
	// 		formattedContent = append(formattedContent, crlf...)
	// 	} else {
	// 		formattedContent = append(formattedContent, c)
	// 	}
	// }
	// formattedContent = append(formattedContent, crlf...)
	// formattedContent = append(formattedContent, crlf...)
	//
	// // send request
	// conn, err := net.Dial("tcp", "93.184.216.34:80")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// _, err = conn.Write(formattedContent)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// buffer := make([]byte, 4096)
	//
	// // get response
	// _, err = conn.Read(buffer)
	// if err != nil {
	// 	fmt.Println("Error reading:", err.Error())
	// }
	//
	// // output response
	// fmt.Println(string(buffer))

	args := os.Args[1:]
	h, err := parser.ParseHurlFile(args[0])
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v", h)

	addrs, err := net.LookupIP(h.Host)
	if err != nil {
		panic(err)
	}

	ipv4 := string(addrs[0])
	address := fmt.Sprintf("%s:%d", ipv4, 80)

	conn, err := net.Dial("tcp", string(addrs[0]))
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write(h.GetRawRequest())
	if err != nil {
		fmt.Println(err)
		return
	}

	buffer := make([]byte, 4096)

	// get response
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// output response
	fmt.Println(string(buffer))
}
