package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func main() {
	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	// Bad Bunny Spotify ID: 4q3ewBCX7sLwd24euuV69X
	artist, err := client.GetArtist(ctx, spotify.ID("4q3ewBCX7sLwd24euuV69X"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Println("Artist name:", artist.Name)
	fmt.Println("Followers:", artist.Followers.Count)

	albums, err := client.GetArtistAlbums(ctx, artist.ID, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	for _, album := range albums.Albums {
		fmt.Println(album.Name)

		tracks, err := client.GetAlbumTracks(ctx, album.ID, spotify.Market("US"))

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		for _, track := range tracks.Tracks {
			fmt.Println("\t", track.ID, ":", track.Name)
		}

	}

}
