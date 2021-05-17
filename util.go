package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

/*
 * Wrapper around client to make testing easier
 * by mocking an http client making calls to a
 * server elsewhere
 */
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

// Default implementation allows the app to use the standard
// http client
func initClient() {
	Client = &http.Client{}
}

func buildDiffList(newList, existingList map[string]interface{}) []Emoji {
	var emojiList []Emoji

	for k, v := range newList {
		if existingList[k] == nil {
			UrlString := fmt.Sprintf("%v", v)
			subs := strings.Split(UrlString, "%v")
			if subs[0] == "alias" {
				emojiList = append(emojiList, Emoji{
					Name: k,
					Type: "alias",
					URL:  subs[1],
				})
			} else {
				emojiList = append(emojiList, Emoji{
					Name: k,
					Type: "image",
					URL:  UrlString,
				})
			}
		}
	}

	return emojiList
}

func parseFile(filename string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	json.Unmarshal(jsonBytes, &result)

	return result["emoji"].(map[string]interface{}), nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}
