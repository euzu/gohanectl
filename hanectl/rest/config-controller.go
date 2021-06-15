package rest

import (
	"encoding/json"
	"github.com/go-chi/render"
	"gohanectl/hanectl/model"
	"net/http"
)

func serverStatus(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory) {
	if status, err := serviceFactory.GetConfigService().GetServerStatus(); err == nil {
		render.JSON(w, r, status)
		return
	}
	http.Error(w, "Cant get server status", http.StatusInternalServerError)
}

func configRooms(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory) {
	if rooms, err := serviceFactory.GetConfigService().GetRoomsConfig(); err == nil {
		render.JSON(w, r, rooms)
		return
	}
	http.Error(w, "Cant get server status", http.StatusInternalServerError)
}

func setUserSettings(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory) {
	if uname, _, err := getCurrentUser(r); err == nil {
		payload := model.UserSettings{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Cant decode payload", http.StatusBadRequest)
			return
		} else {
			if err := serviceFactory.GetUserService().SaveSettings(uname, &payload); err != nil {
				http.Error(w, "Failed to save ", http.StatusInternalServerError)
				return
			}
		}
		render.JSON(w, r, "success")
	} else {
		http.Error(w, "Cant get user", http.StatusUnauthorized)
	}
}

func getUserSettings(w http.ResponseWriter, r *http.Request, serviceFactory model.IServiceFactory) {
	if uname, _, err := getCurrentUser(r); err == nil {
		if settings, err := serviceFactory.GetUserService().GetSettings(uname); err == nil {
			render.JSON(w, r, settings)
		} else {
			http.Error(w, "Cant get user settings", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Cant get user", http.StatusUnauthorized)
	}
}
