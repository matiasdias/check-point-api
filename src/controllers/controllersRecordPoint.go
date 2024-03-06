package controllers

import (
	config "check-point/src/db"
	"check-point/src/models"
	"check-point/src/repository"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// CreateRecordPoint Responsável por criar um novo ponto de registro
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
	if err = point.PreparePoint("cadastro"); err != nil {
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

// ListRecordPoint Responsável por retornar todos os registros
func ListRecordPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	value := strings.ToLower(r.URL.Query().Get("search"))

	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryRecordPoint(db)
	var listRecordPoint []models.RecordWithEmployee

	if value == "" {
		listRecordPoint, err = repository.ListAllRecordPoint(ctx)
	} else {
		listRecordPoint, err = repository.ListRepositoryParamsRecordPoint(ctx, value)
	}

	if err != nil {
		http.Error(w, "Error fetching record Point list", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(listRecordPoint)
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// ListIDRecordPoint Responsável por retornar um registro por id
func ListIDRecordPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	recordID, err := strconv.ParseUint(vars["recordID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryRecordPoint(db)
	existingPoint, err := repository.ListIDRepositoryRecordPoint(ctx, recordID)
	if err != nil {
		http.Error(w, "Error fetching record point", http.StatusInternalServerError)
		return
	}

	if existingPoint.ID == uint64(0) {
		http.Error(w, "Record point not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(existingPoint)
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// DeleteRecordPoint Responsável por deletar um registro de ponto
func DeleteRecordPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	recordID, err := strconv.ParseUint(vars["recordID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryRecordPoint(db)

	// Pesquisa o ponto antes de remover
	existingPoint, err := repository.ListIDRepositoryRecordPoint(ctx, recordID)
	if err != nil {
		http.Error(w, "Error fetching record point", http.StatusInternalServerError)
		return
	}

	// Verifica se esse ponto existe
	if existingPoint.ID == uint64(0) {
		http.Error(w, "Record point not found", http.StatusNotFound)
		return
	}

	if err = repository.DeleteRecordPoint(ctx, recordID); err != nil {
		http.Error(w, "Error deleting record point", http.StatusInternalServerError)
		return
	}

	message := "Recod Point deleted successfully"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	return
}

// UpdateRecordPoint Responsável por atualizar um registro
func UpdateRecordPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	recordID, err := strconv.ParseUint(vars["recordID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	request, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error getting body data", http.StatusUnprocessableEntity)
		return
	}
	var record models.RegisterPointer
	if err = json.Unmarshal(request, &record); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	if err = record.PreparePoint("edição"); err != nil {
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

	existingPoint, err := repository.ListIDRepositoryRecordPoint(ctx, recordID)
	if err != nil {
		http.Error(w, "Error fetching record point", http.StatusInternalServerError)
		return
	}

	if existingPoint.ID == uint64(0) {
		http.Error(w, "Record point not found", http.StatusNotFound)
		return
	}

	if err = repository.UpdateRecordPoint(ctx, recordID, record); err != nil {
		http.Error(w, "Error updating record point", http.StatusInternalServerError)
		return
	}

	message := "record Type updated successfully"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
