package controller

import (
	"encoding/json"
	"net/http"
	"time2do/database"
	"time2do/entity"
)

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	var persons []entity.User
	database.Connector.Find(&persons)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(persons)
}
