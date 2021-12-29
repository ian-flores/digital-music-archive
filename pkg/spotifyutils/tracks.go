package spotifyutils

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

func getArtistTracks(ctx context.Context, client *spotify.Client, artistID spotify.ID) ([]spotify.ID, error) {
	albums, err := client.GetArtistAlbums(ctx, artistID, nil)
	if err != nil {
		return nil, err
	}

	var tracks []spotify.ID
	for _, album := range albums.Albums {
		tracks, err = client.getAlbumTracks(ctx, album.ID, spotify.Market("US"))
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, tracks...)
	}

	return tracks, nil
}
