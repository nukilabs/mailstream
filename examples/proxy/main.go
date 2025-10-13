package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/nukilabs/mailstream"
	"github.com/nukilabs/socks"
)

func main() {
	server := socks.NewServer()
	server.Authentication = socks.UserPass("username", "password")

	go server.ListenAndServe("tcp", ":1080")

	config := mailstream.Config{
		Host:     "imap.example.com",
		Port:     993,
		Email:    "mymail@example.com",
		Password: "password1234",
		ProxyURL: &url.URL{Scheme: "socks5h", Host: "localhost:1080", User: url.UserPassword("username", "password")},
	}
	client, err := mailstream.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create a channel to receive mail updates by subscribing to the client
	listener := client.Subscribe()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Tell the client to start listening for mail updates
	done := client.WaitForUpdates(ctx)

	// We run 1 minute, printing mail updates as they come
	for {
		select {
		case mail := <-listener:
			fmt.Println(mail.Subject)
		case err := <-done:
			log.Fatal(err)
		}
	}
}
