package main

import (
	"encoding/json"
	"log"
	"strings"

	astilkit "github.com/asticode/go-astikit"
	astilectron "github.com/asticode/go-astilectron"
)

func main() {
	var a, _ = astilectron.New(nil, astilectron.Options{
		AppName:           "Emoji Uploader",
		BaseDirectoryPath: ".",
	})

	defer a.Close()

	a.HandleSignals()

	a.Start()

	var w *astilectron.Window
	w, _ = a.NewWindow("./index.html", &astilectron.WindowOptions{
		Center: astilkit.BoolPtr(true),
		Height: astilkit.IntPtr(980),
		Width:  astilkit.IntPtr(1200),
	})

	w.Create()

	w.SendMessage("hello", func(m *astilectron.EventMessage) {
		var s string
		m.Unmarshal(&s)

		log.Printf("received %s\n", s)
	})

	initClient()
	var es *EmojiService
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		var iEvent InboundEvent
		m.Unmarshal(&iEvent)
		log.Printf("Message type: %s\n", iEvent.Type)

		// Workaround to recast into event struct types
		jsonBody, _ := json.Marshal(iEvent.Body)
		switch iEvent.Type {
		case "get-emoji-list":
			var eventBody InboundGetEmojiListEvent
			json.Unmarshal(jsonBody, &eventBody)

			es = newEmojiService(eventBody.Space, eventBody.Token)
			r, err := es.retrieveEmojiToUpload(eventBody.Filename)
			if err != nil {
				log.Println(err.Error())
				return err.Error()
			}

			return r
		case "upload-emoji":
			var eventBody InboundEmojiUploadEvent
			json.Unmarshal(jsonBody, &eventBody)

			imageData, err := downloadFile(eventBody.Emoji.URL)
			if err != nil {
				log.Println(err.Error())
				return err.Error()
			}

			splitStr := strings.Split(eventBody.Emoji.URL, ":")
			if splitStr[0] == "alias" {
				err = es.uploadAlias(eventBody.Emoji.Name, splitStr[1])
			} else {
				err = es.uploadEmoji(eventBody.Emoji, imageData)
			}

			if err != nil {
				log.Println(err.Error())
				return err.Error()
			}

			return 200
		}
		return nil
	})

	// Uncomment to enable console
	//w.OpenDevTools()

	a.Wait()
}
