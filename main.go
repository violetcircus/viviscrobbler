package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
	// "strconv"
	"strings"

	"github.com/violetcircus/viviscrobbler/internal/config"
	"github.com/violetcircus/viviscrobbler/internal/metadata"
	"github.com/violetcircus/viviscrobbler/internal/setup"
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
	config := config.ReadConfig()

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
		// go back to idling after getting status/trackinfo
		// fmt.Fprintln(conn, "idle player")

		// loop checking the state
		elapsed := 0.0
		for {
			time.Sleep(1 * time.Second) // wait one second
			fmt.Fprintf(conn, "status\n")
			status := metadata.GetStatus(reader) // get status
			fmt.Println(status.State)
			switch status.State {
			case "play": // when in play state, calculate percent through song
				elapsed = status.Elapsed
				duration := status.Duration
				percent := elapsed / duration * 100
				fmt.Println("percent:", percent)
				// run scrobble check
				if percent >= config.ScrobbleThreshold && trackInfo.Title != currentlyWatchedTrack {
					currentlyWatchedTrack = trackInfo.Title // set current track to new track
					makeScrobble(trackInfo, status)
					// check if user has the track on repeat, scrobble it if so
				} else if status.Single == 1 && status.Repeat == 1 && status.Elapsed < 1 {
					currentlyWatchedTrack = trackInfo.Title // set current track to new track
					makeScrobble(trackInfo, status)
				}
			case "pause": // when paused, just go back to the start of the loop
				// loop again until something else happens
			case "stop": // when stopped, exit this loop
				elapsed = 0.0
				break
			}
		}
	}
}

func makeScrobble(trackInfo metadata.TrackInfo, status metadata.Status) {
	_ = status
	log.Println("Cleaned artist:", metadata.GetArtist(trackInfo))
}
