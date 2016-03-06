package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const LimitTracks = "100"

type Artist struct {
	Id             int
	Name           string
	Link           string
	Share          string
	Picture        string
	Picture_small  string
	Picture_medium string
	Picture_big    string
	Nb_album       int
	Nb_fan         int
	Radio          bool
	Tracklist      string
}

type Album struct {
	Id              int
	Title           string
	Upc             string
	Link            string
	Share           string
	Cover           string
	Cover_small     string
	Cover_medium    string
	Cover_big       string
	Genres          struct{ Data []Genre }
	Lable           string
	Nb_tracks       int
	Duration        int
	Fans            int
	Rating          int
	Release_date    string
	Record_type     string
	Available       bool
	Tracklist       string
	Explicit_lyrics bool
	Artist          Artist
	Tracks          struct{ Data []Track }
}

type Track struct {
	Id                  int
	Readable            bool
	Title               string
	Title_sort          string
	Title_version       string
	Unseen              bool
	Link                string
	Share               string
	Duration            int
	Track_position      int
	Disk_number         int
	Rank                int
	Release_date        string
	Explicit_lyrics     bool
	Preview             string
	Bpm                 float64
	Gain                float64
	Available_countries []string
	Contributors        []Artist
	Artist              Artist
	Album               Album
}

type Genre struct {
	Id             int
	Name           string
	Picture        string
	Picture_small  string
	Picture_medium string
	Picture_big    string
}

type ErrorData struct {
	Error struct {
		Type    string
		Message string
		Code    int
	}
}

func GetJSON(url string, data interface{}) error {
	log.Println("Deezer: GetJSON: ", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return err2
	}
	var errorData ErrorData
	err2 = json.Unmarshal(body, &errorData)
	if err2 == nil {
		if errorData.Error.Message != "" {
			return errors.New(errorData.Error.Message)
		}
	}
	err2 = json.Unmarshal(body, data)
	if err2 != nil {
		return err2
	}
	return nil
}

func DeezerGetJSON(path string, q url.Values, data interface{}) error {
	log.Println("Deezer: DeezerGetJSON: ", path, " query: ", q)
	u, err := url.Parse(settings.BaseURL + "/" + path)
	if q == nil {
		q = url.Values{}
	}
	if err != nil {
		log.Println("Deezer: DeezerGetJSON: ", err)
		return err
	}
	if settings.HasAccess() {
		q.Set("access_token", settings.AccessToken)
	}
	u.RawQuery = q.Encode()
	return GetJSON(u.String(), data)
}

func GetArtist(id int) (Artist, error) {
	var artist Artist
	err := DeezerGetJSON("artist/"+strconv.Itoa(id), nil, &artist)
	if err != nil {
		return Artist{}, err
	}
	return artist, nil
}

func GetArtistTop(id int) ([]Track, error) {
	type Data struct {
		Data []Track
	}
	var data Data
	err := DeezerGetJSON("artist/"+strconv.Itoa(id)+"/top", url.Values{"limit": {settings.LimitTracksAttribute}}, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}

func GetAlbum(id int) (Album, error) {
	var album Album
	err := DeezerGetJSON("album/"+strconv.Itoa(id), nil, &album)
	if err != nil {
		return Album{}, err
	}
	return album, nil
}

func GetTrack(id int) (Track, error) {
	var track Track
	err := DeezerGetJSON("track/"+strconv.Itoa(id), nil, &track)
	if err != nil {
		return Track{}, err
	}
	return track, nil
}

func GetArtistsFromGenre(genreId string) ([]Artist, error) {
	type Data struct {
		Data []Artist
	}
	var data Data
	err := DeezerGetJSON("genre/"+genreId+"/artists", url.Values{"limit": {settings.LimitResultsAttribute}}, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}

func GetTracksFromAlbum(id int) ([]Track, error) {
	type Data struct {
		Data []Track
		Next string
	}
	var data Data
	err := DeezerGetJSON("album/"+strconv.Itoa(id)+"/tracks", url.Values{"limit": {LimitTracks}}, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}

func QueryTracks(query string) ([]Track, error) {
	query = url.QueryEscape(query)
	type Data struct {
		Data []Track
	}
	var data Data
	err := DeezerGetJSON("search/track", url.Values{"q": {query}, "order": {settings.SortAttribute}, "limit": {settings.LimitResultsAttribute}}, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}

func QueryRecommendedTracks() ([]Track, error) {
	type Data struct {
		Data []Track
	}
	var data Data
	err := DeezerGetJSON("user/me/recommendations/tracks", nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}

func QueryHistoryTracks() ([]Track, error) {
	type Data struct {
		Data []Track
	}
	var data Data
	err := DeezerGetJSON("user/me/history", nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}

func QueryMy25Tracks() ([]Track, error) {
	type Data struct {
		Data []Track
	}
	var data Data
	err := DeezerGetJSON("user/me/charts", nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}

func QueryArtists(query string) ([]Artist, error) {
	query = url.QueryEscape(query)
	type Data struct {
		Data []Artist
	}
	var data Data
	err := DeezerGetJSON("search/artist", url.Values{"q": {query}, "order": {settings.SortAttribute}, "limit": {settings.LimitResultsAttribute}}, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}

func QueryRecommendedArtists() ([]Artist, error) {
	type Data struct {
		Data []Artist
	}
	var data Data
	err := DeezerGetJSON("user/me/recommendations/artists", nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}

func QueryAlbums(query string) ([]Album, error) {
	query = url.QueryEscape(query)
	type Data struct {
		Data []Album
	}
	var data Data
	err := DeezerGetJSON("search/album", url.Values{"q": {query}, "order": {settings.SortAttribute}, "limit": {settings.LimitResultsAttribute}}, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}

func QueryRecommendedAlbums() ([]Album, error) {
	type Data struct {
		Data []Album
	}
	var data Data
	err := DeezerGetJSON("user/me/recommendations/albums", nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return data.Data, errors.New("Empty response")
	}
	return data.Data, nil
}
