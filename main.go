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
	trackInfo string
	status string
	timestamp string
	apiKey string
	sk string
	apiSecret string
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
	trackInfo := GetSong(conn)

	fmt.Println("Artist:", internal.GetArtist(trackInfo))

	// redo this to use yaml
	// var config = internal.ReadConfig()
	// fmt.Println(config.SingleArtist)
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
	trackInfo := make(map[string]string)
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
			trackInfo[key] = value
		}
	} 
	return trackInfo
}
