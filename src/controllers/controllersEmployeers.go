package controllers

import (
	"check-point/src/config"
	"check-point/src/models"
	"check-point/src/repository"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func Create(w http.ResponseWriter, r *http.Request) {
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

	if err = employee.Prepare(models.Cadastro); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	repository := repository.NewRepositoryEmployee(db)
	createEmployee, err := repository.CreateRepositoryEmployee(ctx, employee)
	if err != nil {
		http.Error(w, "Error creating employee", http.StatusInternalServerError)
		return
	}
	responseJSON, err := json.Marshal(createEmployee)
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	value := strings.ToLower(r.URL.Query().Get("search"))

	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryEmployee(db)

	var listEmployee []models.Employee

	if value == "" {
		listEmployee, err = repository.ListAllEmployee(ctx)
	} else {
		listEmployee, err = repository.ListRepositoryParamsEmployee(ctx, value)
	}

	if err != nil {
		http.Error(w, "Error fetching employee list", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(listEmployee)
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

func Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	employeeID, err := strconv.ParseUint(vars["employeeID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

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

	if err = employee.Prepare(models.Edição); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}

	defer db.Close()
	repository := repository.NewRepositoryEmployee(db)
	if err = repository.UpdateRepositoryEmployee(ctx, employeeID, employee); err != nil {
		http.Error(w, "Error update employee", http.StatusInternalServerError)
		return
	}

	message := "Employee updated successfully"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	employeeID, err := strconv.ParseUint(vars["employeeID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}
	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryEmployee(db)
	if err = repository.DeleteRepositoryEmployee(ctx, employeeID); err != nil {
		http.Error(w, "Error deleted employee", http.StatusInternalServerError)
		return
	}

	message := "Employee deleted successfully"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func ListID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	employeeID, err := strconv.ParseUint(vars["employeeID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryEmployee(db)
	employee, err := repository.ListIDRepositoryEmployee(ctx, employeeID)
	if err != nil {
		http.Error(w, "Error fetching employee", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(employee)
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
