package main

import (
	"bufio"
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal/metadata"
	"github.com/violetcircus/viviscrobbler/internal/setup"
	"log"
	"net"
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
		trackInfo := metadata.GetSong(reader)
		// get status
		fmt.Fprintf(conn, "status\n")
		status := metadata.GetStatus(reader)
		// fmt.Println(trackInfo)

		// if the user has told mpd to play:
		state := status.State
		log.Println("state:", state)
		if state == "play" {
			title := trackInfo.Title
			if title != currentlyWatchedTrack { // check if current track != new track
				currentlyWatchedTrack = title // set current track to new track
				log.Println("Cleaned artist:", metadata.GetArtist(trackInfo))
				// } else if status["single"] == 1 && status["repeat"] == 1 && status["elapsed"] < 1 {
			}
		}
	}
}
