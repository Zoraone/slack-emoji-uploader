package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

/*
 * Mock client for tests
 * Creates a custom version of the Do request which is called
 * instead of the standard http one, allowing for a custom
 * response.
 */
type MockDoFunc func(req *http.Request) (*http.Response, error)
type MockClient struct {
	MockDo MockDoFunc
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func TestGetEmojiList(t *testing.T) {
	es := newEmojiService("test_space", "test_token")
	jsonResponse := `{
		"ok": true,
		"emoji": {
			"bowtie": "https://emoji.slack-edge.com/T02G0V63S/bowtie/f3ec6f2bb0.png"
		}	
	}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonResponse)))

	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}

	resp, err := es.getEmojiList()
	if err != nil {
		t.Error(err)
	}

	if resp == nil {
		t.Errorf("Received %v from emoji.list response, expected emoji list array.", resp)
	}
}

func TestUploadEmoji(t *testing.T) {
	es := newEmojiService("test_space", "test_token")
	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"ok":true}`))),
			}, nil
		},
	}

	emoji := Emoji{
		Name: "feelsbadman",
		Type: "image",
		URL:  "https://emoji.slack-edge.com/T02G0V63S/feelsbadman/06040f4acb5c61d4.png",
	}

	err := es.uploadEmoji(emoji, []byte("this is image"))
	if err != nil {
		t.Error(err)
	}
}
