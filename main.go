package main

import (
	"bufio"
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal/configreader"
	"github.com/violetcircus/viviscrobbler/internal/metadata"
	"github.com/violetcircus/viviscrobbler/internal/scrobblelogger"
	"github.com/violetcircus/viviscrobbler/internal/secret"
	"github.com/violetcircus/viviscrobbler/internal/setup"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	secret.GetSecrets()
	log.SetFlags(0)
	setup.Setup() // call setup function (to do)
	config := configreader.ReadConfig()

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
	timestamp := "0"

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
		// go back to idling after getting status/trackinfo
		// fmt.Fprintln(conn, "idle player")

		// loop checking the state
		elapsed := 0.0
		for {
			time.Sleep(1 * time.Second) // wait one second

			// get current song from mpd
			fmt.Fprintf(conn, "currentsong\n")
			trackInfo := metadata.GetSong(reader)

			// get current mpd status
			fmt.Fprintf(conn, "status\n")
			status := metadata.GetStatus(reader)

			fmt.Println("state:", status.State)
			switch status.State {
			case "play": // when in play state, calculate percent through song
				if isRepeat(status) {
					timestamp = strconv.FormatInt(time.Now().Unix(), 10)
				} else if trackInfo.Title != currentlyWatchedTrack {
					timestamp = strconv.FormatInt(time.Now().Unix(), 10)
				}
				elapsed = status.Elapsed
				duration := status.Duration
				percent := elapsed / duration * 100

				fmt.Println("percent:", percent)
				// fmt.Println("title:", trackInfo.Title)
				// fmt.Println("currently watched:", currentlyWatchedTrack)

				// run scrobble check
				if percent >= config.ScrobbleThreshold && trackInfo.Title != currentlyWatchedTrack {
					currentlyWatchedTrack = trackInfo.Title // set current track to new track
					makeScrobble(trackInfo, timestamp)
					// check if user has the track on repeat + the song is within the first second, scrobble it if so
				} else if isRepeat(status) && status.Elapsed < 1 {
					currentlyWatchedTrack = trackInfo.Title // set current track to new track
					makeScrobble(trackInfo, timestamp)
				}
			case "pause": // when paused, just go back to the start of the loop
				// loop again until something else happens
			case "stop": // when stopped, exit this loop
				elapsed = 0.0 // not strictly necessary but like,
				break
			}
		}
	}
}

// check if song is on repeat
func isRepeat(status metadata.Status) bool {
	if status.Single == 1 && status.Repeat == 1 {
		return true
	} else {
		return false
	}
}

func makeScrobble(trackInfo metadata.TrackInfo, timestamp string) {
	s := scrobblelogger.LoggedScrobble{}
	s.Title = trackInfo.Title
	s.Artist = trackInfo.Artist
	s.Album = trackInfo.Album
	s.Timestamp = timestamp
	log.Println("Cleaned artist:", metadata.GetArtist(trackInfo.Artist))
	scrobblelogger.WriteScrobble(s)
	fmt.Println(scrobblelogger.ReadScrobble())
}
