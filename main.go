package main

import (
	"fmt"
	"bufio"
	"log"
	"net"
	"strings"
)

type scrobble struct {
	track_info string
	status string
	timestamp string
	api_key string
	sk string
	api_secret string
	method string
}

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

	track_info := make(map[string]string)
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
		fmt.Println(line)
		key, value, found := strings.Cut(line, ":")
		if found {
			track_info[key] = value
		}
	} 
	fmt.Println(track_info["Last-Modified"])
}
