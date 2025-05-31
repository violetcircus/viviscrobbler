package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/violetcircus/viviscrobbler/internal/metadata"
	"github.com/violetcircus/viviscrobbler/internal/setup"

	// "strconv"
	"strings"
)

// struct for scrobbles
type Scrobble struct {
	trackInfo string
	status    string
	timestamp string
	apiKey    string
	sk        string
	apiSecret string
	method    string
}

func main() {
	log.SetFlags(0)
	setup.Setup() // call setup function (to do)

	// say hi
	fmt.Println("viviscrobbler!")

	// connect to mpd
	conn, err := net.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// set the currently watched track to nothing
	currentlyWatchedTrack := ""

	// main program loop. communicates with mpd
	for {
		// tell mpd to idle and watch for changes in player
		fmt.Fprintln(conn, "idle player")
		// read mpd idle command output until it something in the player changes
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
		// get current song
		fmt.Fprintf(conn, "currentsong\n")
		trackInfo := mapOutput(reader)
		// get status
		fmt.Fprintf(conn, "status\n")
		status := mapOutput(reader)
		// fmt.Println(trackInfo)

		// if the user has told mpd to play:
		state := status["state"]
		if state == "play" {
			title := trackInfo["Title"].(string)
			if title != currentlyWatchedTrack { // check if current track != new track
				currentlyWatchedTrack = title // set current track to new track
				log.Println("state:", state)
				log.Println("title:", title)
				log.Println("Cleaned artist:", metadata.GetArtist(trackInfo))
			} else if status["single"].(float64) == 1 && status["repeat"].(float64) == 1 && status["elapsed"].(int) < 1 {
			}
		}
	}
}

func mapOutput(reader *bufio.Reader) map[string]any {
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// print the status
	fmt.Println("Server:", line)

	// create output map
	output := make(map[string]any)
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
		// fmt.Println(line)
		// put results in the map
		key, value, found := strings.Cut(line, ":")
		if found {
			// if the string is not purely numeric, stick it in the map as a string
			if !isNumeric(strings.TrimSpace(value)) {
				output[strings.TrimSpace(key)] = strings.TrimSpace(value)
			} else {
				// if it is, convert to a 64-bit float and stick it in the map
				value, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Fatal(err)
				}
				output[strings.TrimSpace(key)] = value
			}
		}
	}
	return output
}

// quick helper function to check if a string is purely numeric
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
