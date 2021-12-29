package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spotifyutils "github.com/ian-flores/digital-music-archive/pkg/spotifyutils"

	"cloud.google.com/go/firestore"
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
	spotifyClient := spotify.New(httpClient)

	defer httpClient.CloseIdleConnections()
	// Bad Bunny Spotify ID: 4q3ewBCX7sLwd24euuV69X

	artistData, err := spotifyutils.GetArtistData(ctx, spotifyClient, "4q3ewBCX7sLwd24euuV69X")
	fmt.Println("Artist name:", artistData.ArtistName)

	firestoreClient := createClient(ctx)

	defer firestoreClient.Close()

	doc_string := fmt.Sprintf("artists/%s", artistData.ArtistID)
	fmt.Println(doc_string)
	firestoreDoc := firestoreClient.Doc(doc_string)

	if _, err := firestoreDoc.Get(ctx); status.Code(err) == codes.NotFound {
		wr, err := firestoreDoc.Set(ctx,
			map[string]interface{}{
				"name":          artistData.ArtistName,
				"followers":     artistData.ArtistFollowers,
				"songs":         artistData.ArtistSongs,
				"collaborators": artistData.ArtistCollaborators,
			},
		)

		if err != nil {
			panic(err)
		}
		fmt.Println(wr.UpdateTime)
	}

}

func createClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "digital-music-archive"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}
