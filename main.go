package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal/configreader"
	"github.com/violetcircus/viviscrobbler/internal/metadata"
	"github.com/violetcircus/viviscrobbler/internal/scrobbler"
	"github.com/violetcircus/viviscrobbler/internal/setup"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	log.SetFlags(0)
	setup.Setup()

	args := os.Args
	handleArgs(args)

	// start the other thread that reads the log file and scrobbles the entries in it
	var wg sync.WaitGroup
	wg.Add(1)
	go scrobbler.ReadScrobble(&wg)

	config := configreader.ReadConfig()

	// say hi
	fmt.Println("viviscrobbler!")

	// connect to mpd
	conn, err := net.Dial("tcp", config.ServerAddress+":"+config.ServerPort)
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
		// read mpd idle command output until something in the player changes
		buf := make([]byte, 512) // create reusable buffer to avoid ballooning memory usage since this will loop a lot in the background
		for {
			n, err := reader.Read(buf)
			if err != nil {
				log.Fatal("mpd read error:", err)
			}
			if bytes.Contains(buf[:n], []byte("changed: player")) {
				break
			}
			time.Sleep(time.Second / 2)
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
				if trackInfo.Title != currentlyWatchedTrack {
					fmt.Println("current track:", currentlyWatchedTrack)
					scrobbler.UpdateNowPlaying(trackInfo)
				}
				if isRepeat(status) {
					timestamp = strconv.FormatInt(time.Now().Unix(), 10)
				} else if trackInfo.Title != currentlyWatchedTrack {
					timestamp = strconv.FormatInt(time.Now().Unix(), 10)
				}
				elapsed = status.Elapsed
				duration := status.Duration
				percent := elapsed / duration * 100

				log.Println("percent:", percent)

				// check if user has the track on repeat + the song is within the first second, scrobble it if so
				if isRepeat(status) && status.Elapsed < 1 {
					currentlyWatchedTrack = trackInfo.Title // set current track to new track
					makeScrobble(trackInfo, timestamp)
					// run scrobble
				} else if percent >= config.ScrobbleThreshold && trackInfo.Title != currentlyWatchedTrack {
					currentlyWatchedTrack = trackInfo.Title // set current track to new track
					makeScrobble(trackInfo, timestamp)
				}
			case "pause": // when paused, just go back to the start of the loop
				// loop again until something else happens
			case "stop": // when stopped, exit this loop
				elapsed = 0.0 // not strictly necessary but like, just in case
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

// send song to be written as a scrobble in the log
func makeScrobble(trackInfo metadata.TrackInfo, timestamp string) {
	s := scrobbler.LoggedScrobble{
		Title:     trackInfo.Title,
		Artist:    trackInfo.Artist,
		Album:     trackInfo.Album,
		Timestamp: timestamp,
	}
	log.Println("Cleaned artist:", metadata.GetArtist(trackInfo.Artist))
	scrobbler.WriteScrobble(s)
	// fmt.Println(scrobbler.ReadScrobble())
}

// handle arguments: write default config to file when program is run with config, otherwise run rockbox
// scrobbling on provided filepath
func handleArgs(args []string) {
	if len(args) > 1 {
		arg := args[1]
		switch arg {
		case "config":
			setup.WriteConfig()
		default:
			f, err := os.Stat(arg)
			if err != nil {
				log.Fatal("not a file!! please use the full path to your scrobble log. ~ shortcut for /home/user is okay.", f)
			}
			// allow user to use ~ home shortcut in path
			if strings.HasPrefix(arg, "~") {
				h, err := os.UserHomeDir()
				if err != nil {
					log.Fatal(err)
				}
				scrobbler.ReadRockboxLog(strings.Replace(arg, "~", h, 1))
			} else {
				scrobbler.ReadRockboxLog(arg)
			}
		}
	}
}
