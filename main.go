package main

import (
	"bufio"
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal"
	"log"
	"net"
	"strings"
)

type status struct {
	state      string
	duration   int
	elapsed    int
	time       string
	repeat     int
	single     int
	song       int
	songid     int
	nextsong   int
	nextsongid int
}

type trackInfo struct {
	title       string
	album       string
	albumArtist string
	artist      string
}

type scrobble struct {
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
	internal.Setup() // call setup function (to do)
	fmt.Println("viviscrobbler!")
	conn, err := net.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// set the currently watched track to nothing
	currentlyWatchedTrack := ""

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
		// ask for current song
		fmt.Fprintf(conn, "currentsong\n")
		trackInfo := mapOutput(reader)
		fmt.Fprintf(conn, "status\n")
		status := mapOutput(reader)
		fmt.Println(trackInfo)
		state := status["state"]

		if state == "play" {
			title := trackInfo["Title"]
			if title != currentlyWatchedTrack {
				currentlyWatchedTrack = title // set current track
				log.Println("state:", state)
				log.Println("title:", title)
				log.Println("Cleaned artist:", internal.GetArtist(trackInfo))
				// } else if status["single"] == 1 && status["repeat"] == 1 && status["elapsed"] < 1 {
			}
		}
	}
}

// this could be optimised but im scared of mpd changing stuff around.
func GetSong(reader *bufio.Reader) *trackInfo {
	s := trackInfo{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(line, "Title:") {
			title, found := strings.CutPrefix(line, "Title: ")
			if !found {
				log.Fatal("no title :(")
			}
			s.title = title
		}
		if strings.HasPrefix(line, "Artist:") {
			artist, found := strings.CutPrefix(line, "Artist: ")
			if !found {
				log.Fatal("no title :(")
			}
			s.title = strings.TrimSpace(artist)
		}
		if strings.HasPrefix(line, "Album:") {
			album, found := strings.CutPrefix(line, "Album: ")
			if !found {
				log.Fatal("no album :(")
			}
			s.album = album
		}
		if strings.HasPrefix(line, "AlbumArtist:") {
			albumArtist, found := strings.CutPrefix(line, "AlbumArtist: ")
			if !found {
				log.Fatal("no albumArtist :(")
			}
			s.albumArtist = albumArtist
		}
		if line == "OK" || strings.HasPrefix(line, "ACK") {
			fmt.Println("Response:", line)
			break
		}
	}
	return &s
}

func getStatus(reader *bufio.Reader) *status {
	s := status{}
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
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
