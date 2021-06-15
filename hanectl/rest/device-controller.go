package rest

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"net/http"
	"strings"
)

func deviceAuthoritiesFilter(authorities []string) model.DeviceFilter {
	return func(device *model.Device) bool {
		if len(device.Authorities) > 0 {
			for _, ua := range device.Authorities {
				for _, a := range authorities {
					if strings.Compare(ua, a) == 0 {
						return true
					}
				}
			}
			return false
		}
		return true
	}
}

func deviceConfig(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory) {
	if _, authorities, err := getCurrentUser(r); err == nil {
		if devices, err := serviceFactory.GetDeviceService().GetDevicesDto(deviceAuthoritiesFilter(authorities)); err == nil {
			render.JSON(w, r, devices)
			return
		}
	}
	http.Error(w, "Cant get device configuration", http.StatusInternalServerError)
}

func deviceStates(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory) {
	if state := serviceFactory.GetDeviceService().DeviceStates(); state == nil {
		http.Error(w, "Cant get device state", http.StatusNotFound)
	} else {
		render.JSON(w, r, state)
	}
}

func deviceState(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory) {
	deviceKey := chi.URLParam(r, id)

	if state := serviceFactory.GetDeviceService().DeviceState(deviceKey); state == nil {
		http.Error(w, "Cant get device state", http.StatusNotFound)
	} else {
		render.JSON(w, r, state)
	}
}

func getDeviceForRequest(deviceKey string, r *http.Request, serviceFactory model.IServiceFactory) *model.Device {
	if _, authorities, err := getCurrentUser(r); err == nil {
		if device, err := serviceFactory.GetDeviceService().GetDevice(deviceKey); err == nil {
			if deviceAuthoritiesFilter(authorities)(device) {
				return device
			}
		}
	}
	return nil
}

func deviceCommand(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory) {
	payload := make(model.Dictionary)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Cant decode payload", http.StatusBadRequest)
	} else {
		var deviceKey string

		if value, exists := payload["device"]; exists {
			deviceKey = value.(string)
		} else {
			http.Error(w, "op not supported without device", http.StatusBadRequest)
			return
		}

		if device := getDeviceForRequest(deviceKey, r, serviceFactory); device != nil {
			setDeviceCommand(w, r, serviceFactory, deviceKey, payload)
		} else {
			http.Error(w, "access forbidden ", http.StatusForbidden)
		}
	}
}

func command(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory, cfg config.IConfiguration) {
	cfgToken := cfg.GetStr(config.CommandToken, "")
	if cfgToken != "" {
		token := chi.URLParam(r, token)
		if strings.Compare(cfgToken, token) == 0 {
			deviceKey := chi.URLParam(r, id)
			payload := make(model.Dictionary)
			payload["device"] = deviceKey
			for key, value := range r.URL.Query() {
				payload[key] = value[0]
			}
			setDeviceCommand(w, r, serviceFactory, deviceKey, payload)
		} else {
			http.Error(w, "access forbidden ", http.StatusForbidden)
		}
	} else {
		http.Error(w, "access forbidden ", http.StatusForbidden)
	}
}

func groupCommand(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory, cfg config.IConfiguration) {
	cfgToken := cfg.GetStr(config.CommandToken, "")
	if cfgToken != "" {
		token := chi.URLParam(r, token)
		if strings.Compare(cfgToken, token) == 0 {
			groupKey := chi.URLParam(r, id)
			payload := make(model.Dictionary)
			for key, value := range r.URL.Query() {
				payload[key] = value[0]
			}
			setDeviceGroupCommand(w, r, serviceFactory, groupKey, payload)
		} else {
			http.Error(w, "access forbidden ", http.StatusForbidden)
		}
	} else {
		http.Error(w, "access forbidden ", http.StatusForbidden)
	}
}

func setDeviceCommand(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory, deviceKey string, payload model.Dictionary) {
	state := serviceFactory.GetDeviceService().DeviceCommand(deviceKey, payload)
	if state {
		render.JSON(w, r, state)
	} else {
		http.Error(w, "Cant execute device command", http.StatusNotFound)
	}
}


func setDeviceGroupCommand(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory, groupKey string, payload model.Dictionary) {
	state := serviceFactory.GetDeviceService().DeviceGroupCommand(groupKey, payload)
	if state {
		render.JSON(w, r, state)
	} else {
		http.Error(w, "Cant execute device command", http.StatusNotFound)
	}
}
