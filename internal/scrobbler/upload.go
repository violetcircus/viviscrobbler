package scrobbler

// scrobble to be uploaded to lastfm
type Scrobble struct {
	Artist    string
	Album     string
	Timestamp string
	Title     string
	ApiKey    string
	Secret    string
	SKey      string
}

func UploadScrobbles() {

}

func UpdateNowPlaying() {

}
