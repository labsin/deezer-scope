package main

import (
	//"encoding/json"
	//"launchpad.net/go-onlineaccounts/v1"
	"launchpad.net/go-unityscopes/v2"
	"log"
	"strconv"
)

const loginNagTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "grid",
    "card-total_results": "large",
    "card-background": "color:///#DD4814"
  },
  "components": {
    "title": "title",
    "background": "background",
    "art" : {
      "aspect-ratio": 100.0
    }
  }
}`

const trackCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "grid",
    "card-size": "medium"
  },
  "components": {
    "title": "title",
    "subtitle": "album",
    "attributes": "artistAttr",
	"art": {
		"field": "art"
	}
  }
}`

const artistCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "grid",
    "card-size": "medium"
  },
  "components": {
    "title": "title",
	"attributes": "attributes",
	"art": {
		"field": "art"
	}
  }
}`

const albumCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "grid",
    "card-size": "medium"
  },
  "components": {
    "title": "title",
    "subtitle": "artist",
	"attributes": "attributes",
	"art": {
		"field": "art"
	}
  }
}`

const artistFromGenreCategoryTemplate = `{
  "schema-version": 1,
  "template": {
    "category-layout": "grid",
    "card-size": "medium"
  },
  "components": {
    "title": "title",
	"art": {
		"field": "art"
	}
  }
}`

var genreIds = [24]int{
	0,
	132,
	116,
	152,
	113,
	165,
	85,
	106,
	466,
	144,
	129,
	52,
	98,
	173,
	464,
	2,
	12,
	16,
	153,
	75,
	459,
	81,
	95,
	197,
}

var genreName = [24]string{
	"All",
	"Pop",
	"Rap/Hip Hop",
	"Rock",
	"Dance",
	"R&B/Soul/Funk",
	"Alternative",
	"Electro",
	"Folk",
	"Reggae",
	"Jazz",
	"French Chanson",
	"Classical",
	"Films/Games",
	"Metal",
	"African Music",
	"Arabic Music",
	"Asian Music",
	"Blues",
	"Brazilian Music",
	"German music",
	"Indian Music",
	"Kids",
	"Latin Music",
}

// SCOPE ***********************************************************************

var scope_interface scopes.Scope

type DeezerScope struct {
	ClientId string
	//TODO:	Accounts *accounts.Watcher
	base *scopes.ScopeBase
}

func (s *DeezerScope) SetScopeBase(base *scopes.ScopeBase) {
	s.base = base
}

