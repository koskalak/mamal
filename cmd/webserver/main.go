package main

import (
	"context"
	"log"

	"github.com/spotify-bot/server/internal/config"
	"github.com/spotify-bot/server/internal/spotify"
	"github.com/spotify-bot/server/internal/webserver"
	"golang.org/x/oauth2"
	spotifyOauth "golang.org/x/oauth2/spotify"
)

func main() {
	ctx := context.Background()

	authConf := &oauth2.Config{
		ClientID:     config.AppConfig.Spotify.SpotifyClientID,
		ClientSecret: config.AppConfig.Spotify.SpotifyClientSecret,
		Scopes: []string{
			"user-read-currently-playing",
			"user-read-recently-played",
			"user-modify-playback-state",
		},
		Endpoint:    spotifyOauth.Endpoint,
		RedirectURL: "http://" + config.AppConfig.Spotify.ApiServerAddress + "/auth/callback", //FIXME
	}

	s, err := spotify.New(ctx, spotify.ProviderOptions{
		DatabaseDSN: config.AppConfig.Webserver.MongoDSN,
		AuthConfig:  authConf,
	})
	if err != nil {
		log.Fatal("Failed to connect to Mongo Storage", err)
	}

	w := webserver.New(webserver.WebServerOptions{
		Spotify: s,
	})

	err = w.Start(config.AppConfig.Webserver.Address)

	if err != nil {
		log.Fatalln("Cannot start server", err)
	}
}
