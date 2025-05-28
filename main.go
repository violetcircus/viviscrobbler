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
	reader := bufio.NewReader(conn)

	for {
		// tell mpd we're idling
		fmt.Fprintln(conn, "idle player") 
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Println("yurp we died")
				log.Fatal(err)
			}
			if strings.HasPrefix(line, "changed: player") {
				break
			}
		}

		trackInfo := GetSong(conn, reader)

		log.Println(internal.GetArtist(trackInfo))
	}
}


func GetSong (conn net.Conn, reader *bufio.Reader) map[string]string {
	// ask for current song
	fmt.Fprintf(conn, "currentsong\n")

	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// print the status
	fmt.Println("Server:", line)

	// create track info map
	trackInfo := make(map[string]string)
	// loop over song info
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
