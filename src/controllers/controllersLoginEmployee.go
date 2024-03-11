package controllers

import (
	config "check-point/src/db"
	"check-point/src/models"
	"check-point/src/repository"
	"check-point/src/security"
	"check-point/src/token"
	"check-point/src/utils"
	"encoding/json"
	"io"
	"net/http"
)

// Login responsável pelo login do funcionǽrio
func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request, err := io.ReadAll(r.Body)
	if err != nil {
		err := utils.UnprocessableEntityError("Error getting body data")
		utils.RespondWithError(w, err)
		return
	}

	var employee models.Employee
	if err = json.Unmarshal(request, &employee); err != nil {
		err := utils.BadRequestError("Error decoding JSON")
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

	repository := repository.NewRepositoryEmployee(db)
	employeeSaveBank, err := repository.FindByEmail(ctx, employee.Email)
	if err != nil {
		err := utils.InternalServerError("Error fetching employee")
		utils.RespondWithError(w, err)
		return
	}

	if err = security.VerifyPassWord(employeeSaveBank.PassWord, employee.PassWord); err != nil {
		err := utils.UnauthorizedError("Error when verifying the user's password saved in the bank")
		utils.RespondWithError(w, err)
		return
	}

	token, err := token.CreateToken(employeeSaveBank.ID)
	if err != nil {
		err := utils.InternalServerError("Error creating to token")
		utils.RespondWithError(w, err)
		return
	}

	responseToken := token
	utils.RespondWithSuccessJSON(w, map[string]string{"Token": responseToken})

}
