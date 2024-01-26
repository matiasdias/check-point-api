package controllers

import (
	"check-point/src/config"
	"check-point/src/models"
	"check-point/src/repository"
	"encoding/json"
	"io"
	"net/http"
)

func CreateRecordPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error getting body data", http.StatusUnprocessableEntity)
		return
	}

	var point models.RegisterPointer

	if err = json.Unmarshal(request, &point); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	if err = point.PreparePoint(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	repository := repository.NewRepositoryRecordPoint(db)
	createPoint, err := repository.CreateRecordPoint(ctx, point)
	if err != nil {
		http.Error(w, "Error creating record point", http.StatusInternalServerError)
		return
	}
	responseJSON, err := json.Marshal(createPoint)
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}
