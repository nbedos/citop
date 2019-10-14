package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_RecentRepoBuilds(t *testing.T) {
	token := os.Getenv("TRAVIS_API_TOKEN")
	if token == "" {
		t.Fatal("Environment variable TRAVIS_API_TOKEN is not set")
	}

	authUrl := "https://api.travis-ci.org/pusher/auth"
	wsURL := PusherUrl("ws.pusherapp.com", "5df8ac576dcccf4fd076")
	//channel := "private-user-1842548"
	channel := "repo-25564643"

	authHeader := map[string]string{
		"Authorization": fmt.Sprintf("token %s", token),
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	p, err := NewPusherClient(ctx, wsURL, authUrl, authHeader)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

forever:
	for {
		event, err := p.NextEvent(ctx)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			break forever
		case err != nil:
			t.Fatal(err)
		}
		fmt.Printf("received event %s (%s)\n", event.Event, string(event.Data))
		switch event.Event {
		case ConnectionEstablished:
			err := p.Subscribe(ctx, channel)
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				break forever
			case err != nil:
				t.Fatal(err)
			}
		}
	}
}