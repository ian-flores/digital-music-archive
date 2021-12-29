package spotifyutils

import (
	"context"
	"fmt"
	"os"

	"github.com/zmb3/spotify/v2"
)

type ArtistData struct {
	ArtistID            spotify.ID
	ArtistName          string
	ArtistFollowers     uint
	ArtistSongs         []ArtistTrack
	ArtistCollaborators map[string]ArtistCollaborator
}

type ArtistTrack struct {
	TrackID   spotify.ID
	TrackName string
	Artists   map[string]ArtistCollaborator
}

type ArtistCollaborator struct {
	ArtistID   spotify.ID
	ArtistName string
}

// GetArtistData returns the data for an artist
// It accepts the context of the app, the Spotify client, and the artist ID
// It returns the ArtistData struct
func GetArtistData(ctx context.Context, client *spotify.Client, artistID string) (ArtistData, error) {

	artist, err := client.GetArtist(ctx, spotify.ID(artistID))

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return ArtistData{}, err
	}

	var artistData ArtistData

	artistData.ArtistID = artist.ID
	artistData.ArtistName = artist.Name
	artistData.ArtistFollowers = artist.Followers.Count

	albums, err := client.GetArtistAlbums(ctx, artist.ID, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return ArtistData{}, err
	}

	all_artists := make(map[string]ArtistCollaborator)

	for _, album := range albums.Albums {
		tracks, _ := client.GetAlbumTracks(ctx, album.ID, spotify.Market("US"))
		for _, track := range tracks.Tracks {
			artists := getArtistsInTrack(track)
			artistData.ArtistSongs = append(artistData.ArtistSongs,
				ArtistTrack{
					track.ID,
					track.Name,
					artists,
				})

			for _, artist := range artists {
				if _, ok := all_artists[string(artist.ArtistID)]; !ok {
					all_artists[string(artist.ArtistID)] = artist
				}
			}

		}
	}

	artistData.ArtistCollaborators = all_artists

	return artistData, nil
}

func getArtistsInTrack(track spotify.SimpleTrack) map[string]ArtistCollaborator {
	artists := make(map[string]ArtistCollaborator)

	for _, artist := range track.Artists {
		if _, ok := artists[string(artist.ID)]; !ok {
			artists[string(artist.ID)] = ArtistCollaborator{
				ArtistID:   artist.ID,
				ArtistName: artist.Name,
			}
		}
	}
	return artists
}
