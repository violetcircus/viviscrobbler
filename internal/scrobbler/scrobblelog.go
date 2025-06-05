package scrobbler

import (
	"encoding/csv"
	"github.com/violetcircus/viviscrobbler/internal/metadata"
	"log"
	"os"
)

type LoggedScrobble struct {
	Artist    string
	Title     string
	Album     string
	Timestamp string
}

func WriteScrobble(scrobble LoggedScrobble) {
	f := "/home/violet/.config/vvscrob/logFile.tsv"
	log.Println("writing scrobble")
	logFile, err := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	w := csv.NewWriter(logFile)
	w.Comma = '\t'
	row := []string{scrobble.Artist, scrobble.Album, scrobble.Title, scrobble.Timestamp}
	if err := w.Write(row); err != nil {
		log.Fatal("error writing to file")
	}
	w.Flush()
}

func ReadScrobble() LoggedScrobble {
	s := LoggedScrobble{}
	f := "/home/violet/.config/vvscrob/logFile.tsv"
	logFile, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	r := csv.NewReader(logFile)
	r.Comma = '\t'
	scrobbles, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("scrobbles:", scrobbles)
	scrobble := scrobbles[0]
	s.Artist = metadata.GetArtist(scrobble[0])
	s.Album = scrobble[1]
	s.Title = scrobble[2]
	s.Timestamp = scrobble[3]
	return s
}
