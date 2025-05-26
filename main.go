package main

import (
	"fmt"
	"bufio"
	"log"
	"net"
	"strings"
)

func main() {
	fmt.Println("viviscrobbler!")
	conn, err := net.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// create bufio reader that reads what comes from the connection
	reader := bufio.NewReader(conn)

	// read each new line
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	//print the output
	fmt.Println("Server:", line)

	fmt.Fprintf(conn, "currentsong\n")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		if line == "OK" || strings.HasPrefix(line, "ACK") {
			fmt.Println("Response:", line)
			break
		}
		fmt.Println("status:", line)
	} 
}
