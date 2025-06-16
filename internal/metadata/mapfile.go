package metadata

import (
	"encoding/csv"
	"github.com/violetcircus/viviscrobbler/internal/configreader"
	"log"
	"os"
)

// artists stored in the map
type artistEntry struct {
	oldString string
	cleaned   string
}

// check map file for cleaned artist name for received metadata string
func checkMapFile(artist string) string {
	f := configreader.GetConfigDir() + "mapFile.tsv"
	mapFile, err := os.OpenFile(f, os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer mapFile.Close()

	r := csv.NewReader(mapFile)
	r.Comma = '\t'
	entries, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// there's probably a less expensive way to do this - i wonder if there's a way to have ReadAll produce a map instead of a slice? idk
	if len(entries) > 0 {
		for _, entry := range entries {
			if entry[0] == artist {
				log.Println("cleaned artist found!", entry[1])
				return entry[1]
			}
		}
	} else {
		log.Print("artist map empty!")
	}

	return ""
}

// write cleaned artist name to map file
func writeMapFile(artist string, cleanedArtist string) {
	s := artistEntry{
		oldString: artist,
		cleaned:   cleanedArtist,
	}
	f := configreader.GetConfigDir() + "mapFile.tsv"
	mapFile, err := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer mapFile.Close()

	w := csv.NewWriter(mapFile)
	w.Comma = '\t'

	row := []string{s.oldString, s.cleaned}
	if err := w.Write(row); err != nil {
		log.Println("error writing artist to file")
	}
	w.Flush()
}
