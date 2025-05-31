package main

import (
	"bufio"
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal/metadata"
	"github.com/violetcircus/viviscrobbler/internal/setup"
	"log"
	"net"
	// "reflect"
	"strconv"
	"strings"
)

// struct for status info reported by mpd
type Status struct {
	State      string
	Duration   int
	Elapsed    int
	Time       string
	Repeat     int
	Single     int
	Song       int
	Songid     int
	Nextsong   int
	Nextsongid int
}

// struct for song info reported by mpd
type TrackInfo struct {
	Title       string
	Album       string
	AlbumArtist string
	Artist      string
}

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
		trackInfo := getSong(reader)
		// get status
		fmt.Fprintf(conn, "status\n")
		status := mapOutput(reader)
		// fmt.Println(trackInfo)

		// if the user has told mpd to play:
		state := status["state"]
		if state == "play" {
			title := trackInfo.Title
			if title != currentlyWatchedTrack { // check if current track != new track
				currentlyWatchedTrack = title // set current track to new track
				log.Println("state:", state)
				log.Println("title:", title)
				// log.Println("Cleaned artist:", metadata.GetArtist(trackInfo))
				// } else if status["single"] == 1 && status["repeat"] == 1 && status["elapsed"] < 1 {
			}
		}
	}
}

// this could be optimised but im scared of mpd changing stuff around.
func getSong(reader *bufio.Reader) *TrackInfo {
	s := TrackInfo{}
	for key, value := range mapOutput(reader) {
		switch key {
		case "Title":
			s.Title = value
		case "Album":
			s.Album = value
		case "Artist":
			s.Artist = value
		case "AlbumArtist":
			s.AlbumArtist = value
		}
	}
	return &s
}

func getStatus(reader *bufio.Reader) *Status {
	s := Status{}
	for key, value := range mapOutput(reader) {
		switch key {
		case "state":
			s.State = value
		case "repeat":
			repeat, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal(err)
			}
			s.Repeat = repeat
		case "single":
			single, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal(err)
			}
			s.Single = single
		}
	}
	return &s
}

func mapOutput(reader *bufio.Reader) map[string]string {
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
		// fmt.Println(line)
		// put results in the map
		key, value, found := strings.Cut(line, ":")
		if found {
			trackInfo[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	return trackInfo
}
