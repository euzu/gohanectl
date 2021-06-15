package service

import (
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/model"
	"io/ioutil"
	"net/http"
)

type RestService struct {
	messageHandler model.RestMessageHandler
}

func (r *RestService) GetRequest(url string, device *model.Device) bool {
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			if bodyBytes, err := ioutil.ReadAll(resp.Body); err == nil {
				payload := string(bodyBytes)
				r.HandleMessage(device, payload)
				return true
			} else {
				log.Error().Msgf("Failed to parse response: %v", err)
			}
		} else {
			log.Error().Msgf("Failed request with status: %d", resp.StatusCode)
		}
	} else {
		log.Debug().Msgf("Failed request: %v", err)
	}
	return false
}

func (r *RestService) SetMessageHandler(handler model.RestMessageHandler) {
	r.messageHandler = handler
}

func (r *RestService) HandleMessage(device *model.Device, payload string) {
	if r.messageHandler != nil {
		r.messageHandler(device, payload)
	}
}

func NewRestService() model.IRestService {
	return new(RestService)
}
