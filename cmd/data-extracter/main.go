package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"

	spotifyutils "github.com/ian-flores/puerto-rico-digital-music-archive/pkg/spotifyutils"

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

	artistData, err := spotifyutils.GetArtistData(ctx, client, "4q3ewBCX7sLwd24euuV69X")
	fmt.Println("Artist name:", artistData.ArtistName)
	fmt.Println("Followers:", artistData.ArtistFollowers)
	fmt.Println("Artist Collaborators:", artistData.ArtistCollaborators)

}
