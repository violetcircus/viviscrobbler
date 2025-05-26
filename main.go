package main

import (
	"fmt"
	"bufio"
	"log"
	"net"
	"strings"
	"github.com/violetcircus/viviscrobbler/internal"
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
	internal.Setup()
	fmt.Println("viviscrobbler!")
	conn, err := net.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	track_info := GetSong(conn)
	fmt.Println("Artist:", internal.GetArtist(track_info["Artist"]))
	var config = internal.ReadConfig()
	fmt.Println(config.SingleArtist)
}

func GetSong (conn net.Conn) map[string]string {
	// create bufio reader that reads what comes from the connection
	reader := bufio.NewReader(conn)

	// read each new line
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// print the output
	fmt.Println("Server:", line)

	fmt.Fprintf(conn, "currentsong\n")

	// create track info map
	track_info := make(map[string]string)
	// loop over lines returned by http request
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		// check if response finished, stop if it has
		if line == "OK" || strings.HasPrefix(line, "ACK") {
			fmt.Println("Response:", line)
			break
		}
		// output results to terminal
		fmt.Println(line)
		// put results in the map
		key, value, found := strings.Cut(line, ":")
		if found {
			track_info[key] = value
		}
	} 
	return track_info
}
