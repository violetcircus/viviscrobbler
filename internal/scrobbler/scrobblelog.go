package scrobbler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal/configreader"
	"github.com/violetcircus/viviscrobbler/internal/metadata"
	"io"
	"log"
	"os"
	"sync"
)

type LoggedScrobble struct {
	Artist    string
	Title     string
	Album     string
	Timestamp string
}

var m sync.Mutex

// write scrobble to file
func WriteScrobble(scrobble LoggedScrobble) {
	m.Lock()
	f := configreader.ConfigLocation + "logFile.tsv"
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
	m.Unlock()
}

// read scrobbles from file
func ReadScrobble(wg *sync.WaitGroup) LoggedScrobble {
	defer wg.Done()
	f := configreader.ConfigLocation + "logFile.tsv"

	for {
		m.Lock()
		s := LoggedScrobble{}
		logFile, err := os.OpenFile(f, os.O_RDWR, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		r := csv.NewReader(logFile)
		r.Comma = '\t'
		scrobbles, err := r.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		if len(scrobbles) > 0 {
			scrobble := scrobbles[0]
			s = LoggedScrobble{
				Artist:    metadata.GetArtist(scrobble[0]),
				Album:     scrobble[1],
				Title:     scrobble[2],
				Timestamp: scrobble[3],
			}
			log.Printf("scrobble: %s", s)
			if UploadScrobbles(s) {
				popLine(logFile)
			}
		} else {
			logFile.Close()
			m.Unlock()
			continue
		}
		logFile.Close()
		m.Unlock()
	}
}

// delete first line of file
func popLine(f *os.File) {
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Size() > 0 {
		buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

		_, err = f.Seek(0, io.SeekStart) // move file pointer to start
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(buf, f)
		if err != nil {
			log.Fatal(err)
		}

		line, err := buf.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("popped line:", string(line))

		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}

		nw, err := io.Copy(f, buf)
		if err != nil {
			fmt.Println("copy error")
			log.Fatal(err)
		}

		err = f.Truncate(nw)
		if err != nil {
			fmt.Println("truncate error")
			log.Fatal(err)
		}
		err = f.Sync()
		if err != nil {
			fmt.Println("sync error")
			log.Fatal(err)
		}

	}
}
