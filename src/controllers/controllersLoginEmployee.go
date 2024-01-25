package controllers

import (
	"check-point/src/auth"
	"check-point/src/config"
	"check-point/src/models"
	"check-point/src/repository"
	"check-point/src/security"
	"encoding/json"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error getting body data", http.StatusUnprocessableEntity)
		return
	}

	var employee models.Employee
	if err = json.Unmarshal(request, &employee); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryEmployee(db)
	employeeSaveBank, err := repository.FindByEmail(ctx, employee.Email)
	if err != nil {
		http.Error(w, "Error fetching employee", http.StatusInternalServerError)
		return
	}

	if err = security.VerifyPassWord(employeeSaveBank.PassWord, employee.PassWord); err != nil {
		http.Error(w, "Error when verifying the user's password saved in the bank", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateToken(employeeSaveBank.ID, employeeSaveBank.Email)
	if err != nil {
		http.Error(w, "Error creating to token", http.StatusInternalServerError)
		return
	}

	reponseToken := token
	responseJSON, err := json.Marshal(map[string]string{"token:": reponseToken})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}