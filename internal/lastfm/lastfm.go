package lastfm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3/client"
)

type LastFM struct {
	ApiKey string
}

func NewLastFM(apiKey string) *LastFM {
	return &LastFM{
		ApiKey: apiKey,
	}
}

type lastfmImageData struct {
	Size string `json:"size"`
	Url  string `json:"#text"`
}

type lastfmRecentTrack struct {
	RecentTracks struct {
		Track []struct {
			Artist struct {
				Name string `json:"#text"`
			} `json:"artist"`
			Images []lastfmImageData `json:"image"`
			Album  struct {
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

type CurrentlyListening struct {
	IsListening bool   `json:"isListening"`
	Track       *Track `json:"track,omitempty"`
}

type lastfmAlbumInfo struct {
	Album struct {
		Artist string `json:"artist"`

		Tags struct {
			Tag []struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"tag"`
		} `json:"tags"`

		Tracks struct {
			Track []struct {
				Duration int64  `json:"duration"`
				Url      string `json:"url"`
				Name     string `json:"name"`
				Attr     struct {
					Rank int64 `json:"rank"`
				} `json:"@attr"`
				Artist struct {
					Url  string `json:"url"`
					Name string `json:"name"`
				} `json:"artist"`
			} `json:"track"`
		} `json:"tracks"`
	} `json:"album"`
}

type Album struct {
	Artist    string `json:"artist"`
	Url       string `json:"url"`
	AlbumName string `json:"name"`
	Image     string `json:"image"`
	Playcount string `json:"playcount"`
}

type Track struct {
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	TrackName string `json:"name"`
	Image     string `json:"image"`
	Url       string `json:"url"`
}

func (lfm *LastFM) GetCurrentlyPlaying() (*CurrentlyListening, error) {
	dataFormat := &CurrentlyListening{
		IsListening: false,
	}

	cc := client.New()
	resp, err := cc.Get(fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=zigl3ur&api_key=%s&format=json", lfm.ApiKey))
	if err != nil {
		return dataFormat, err
	}

	defer resp.Close()

	var data lastfmRecentTrack
	if err = json.Unmarshal(resp.Body(), &data); err != nil {
		return dataFormat, err
	}

	tracks := data.RecentTracks.Track
	dataFormat.IsListening = len(tracks) > 0 && tracks[0].Attr.IsPlaying == "true"

	if dataFormat.IsListening {
		dataFormat.Track = &Track{
			Artist:    tracks[0].Artist.Name,
			Album:     tracks[0].Album.Name,
			TrackName: tracks[0].TrackName,
			Url:       tracks[0].TrackUrl,
		}

		if len(tracks[0].Images) > 0 {
			for i := range tracks[0].Images {
				if tracks[0].Images[i].Size == "large" {
					dataFormat.Track.Image = getAr0ImageUrl(tracks[0].Images[i].Url)
					break
				}
			}
		}
	}

	return dataFormat, nil
}

func (lfm *LastFM) GetTopAlbums() ([]Album, error) {
	cc := client.New()
	resp, err := cc.Get(fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?method=user.getTopAlbums&user=zigl3ur&api_key=%s&limit=100&format=json", lfm.ApiKey))
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var raw struct {
		TopAlbums struct {
			Albums []struct {
				Artist struct {
					Name string `json:"name"`
				} `json:"artist"`
				Url       string            `json:"url"`
				AlbumName string            `json:"name"`
				Images    []lastfmImageData `json:"image"`
				Playcount string            `json:"playcount"`
			} `json:"album"`
		} `json:"topalbums"`
	}

	if err = json.Unmarshal(resp.Body(), &raw); err != nil {
		return nil, err
	}

	albums := make([]Album, 0, len(raw.TopAlbums.Albums))
	for _, a := range raw.TopAlbums.Albums {
		album := Album{
			Artist:    a.Artist.Name,
			Url:       a.Url,
			AlbumName: a.AlbumName,
			Playcount: a.Playcount,
		}
		for _, img := range a.Images {
			if img.Size == "large" {
				album.Image = getAr0ImageUrl(img.Url)
				break
			}
		}
		albums = append(albums, album)
	}

	return albums, nil
}

func (lfm *LastFM) GetAlbumInfo(artist, album string) (*lastfmAlbumInfo, error) {
	dataFormat := &lastfmAlbumInfo{}

	cc := client.New()
	resp, err := cc.Get(fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?method=album.getinfo&api_key=%s&artist=%s&album=%s&format=json", lfm.ApiKey, artist, album))
	if err != nil {
		return dataFormat, err
	}

	defer resp.Close()

	if err = json.Unmarshal(resp.Body(), &dataFormat); err != nil {
		return dataFormat, err
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
