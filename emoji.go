package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strconv"
	"time"
)

type EmojiService struct {
	UrlEmojiList string
	UrlEmojiAdd  string

	Token string
}

func newEmojiService(space, token string) *EmojiService {
	return &EmojiService{
		UrlEmojiList: fmt.Sprintf("https://%s.slack.com/api/emoji.list", space),
		UrlEmojiAdd:  fmt.Sprintf("https://%s.slack.com/api/emoji.add", space),
		Token:        token,
	}
}

func (es *EmojiService) retrieveEmojiToUpload(filename string) ([]Emoji, error) {
	newEmojiJSON, err := parseFile(filename)
	if err != nil {
		return nil, err
	}

	existingEmojiJSON, err := es.getEmojiList()
	if err != nil {
		return nil, err
	}

	diffEmojiList := buildDiffList(newEmojiJSON, existingEmojiJSON)
	if err != nil {
		return nil, err
	}

	return diffEmojiList, nil
}

func (es *EmojiService) getEmojiList() (map[string]interface{}, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	mode, err := w.CreateFormField("token")
	if err != nil {
		return nil, err
	}
	mode.Write([]byte(es.Token))

	defer w.Close()

	req, err := http.NewRequest("GET", es.UrlEmojiList, buf)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("token", es.Token)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseJson map[string]interface{}
	json.Unmarshal(content, &responseJson)

	return responseJson["emoji"].(map[string]interface{}), nil
}

func (es *EmojiService) uploadEmoji(emoji Emoji, imageData []byte) error {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	mode, err := w.CreateFormField("mode")
	if err != nil {
		return err
	}
	mode.Write([]byte("data"))

	nameWriter, err := w.CreateFormField("name")
	nameWriter.Write([]byte(emoji.Name))
	tokenWriter, err := w.CreateFormField("token")
	tokenWriter.Write([]byte(es.Token))

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", filepath.Base(emoji.URL)))
	h.Set("Content-Type", fmt.Sprintf("image/%s", filepath.Ext(emoji.URL)[1:]))

	image, err := w.CreatePart(h)
	if err != nil {
		return err
	}
	image.Write(imageData)

	w.Close()

	retry := true
	for retry {
		retry, err = es.upload(emoji.Name, w, *buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *EmojiService) uploadAlias(name, aliasFor string) error {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	mode, err := w.CreateFormField("mode")
	if err != nil {
		return err
	}
	mode.Write([]byte("alias"))

	nameWriter, err := w.CreateFormField("name")
	nameWriter.Write([]byte(name))
	tokenWriter, err := w.CreateFormField("token")
	tokenWriter.Write([]byte(es.Token))
	aliasWriter, err := w.CreateFormField("alias_for")
	aliasWriter.Write([]byte(aliasFor))

	w.Close()

	retry := true
	for retry {
		retry, err = es.upload(name, w, *buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *EmojiService) upload(name string, w *multipart.Writer, buf bytes.Buffer) (bool, error) {
	req, err := http.NewRequest("POST", es.UrlEmojiAdd, &buf)
	if err != nil {
		return false, err
	}

	req.SetBasicAuth("api", es.Token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		break
	case 429:
		wait, err := strconv.Atoi(resp.Header.Get("retry-after"))
		if err != nil {
			return false, err
		}
		fmt.Printf("429 Too many requests, sleeping for %d seconds %s\n", wait, name)
		time.Sleep(time.Duration(wait) * time.Second)
		return true, nil
	default:
		return false, fmt.Errorf("Unexpected status code: %v for %s", resp.StatusCode, name)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var responseJson map[string]interface{}
	json.Unmarshal(content, &responseJson)

	if responseJson["ok"] != true {
		if responseJson["error"] == "not_authed" {

		}
		return false, fmt.Errorf("Error uploading emoji %s %s", name, responseJson["error"])
	}

	fmt.Printf("Success uploading %s\n", name)
	return false, nil
}
