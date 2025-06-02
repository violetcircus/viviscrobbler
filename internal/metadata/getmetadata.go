package metadata

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// struct for status info reported by mpd
type Status struct {
	State      string
	Duration   float64
	Elapsed    float64
	Time       float64
	Repeat     int
	Single     int
	Song       int
	SongID     int
	NextSong   int
	NextSongID int
}

// struct for song info reported by mpd
type TrackInfo struct {
	Title       string
	Album       string
	AlbumArtist string
	Artist      string
}

// populate track info struct
func GetSong(reader *bufio.Reader) TrackInfo {
	s := TrackInfo{}
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// print the status
	fmt.Println("Server:", line)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line) // sanitise lines
		// break loop when track info is complete or an error is given
		if line == "OK" || strings.HasPrefix(line, "ACK") {
			fmt.Println("song response:", line)
			break
		}
		// populate struct
		key, value, found := strings.Cut(line, ": ")
		if found {
			switch key {
			case "Title":
				s.Title = strings.TrimSpace(value)
			case "Album":
				s.Album = strings.TrimSpace(value)
			case "Artist":
				s.Artist = strings.TrimSpace(value)
			case "AlbumArtist":
				s.AlbumArtist = strings.TrimSpace(value)
			}
		}
	}
	return s
}

func GetStatus(reader *bufio.Reader) Status {
	s := Status{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line) // sanitise lines
		// break loop when status is retrieved or error is given
		if line == "OK" || strings.HasPrefix(line, "ACK") {
			fmt.Println("status response:", line)
			break
		}
		// populate struct
		key, value, found := strings.Cut(line, ": ")
		if found {
			switch key {
			case "state":
				s.State = value
			case "time":
				// split on the colon and convert the left value to a float
				t, _, found := strings.Cut(value, ":")
				if found {
					time, err := strconv.ParseFloat(t, 64)
					if err != nil {
						log.Fatal(err)
					}
					s.Time = time
				}
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
			case "duration":
				duration, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Fatal(err)
				}
				s.Duration = duration
			case "elapsed":
				elapsed, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Fatal(err)
				}
				s.Elapsed = elapsed
			case "song":
				song, err := strconv.Atoi(value)
				if err != nil {
					log.Fatal(err)
				}
				s.Song = song
			case "songid":
				songid, err := strconv.Atoi(value)
				if err != nil {
					log.Fatal(err)
				}
				s.SongID = songid
			case "nextsong":
				nextsong, err := strconv.Atoi(value)
				if err != nil {
					log.Fatal(err)
				}
				s.NextSong = nextsong
			case "nextsongid":
				nextsongid, err := strconv.Atoi(value)
				if err != nil {
					log.Fatal(err)
				}
				s.NextSongID = nextsongid
			}
		}
	}
	return s
}

// func mapOutput(reader *bufio.Reader) map[string]string {
// 	line, err := reader.ReadString('\n')
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// print the status
// 	fmt.Println("Server:", line)
//
// 	// create track info map
// 	trackInfo := make(map[string]string)
// 	// loop over song info
// 	for {
// 		line, err := reader.ReadString('\n')
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		line = strings.TrimSpace(line)
// 		// check if response finished, stop if it has
// 		if line == "OK" || strings.HasPrefix(line, "ACK") {
// 			fmt.Println("Response:", line)
// 			break
// 		}
// 		// output results to terminal
// 		// fmt.Println(line)
// 		// put results in the map
// 		key, value, found := strings.Cut(line, ":")
// 		if found {
// 			trackInfo[strings.TrimSpace(key)] = strings.TrimSpace(value)
// 		}
// 	}
// 	return trackInfo
// }
