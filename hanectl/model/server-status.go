package model

type ServerStatus struct {
	Websocket bool `json:"websocket"`
	Mqtt      int  `json:"mqtt"`
}
