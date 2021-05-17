package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type Env struct {
	Space string
	Token string
}

func (e *Env) isMissingConfig() bool {
	if e.Space == "" && e.Token == "" {
		return true
	}
	return false
}

func loadEnvVariables() (Env, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not found, using set environment variables.")
	}

	e := Env{
		Space: os.Getenv("SLACK_SPACE"),
		Token: os.Getenv("SLACK_TOKEN"),
	}

	if e.isMissingConfig() {
		return Env{}, fmt.Errorf("Missing environment variables, make sure they have been set correctly!")
	}
	return e, nil
}

func TestGetEmojiList(t *testing.T) {
	env, err := loadEnvVariables()
	if err != nil {
		t.Error(err)
	}
	es := newEmojiService(env.Space, env.Token)

	resp, err := es.getEmojiList()
	if err != nil {
		t.Error(err)
	}

	if resp == nil {
		t.Errorf("Received %v from emoji.list response, expected emoji list array.", resp)
	}
}

func TestUploadEmoji(t *testing.T) {
	env, err := loadEnvVariables()
	if err != nil {
		t.Error(err)
	}
	es := newEmojiService(env.Space, env.Token)

	emoji := Emoji{
		Name: "feelsbadman",
		Type: "image",
		URL:  "https://emoji.slack-edge.com/T02G0V63S/feelsbadman/06040f4acb5c61d4.png",
	}

	imageData, err := downloadFile(emoji.URL)
	if err != nil {
		t.Error(err)
	}

	err = es.uploadEmoji(emoji, imageData)
	if err != nil {
		t.Error(err)
	}

	// Cleanup
	err = es.removeEmoji(emoji.Name)
	if err != nil {
		t.Error(err)
	}
}
