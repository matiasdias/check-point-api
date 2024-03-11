package controllers

import (
	config "check-point/src/db"
	"check-point/src/models"
	"check-point/src/repository"
	"check-point/src/utils"
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
		err := utils.InternalServerError("Error getting body data")
		utils.RespondWithError(w, err)
		return
	}

	var point models.RegisterPointer

	if err = json.Unmarshal(request, &point); err != nil {
		err := utils.BadRequestError("Error decoding JSON")
		utils.RespondWithError(w, err)
		return
	}
	if err = point.PreparePoint("cadastro"); err != nil {
		err := utils.BadRequestError(err.Error())
		utils.RespondWithError(w, err)
		return
	}

	db, err := config.Connection()
	if err != nil {
		err := utils.InternalServerError("Error connecting to database")
		utils.RespondWithError(w, err)
		return
	}

	defer db.Close()

	repository := repository.NewRepositoryRecordPoint(db)
	createPoint, err := repository.CreateRecordPoint(ctx, point)
	if err != nil {
		err := utils.InternalServerError("Error creating record point")
		utils.RespondWithError(w, err)
		return
	}
	utils.RespondWithSuccessJSON(w, createPoint)

}

// ListRecordPoint Responsável por retornar todos os registros
func ListRecordPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	value := strings.ToLower(r.URL.Query().Get("search"))

	db, err := config.Connection()
	if err != nil {
		err := utils.InternalServerError("Error connecting to database")
		utils.RespondWithError(w, err)
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
		err := utils.InternalServerError("Error fetching record Point list")
		utils.RespondWithError(w, err)
		return
	}

	utils.RespondWithSuccessJSON(w, listRecordPoint)
}

// ListIDRecordPoint Responsável por retornar um registro por id
func ListIDRecordPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	recordID, err := strconv.ParseUint(vars["recordID"], 10, 64)
	if err != nil {
		err := utils.BadRequestError("Invalid record ID")
		utils.RespondWithError(w, err)
		return
	}

	db, err := config.Connection()
	if err != nil {
		err := utils.InternalServerError("Error connecting to database")
		utils.RespondWithError(w, err)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryRecordPoint(db)
	existingPoint, err := repository.ListIDRepositoryRecordPoint(ctx, recordID)
	if err != nil {
		err := utils.InternalServerError("Error fetching record point")
		utils.RespondWithError(w, err)
		return
	}

	if existingPoint.ID == uint64(0) {
		err := utils.NotFoundError("Record point not found")
		utils.RespondWithError(w, err)
		return
	}

	utils.RespondWithSuccessJSON(w, existingPoint)
}

// DeleteRecordPoint Responsável por deletar um registro de ponto
func DeleteRecordPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	recordID, err := strconv.ParseUint(vars["recordID"], 10, 64)
	if err != nil {
		err := utils.BadRequestError("Invalid record ID")
		utils.RespondWithError(w, err)
		return
	}

	db, err := config.Connection()
	if err != nil {
		err := utils.InternalServerError("Error connecting to database")
		utils.RespondWithError(w, err)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryRecordPoint(db)

	// Pesquisa o ponto antes de remover
	existingPoint, err := repository.ListIDRepositoryRecordPoint(ctx, recordID)
	if err != nil {
		err := utils.InternalServerError("Error fetching record poin")
		utils.RespondWithError(w, err)
		return
	}

	// Verifica se esse ponto existe
	if existingPoint.ID == uint64(0) {
		err := utils.NotFoundError("Record point not found")
		utils.RespondWithError(w, err)
		return
	}

	if err = repository.DeleteRecordPoint(ctx, recordID); err != nil {
		err := utils.InternalServerError("Error deleting record point")
		utils.RespondWithError(w, err)
		return
	}

	utils.RespondWithSuccessMessage(w, "Recod Point deleted successfully")

}

// UpdateRecordPoint Responsável por atualizar um registro
func UpdateRecordPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	recordID, err := strconv.ParseUint(vars["recordID"], 10, 64)
	if err != nil {
		err := utils.BadRequestError("Invalid record ID")
		utils.RespondWithError(w, err)
		return
	}

	request, err := io.ReadAll(r.Body)
	if err != nil {
		err := utils.UnprocessableEntityError("Error getting body data")
		utils.RespondWithError(w, err)
		return
	}
	var record models.RegisterPointer
	if err = json.Unmarshal(request, &record); err != nil {
		err := utils.BadRequestError("Error decoding JSON")
		utils.RespondWithError(w, err)
		return
	}

	if err = record.PreparePoint("edição"); err != nil {
		err := utils.BadRequestError(err.Error())
		utils.RespondWithError(w, err)
		return
	}

	db, err := config.Connection()
	if err != nil {
		err := utils.InternalServerError("Error connecting to database")
		utils.RespondWithError(w, err)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryRecordPoint(db)

	existingPoint, err := repository.ListIDRepositoryRecordPoint(ctx, recordID)
	if err != nil {
		err := utils.InternalServerError("Error fetching record point")
		utils.RespondWithError(w, err)
		return
	}

	if existingPoint.ID == uint64(0) {
		err := utils.NotFoundError("Record point not found")
		utils.RespondWithError(w, err)
		return
	}

	if err = repository.UpdateRecordPoint(ctx, recordID, record); err != nil {
		err := utils.InternalServerError("Error updating record point")
		utils.RespondWithError(w, err)
		return
	}

	utils.RespondWithSuccessMessage(w, "Record Type updated successfully")

}
