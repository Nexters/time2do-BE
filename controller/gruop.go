package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time2do/database"
	"time2do/entity"
)

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var group entity.Group
	json.Unmarshal(requestBody, &group)
	database.Connector.Create(group)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func GetAllGroup(w http.ResponseWriter, r *http.Request) {
	var gruops []entity.Group
	database.Connector.Find(&gruops)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gruops)
}
