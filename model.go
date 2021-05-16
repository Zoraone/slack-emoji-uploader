package main

type Emoji struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type InboundEvent struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type InboundGetEmojiListEvent struct {
	Token    string `json:"token"`
	Filename string `json:"filename"`
	Space    string `json:"space"`
}

type InboundEmojiUploadEvent struct {
	Emoji Emoji `json:"emoji"`
}

type OutboundEvent struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}