func (s *DeezerScope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {
	var resultType string
	err := result.Get("type", &resultType)
	if err != nil {
		return err
	}
	var resultId int
	err = result.Get("id", &resultId)
	if err != nil {
		return err
	}
	switch resultType {
	case "track":
		{
			track, err := GetTrack(resultId)
			if err != nil {
				return err
			}
			widget_header := scopes.NewPreviewWidget("hdr", "header")
			widget_header.AddAttributeValue("title", track.Title)
			widget_header.AddAttributeValue("subtitle", track.Artist.Name)
			widget_picture := scopes.NewPreviewWidget("picture", "image")
			widget_picture.AddAttributeValue("source", track.Album.Cover_big)
			widget_release := scopes.NewPreviewWidget("release", "text")
			widget_release.AddAttributeValue("text", "Released on "+track.Release_date)
			widget_duration := scopes.NewPreviewWidget("duration", "text")
			widget_duration.AddAttributeValue("text", "Duration of "+SecondsToString(track.Duration))
			type att map[string]interface{}
			widget_actions := scopes.NewPreviewWidget("actions", "actions")
			actions_atts := make([]att, 0, 10)
			actions_atts = append(actions_atts, att{
				"id":    "link",
				"label": "Link",
			})
			widget_actions.AddAttributeValue("actions", actions_atts)
			widget_audio := scopes.NewPreviewWidget("audio", "audio")
			audio_tracks_atts := make([]att, 0, 1)
			audio_tracks_atts = append(audio_tracks_atts, att{
				"title":    track.Title,
				"subtitle": track.Album.Title,
				"source":   track.Preview,
				"length":   track.Duration,
			})
			widget_audio.AddAttributeValue("tracks", audio_tracks_atts)
			reply.PushWidgets(widget_header, widget_picture,
				widget_release, widget_duration, widget_actions, widget_audio)
		}
	case "album":
		{
			album, err := GetAlbum(resultId)
			if err != nil {
				return err
			}
			widget_header := scopes.NewPreviewWidget("hdr", "header")
			widget_header.AddAttributeValue("title", album.Title)
			widget_header.AddAttributeValue("subtitle", album.Artist.Name)
			widget_picture := scopes.NewPreviewWidget("picture", "image")
			widget_picture.AddAttributeValue("source", album.Cover_big)
			widget_no_tracks := scopes.NewPreviewWidget("no_tracks", "text")
			widget_no_tracks.AddAttributeValue("text", ` <img src="image://theme/stock_music" align="middle" width="15" height="15">`+strconv.Itoa(album.Nb_tracks))
			widget_fans := scopes.NewPreviewWidget("fans", "text")
			widget_fans.AddAttributeValue("text", ` <img src="image://theme/starred" align="middle" width="15" height="15">`+strconv.Itoa(album.Fans))
			widget_release := scopes.NewPreviewWidget("release", "text")
			widget_release.AddAttributeValue("text", "Released on "+album.Release_date)
			widget_duration := scopes.NewPreviewWidget("duration", "text")
			widget_duration.AddAttributeValue("text", "Duration of "+SecondsToString(album.Duration))
			type att map[string]interface{}
			widget_actions := scopes.NewPreviewWidget("actions", "actions")
			actions_atts := make([]att, 0, 10)
			actions_atts = append(actions_atts, att{
				"id":    "link",
				"label": "Link",
			})
			widget_actions.AddAttributeValue("actions", actions_atts)
			widget_audio := scopes.NewPreviewWidget("audio", "audio")
			audio_tracks_atts := make([]att, 0, 10)
			tracks, err2 := GetTracksFromAlbum(resultId)
			if err2 != nil {
				return err2
			}
			tracks_in_disc := make(map[int]int)
			for _, track := range tracks {
				no_tracks, ok := tracks_in_disc[track.Disk_number]
				if ok {
					if no_tracks < track.Track_position {
						tracks_in_disc[track.Disk_number] = track.Track_position
					}
				} else {
					tracks_in_disc[track.Disk_number] = track.Track_position
				}
			}
			var max_tracks_in_disc int
			for _, no_tracks := range tracks_in_disc {
				if no_tracks > max_tracks_in_disc {
					max_tracks_in_disc = no_tracks
				}
			}
			discs := len(tracks_in_disc)
			for _, track := range tracks {
				var trackTitle string
				if discs > 1 {
					trackTitle = "<b>" + PadWith0(track.Disk_number, discs)
					trackTitle += PadWith0(track.Track_position, max_tracks_in_disc) + "</b> " + track.Title
				} else {
					trackTitle = "<b>" + PadWith0(track.Track_position, tracks_in_disc[track.Disk_number]) + "</b> " + track.Title
				}
				audio_tracks_atts = append(audio_tracks_atts, att{
					"title":    trackTitle,
					"subtitle": track.Album.Title,
					"source":   track.Preview,
					"length":   track.Duration,
				})
			}
			widget_audio.AddAttributeValue("tracks", audio_tracks_atts)
			layout_1col := scopes.NewColumnLayout(1)
			layout_1col.AddColumn("hdr", "picture", "release", "duration", "no_tracks", "fans", "actions", "audio")
			layout_2col := scopes.NewColumnLayout(2)
			layout_2col.AddColumn("hdr", "picture", "release", "duration", "no_tracks", "fans", "actions")
			layout_2col.AddColumn("audio")
			reply.RegisterLayout(layout_1col, layout_2col)
			reply.PushWidgets(widget_header, widget_picture, widget_no_tracks,
				widget_fans, widget_release, widget_duration, widget_actions, widget_audio)
		}
	case "artist":
		{
			artist, err := GetArtist(resultId)
			if err != nil {
				return err
			}
			var top_tracks []Track
			top_tracks, err = GetArtistTop(resultId)
			if err != nil {
				return err
			}
			widget_header := scopes.NewPreviewWidget("hdr", "header")
			widget_header.AddAttributeValue("title", artist.Name)
			widget_picture := scopes.NewPreviewWidget("picture", "image")
			widget_picture.AddAttributeValue("source", artist.Picture_big)
			widget_albums := scopes.NewPreviewWidget("no_albums", "text")
			widget_albums.AddAttributeValue("text", ` <img src="`+s.base.ScopeDirectory()+`/album.svg`+`" align="middle" width="15" height="15">`+strconv.Itoa(artist.Nb_album))
			widget_fans := scopes.NewPreviewWidget("fans", "text")
			widget_fans.AddAttributeValue("text", ` <img src="image://theme/starred" align="middle" width="15" height="15">`+strconv.Itoa(artist.Nb_fan))
			type att map[string]interface{}
			widget_actions := scopes.NewPreviewWidget("actions", "actions")
			actions_atts := make([]att, 0, 10)
			actions_atts = append(actions_atts, att{
				"id":    "link",
				"label": "Link",
			})
			widget_actions.AddAttributeValue("actions", actions_atts)
			widget_audio := scopes.NewPreviewWidget("audio", "audio")
			audio_tracks_atts := make([]att, 0, 10)
			for _, track := range top_tracks {
				audio_tracks_atts = append(audio_tracks_atts, att{
					"title":    track.Title,
					"subtitle": track.Album.Title,
					"source":   track.Preview,
					"length":   track.Duration,
				})
			}
			widget_audio.AddAttributeValue("tracks", audio_tracks_atts)
			layout_1col := scopes.NewColumnLayout(1)
			layout_1col.AddColumn("hdr", "picture", "no_albums", "fans", "actions", "audio")
			layout_2col := scopes.NewColumnLayout(2)
			layout_2col.AddColumn("hdr", "picture", "no_albums", "fans", "actions")
			layout_2col.AddColumn("audio")
			reply.RegisterLayout(layout_1col, layout_2col)
			reply.PushWidgets(widget_header, widget_picture, widget_albums,
				widget_fans, widget_actions, widget_audio)
		}
	}
	return nil
}

func (s *DeezerScope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {

	var queryString string = query.QueryString()
	//var queryString string = `a`

	if queryString != "" {
		var tracks, err = QueryTracks(queryString)
		if err != nil {
			return err
		}
		category := reply.RegisterCategory("query_track", "Found tracks", "", trackCategoryTemplate)
		for _, track := range tracks {
			result := scopes.NewCategorisedResult(category)
			result.Set("type", "track")
			result.Set("id", track.Id)
			result.SetArt(track.Album.Cover)
			result.SetTitle(track.Title)
			result.Set("album", track.Album.Title)
			result.Set("artist", track.Artist.Name)
			var artistAttr [1]map[string]string
			artistAttr[0] = make(map[string]string)
			artistAttr[0]["value"] = track.Artist.Name
			result.Set("artistAttr", artistAttr)
			uriTail := track.Artist.Name + "/" + track.Album.Title + "/" + track.Title
			result.SetURI(track.Link)
			result.SetDndURI("http://localhost_dndrui/" + uriTail)

			err = reply.Push(result)
			if err != nil {
				return err
			}
		}
		var artists, err2 = QueryArtists(queryString)
		if err2 != nil {
			return err2
		}
		category = reply.RegisterCategory("query_artist", "Found artists", "", artistCategoryTemplate)
		for _, artist := range artists {
			result := scopes.NewCategorisedResult(category)
			result.Set("type", "artist")
			result.Set("id", artist.Id)
			result.SetArt(artist.Picture)
			result.SetTitle(artist.Name)
			uriTail := artist.Name
			result.SetURI(artist.Link)
			result.SetDndURI("http://localhost_dndrui/" + uriTail)
			var attributes [2]map[string]string
			attributes[0] = make(map[string]string)
			attributes[0]["value"] = strconv.Itoa(artist.Nb_album)
			attributes[0]["icon"] = s.base.ScopeDirectory() + "/album.svg"
			attributes[1] = make(map[string]string)
			attributes[1]["value"] = strconv.Itoa(artist.Nb_fan)
			attributes[1]["icon"] = "image://theme/starred"
			result.Set("attributes", attributes)
			err = reply.Push(result)
			if err != nil {
				return err
			}
		}
		var albums, err3 = QueryAlbums(queryString)
		if err3 != nil {
			return err3
		}
		category = reply.RegisterCategory("query_album", "Found albums", "", albumCategoryTemplate)
		for _, album := range albums {
			result := scopes.NewCategorisedResult(category)
			result.Set("type", "album")
			result.Set("id", album.Id)
			result.SetArt(album.Cover)
			result.SetTitle(album.Title)
			result.Set("artist", album.Artist.Name)
			uriTail := album.Artist.Name + "/" + album.Title
			result.SetURI(album.Link)
			result.SetDndURI("http://localhost_dndrui/" + uriTail)
			var attributes [1]map[string]string
			attributes[0] = make(map[string]string)
			attributes[0]["value"] = strconv.Itoa(album.Nb_tracks)
			attributes[0]["icon"] = "image://theme/stock_music"
			result.Set("attributes", attributes)
			err = reply.Push(result)
			if err != nil {
				return err
			}
		}
	} else {
		// Surfacing mode
		accessToken := ""
		/*TODO:
		services := s.Accounts.EnabledServices()
		log.Printf("Number of enabled services: %v\n", len(services))
		if len(services) > 0 {
			service := services[0]
			log.Printf("Service: %#v\n", service)
			// If the service is in an error state, try
			// and refresh it.
			if service.Error != nil {
				::service = s.Accounts.Refresh(service.AccountId, false)
				log.Printf("Refreshed: %#v\n", service)
			}
			if service.Error == nil {
				accessToken = service.AccessToken
			}
		}
		*/

		if accessToken != "" {
			var tracks, err = QueryRecommendedTracks(accessToken)
			if err != nil {
				return err
			}
			category := reply.RegisterCategory("query_track", "Found tracks", "", trackCategoryTemplate)
			for _, track := range tracks {
				result := scopes.NewCategorisedResult(category)
				result.Set("type", "track")
				result.Set("id", track.Id)
				result.SetArt(track.Album.Cover)
				result.SetTitle(track.Title)
				result.Set("album", track.Album.Title)
				result.Set("artist", track.Artist.Name)
				var artistAttr [1]map[string]string
				artistAttr[0] = make(map[string]string)
				artistAttr[0]["value"] = track.Artist.Name
				result.Set("artistAttr", artistAttr)
				uriTail := track.Artist.Name + "/" + track.Album.Title + "/" + track.Title
				result.SetURI(track.Link)
				result.SetDndURI("http://localhost_dndrui/" + uriTail)

				err = reply.Push(result)
				if err != nil {
					return err
				}
			}
			var artists, err2 = QueryRecommendedArtists(accessToken)
			if err2 != nil {
				return err2
			}
			category = reply.RegisterCategory("query_artist", "Found artists", "", artistCategoryTemplate)
			for _, artist := range artists {
				result := scopes.NewCategorisedResult(category)
				result.Set("type", "artist")
				result.Set("id", artist.Id)
				result.SetArt(artist.Picture)
				result.SetTitle(artist.Name)
				uriTail := artist.Name
				result.SetURI(artist.Link)
				result.SetDndURI("http://localhost_dndrui/" + uriTail)
				var attributes [2]map[string]string
				attributes[0] = make(map[string]string)
				attributes[0]["value"] = strconv.Itoa(artist.Nb_album)
				attributes[0]["icon"] = s.base.ScopeDirectory() + "/album.svg"
				attributes[1] = make(map[string]string)
				attributes[1]["value"] = strconv.Itoa(artist.Nb_fan)
				attributes[1]["icon"] = "image://theme/starred"
				result.Set("attributes", attributes)
				err = reply.Push(result)
				if err != nil {
					return err
				}
			}
			var albums, err3 = QueryRecommendedAlbums(accessToken)
			if err3 != nil {
				return err3
			}
			category = reply.RegisterCategory("query_album", "Found albums", "", albumCategoryTemplate)
			for _, album := range albums {
				result := scopes.NewCategorisedResult(category)
				result.Set("type", "album")
				result.Set("id", album.Id)
				result.SetArt(album.Cover)
				result.SetTitle(album.Title)
				result.Set("artist", album.Artist.Name)
				uriTail := album.Artist.Name + "/" + album.Title
				result.SetURI(album.Link)
				result.SetDndURI("http://localhost_dndrui/" + uriTail)
				var attributes [1]map[string]string
				attributes[0] = make(map[string]string)
				attributes[0]["value"] = strconv.Itoa(album.Nb_tracks)
				attributes[0]["icon"] = "image://theme/stock_music"
				result.Set("attributes", attributes)
				err = reply.Push(result)
				if err != nil {
					return err
				}
			}
		} else {
			root_department := s.CreateDepartments(query, metadata, reply)
			reply.RegisterDepartments(root_department)

			var genreIdStr string
			if query.DepartmentID() == "" {
				genreIdStr = "0"
			} else {
				genreIdStr = query.DepartmentID()
			}
			// Not logged in, so add nag
			cat := reply.RegisterCategory("", "Nag", "", loginNagTemplate)
			result := scopes.NewCategorisedResult(cat)
			result.SetTitle("Log in to deezer")
			reg_err := scopes.RegisterAccountLoginResult(result, query, "deezer-scope.labsin_deezer-scope", "deezer-scope.labsin_deezer-scope", "deezer-scope.labsin_account", scopes.PostLoginInvalidateResults, scopes.PostLoginDoNothing)
			if reg_err != nil {
				return reg_err
			}

			var rep_err = reply.Push(result)
			if rep_err != nil {
				return rep_err
			}

			var artists, err = GetArtistsFromGenre(genreIdStr)
			if err != nil {
				//err2 := reply.PushSurfacingResultsFromCache()
				//if err2 != nil {
				//	return err2
				//}
				return err
			}
			category := reply.RegisterCategory("artists", "Artists", "", artistFromGenreCategoryTemplate)
			for _, artist := range artists {
				result := scopes.NewCategorisedResult(category)
				result.Set("type", "artist")
				result.Set("id", artist.Id)
				result.SetArt(artist.Picture)
				result.SetTitle(artist.Name)
				uriTail := artist.Name
				result.SetURI(artist.Link)
				result.SetDndURI("http://localhost_dndrui/" + uriTail)
				err = reply.Push(result)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *DeezerScope) CreateDepartments(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply) *scopes.Department {
	department, _ := scopes.NewDepartment("", query, "Browse Music")
	department.SetAlternateLabel("Browse All Music")
	for i, name := range genreName {
		if i == 0 {
			continue
		}
		cur_dep, _ := scopes.NewDepartment(strconv.Itoa(genreIds[i]), query, name)
		department.AddSubdepartment(cur_dep)
	}
	return department
}

// MAIN ************************************************************************

func main() {
	/*log.Println("Deezer: blabla Setting up accounts")
	watcher := accounts.NewWatcher("deezer-scope.labsin_deezer-scope", []string{"deezer-scope.labsin_account"})
	watcher.Settle()*/
	scope := &DeezerScope{
		ClientId: "172955",
		//TODO: Accounts: watcher,
	}
	if err := scopes.Run(scope); err != nil {
		log.Fatalln(err)
	}
}
