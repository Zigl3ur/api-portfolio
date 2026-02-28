package lastfm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3/client"
)

type lasfmData struct {
	RecentTracks struct {
		Track []struct {
			Artist struct {
				Name string `json:"#text"`
			} `json:"artist"`
			Images []struct {
				Size string `json:"size"`
				Url  string `json:"#text"`
			} `json:"image"`
			Album struct {
				Name string `json:"#text"`
			} `json:"album"`
			TrackName string `json:"name"`
			Attr      struct {
				IsPlaying string `json:"nowplaying"`
			} `json:"@attr"`
			TrackUrl string `json:"url"`
		} `json:"track"`
	} `json:"recenttracks"`
}

type FormatedData struct {
	IsListening bool `json:"isListening"`
	Track       *struct {
		Artist    string `json:"artist,omitempty"`
		Album     string `json:"album,omitempty"`
		TrackName string `json:"name,omitempty"`
		Image     string `json:"image,omitempty"`
		Url       string `json:"url,omitempty"`
	} `json:"track,omitempty"`
}

func MusicHandler(apiKey string) (*FormatedData, error) {

	dataFormat := &FormatedData{
		IsListening: false,
	}

	cc := client.New()
	resp, err := cc.Get(fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=zigl3ur&api_key=%s&format=json", apiKey))
	if err != nil {
		return dataFormat, err
	}

	defer resp.Close()

	var data lasfmData
	if err = json.Unmarshal(resp.Body(), &data); err != nil {
		return dataFormat, err
	}

	tracks := data.RecentTracks.Track
	dataFormat.IsListening = len(tracks) > 0 && tracks[0].Attr.IsPlaying == "true"

	if dataFormat.IsListening {
		if len(tracks[0].Images) > 0 {
			for i := range tracks[0].Images {
				if tracks[0].Images[i].Size == "large" {
					dataFormat.Track.Image = getAr0ImageUrl(tracks[0].Images[i].Url)
					break
				}
			}
		}

		dataFormat.Track.TrackName = tracks[0].TrackName
		dataFormat.Track.Album = tracks[0].Album.Name
		dataFormat.Track.Artist = tracks[0].Artist.Name
		dataFormat.Track.Url = tracks[0].TrackUrl
	}

	return dataFormat, nil
}

// we replace the tag '174s' with 'ar0' in the image url to have a better quality,
// it doesnt seems that we can get it from the api directly
func getAr0ImageUrl(baseImage string) string {
	parts := strings.Split(baseImage, "174s")
	if len(parts) != 2 {
		return baseImage
	}

	return strings.Join(parts, "ar0")
}
