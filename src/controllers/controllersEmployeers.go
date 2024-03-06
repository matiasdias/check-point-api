package controllers

import (
	config "check-point/src/db"
	"check-point/src/models"
	"check-point/src/repository"
	"check-point/src/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Create Responsável por criar um novo funcionário
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
		log.Printf("Error connecting to database: %v\n", err)
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

// List Responsável por listar um novo funcionário
func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	value := strings.ToLower(r.URL.Query().Get("search"))

	db, err := config.Connection()
	if err != nil {
		log.Printf("Error connecting to database: %v\n", err)
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

	if value != "" && len(listEmployee) == 0 {
		http.Error(w, "Employee not found", http.StatusNotFound)
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

// Update Responsável por atualizar um funcionário
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

	// Abre a conexão com o banco
	db, err := config.Connection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryEmployee(db)

	// Pesquisa o funcionario antes de atualizar
	existingEmployee, err := repository.ListIDRepositoryEmployee(ctx, employeeID)
	if err != nil {
		http.Error(w, "Error fetching employee", http.StatusInternalServerError)
		return
	}

	// Verifica se esse funcionario existe
	if existingEmployee.ID == uint64(0) {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	// Verifica se o email existe ou não no banco de dados
	if exists, err := repository.FindByEmailExists(ctx, employee.Email); err != nil {
		http.Error(w, "Error checking email existence", http.StatusInternalServerError)
		return
	} else if exists {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// atualiza o funcionario
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

// Delete Responsável por remover um funcionǽrio
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

	// Pesquisa o funcionario antes de remover
	existingEmployee, err := repository.ListIDRepositoryEmployee(ctx, employeeID)
	if err != nil {
		http.Error(w, "Error fetching employee", http.StatusInternalServerError)
		return
	}

	// Verifica se esse funcionario existe
	if existingEmployee.ID == uint64(0) {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	// Remove o funcionario caso ele exista
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

// ListID Responsável por listar um funcionário por id
func ListID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	employeeID, err := strconv.ParseUint(vars["employeeID"], 10, 64)
	if err != nil {
		log.Printf("Invalid parameter: %v", err)
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

	if employee.ID == uint64(0) {
		http.Error(w, "Employee not found", http.StatusNotFound)
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

// UpdatePassWord Responsável por atualizar a senha do funcionário
func UpdatePassWord(w http.ResponseWriter, r *http.Request) {
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

	var pass models.UpdatePassWord

	if err = json.Unmarshal(request, &pass); err != nil {
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

	// Pesquisa o funcionario antes de remover
	existingEmployee, err := repository.ListIDRepositoryEmployee(ctx, employeeID)
	if err != nil {
		http.Error(w, "Error fetching employee", http.StatusInternalServerError)
		return
	}

	// Verifica se esse funcionario existe
	if existingEmployee.ID == uint64(0) {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	// Retorna a senha salva no banco do funcionairo
	passSaveBank, err := repository.ListByPass(ctx, employeeID)
	if err != nil {
		http.Error(w, "Error search passWord employee", http.StatusInternalServerError)
		return
	}

	// Verifica se a senha salva banco é igual a senha atual
	if err = security.VerifyPassWord(passSaveBank, pass.CurrentPassWord); err != nil {
		http.Error(w, "The password saved in the bank is different from the current password", http.StatusUnauthorized)
		return
	}

	newPassWordHash, err := security.Hash(pass.NewPassWord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = repository.UpdatePassWord(ctx, employeeID, string(newPassWordHash)); err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	message := "Employee updated passWord successfully"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
