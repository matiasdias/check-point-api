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
		err := utils.UnprocessableEntityError("Error getting body data")
		utils.RespondWithError(w, err)
		return
	}

	var employee models.Employee

	employeeIDAdmin, err := token.ExtractEmployeeIDAdmin(r)
	if err != nil {
		err := utils.BadRequestError("Failed to process is_admin token")
		utils.RespondWithError(w, err)
		return
	}

	if employeeIDAdmin == false {
		err := utils.ForbiddenError("Only admins can create new employees")
		utils.RespondWithError(w, err)
		return
	}

	if err = json.Unmarshal(request, &employee); err != nil {
		err := utils.BadRequestError("Error decoding JSON")
		utils.RespondWithError(w, err)
		return
	}

	if err = employee.Prepare(models.Cadastro); err != nil {
		err := utils.BadRequestError(err.Error())
		utils.RespondWithError(w, err)
		return
	}

	db, err := config.Connection()
	if err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		err := utils.InternalServerError("Error connecting to database")
		utils.RespondWithError(w, err)
		return
	}

	defer db.Close()

	repository := repository.NewRepositoryEmployee(db)
	createEmployee, err := repository.CreateRepositoryEmployee(ctx, employee)
	if err != nil {
		err := utils.InternalServerError("Error creating employee")
		utils.RespondWithError(w, err)
		return
	}

	utils.RespondWithSuccessJSON(w, createEmployee)
}

// List Responsável por listar um novo funcionário
func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	value := strings.ToLower(r.URL.Query().Get("search"))
	employeeIDAdmin, err := token.ExtractEmployeeIDAdmin(r)
	if err != nil {
		err := utils.BadRequestError("Failed to process is_admin token")
		utils.RespondWithError(w, err)
		return
	}

	db, err := config.Connection()
	if err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		err := utils.InternalServerError("Error connecting to database")
		utils.RespondWithError(w, err)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryEmployee(db)

	var listEmployee []models.Employee

	if value == "" {
		listEmployee, err = repository.ListAllEmployee(ctx, employeeIDAdmin)
	} else {
		listEmployee, err = repository.ListRepositoryParamsEmployee(ctx, value, employeeIDAdmin)
	}
	if err != nil {
		err := utils.InternalServerError("Error fetching employee list")
		utils.RespondWithError(w, err)
		return
	}

	if value != "" && len(listEmployee) == 0 {
		err := utils.NotFoundError("Employee not found")
		utils.RespondWithError(w, err)
		return
	}

	utils.RespondWithSuccessJSON(w, listEmployee)

}

// Update Responsável por atualizar um funcionário
func Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	employeeID, err := strconv.ParseUint(vars["employeeID"], 10, 64)
	if err != nil {
		err := utils.BadRequestError("Invalid employee ID")
		utils.RespondWithError(w, err)
		return
	}

	employeeIDAdmin, err := token.ExtractEmployeeIDAdmin(r)
	if err != nil {
		err := utils.BadRequestError("Failed to process is_admin token")
		utils.RespondWithError(w, err)
		return
	}

	if employeeIDAdmin == false {
		err := utils.ForbiddenError("Only admins can update employees")
		utils.RespondWithError(w, err)
		return
	}

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

	if err = employee.Prepare(models.Edição); err != nil {
		err := utils.BadRequestError(err.Error())
		utils.RespondWithError(w, err)
		return
	}

	// Abre a conexão com o banco
	db, err := config.Connection()
	if err != nil {
		err := utils.InternalServerError("Error connecting to database")
		utils.RespondWithError(w, err)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryEmployee(db)

	// Pesquisa o funcionario antes de atualizar
	existingEmployee, err := repository.ListIDRepositoryEmployee(ctx, employeeID)
	if err != nil {
		err := utils.InternalServerError("Error fetching employee")
		utils.RespondWithError(w, err)
		return
	}

	// Verifica se esse funcionario existe
	if existingEmployee.ID == uint64(0) {
		err := utils.NotFoundError("Employee not found")
		utils.RespondWithError(w, err)
		return
	}

	if employee.Email != existingEmployee.Email {
		// Verifica se o email existe ou não no banco de dados
		exists, err := repository.FindByEmailExists(ctx, employee.Email)
		if err != nil {
			err := utils.InternalServerError("Error checking email existence")
			utils.RespondWithError(w, err)
			return
		}
		if exists {
			err := utils.ConflitError("Email already exists")
			utils.RespondWithError(w, err)
			return
		}
	}

	// atualiza o funcionario
	if err = repository.UpdateRepositoryEmployee(ctx, employeeID, employee); err != nil {
		err := utils.InternalServerError("Error update employee")
		utils.RespondWithError(w, err)
		return
	}

	utils.RespondWithSuccessMessage(w, "Employee updated successfully")

}

// Delete Responsável por remover um funcionǽrio
func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	employeeID, err := strconv.ParseUint(vars["employeeID"], 10, 64)
	if err != nil {
		err := utils.BadRequestError("Invalid employee ID")
		utils.RespondWithError(w, err)
		return
	}

	employeeIDAdmin, err := token.ExtractEmployeeIDAdmin(r)
	if err != nil {
		err := utils.BadRequestError("Failed to process is_admin token")
		utils.RespondWithError(w, err)
		return
	}

	if employeeIDAdmin == false {
		err := utils.ForbiddenError("Only admins can deleted employees")
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

	// Pesquisa o funcionario antes de remover
	existingEmployee, err := repository.ListIDRepositoryEmployee(ctx, employeeID)
	if err != nil {
		err := utils.InternalServerError("Error fetching employee")
		utils.RespondWithError(w, err)
		return
	}

	// Verifica se esse funcionario existe
	if existingEmployee.ID == uint64(0) {
		err := utils.NotFoundError("Employee not found")
		utils.RespondWithError(w, err)
		return
	}

	// Remove o funcionario caso ele exista
	if err = repository.DeleteRepositoryEmployee(ctx, employeeID); err != nil {
		err := utils.InternalServerError("Error deleted employee")
		utils.RespondWithError(w, err)
		return
	}

	utils.RespondWithSuccessMessage(w, "Employee deleted successfully")
}

// ListID Responsável por listar um funcionário por id
func ListID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	employeeID, err := strconv.ParseUint(vars["employeeID"], 10, 64)
	if err != nil {
		log.Printf("Invalid parameter: %v", err)
		err := utils.BadRequestError("Invalid employee ID")
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

	employee, err := repository.ListIDRepositoryEmployee(ctx, employeeID)
	if err != nil {
		err := utils.InternalServerError("Error fetching employee")
		utils.RespondWithError(w, err)
		return
	}

	if employee.ID == uint64(0) {
		err := utils.NotFoundError("Employee not found")
		utils.RespondWithError(w, err)
		return
	}

	utils.RespondWithSuccessJSON(w, employee)

}

// UpdatePassWord Responsável por atualizar a senha do funcionário
func UpdatePassWord(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	employeeID, err := strconv.ParseUint(vars["employeeID"], 10, 64)
	if err != nil {
		err := utils.BadRequestError("Invalid employee ID")
		utils.RespondWithError(w, err)
		return
	}
	employeeIDAdmin, err := token.ExtractEmployeeIDAdmin(r)
	if err != nil {
		err := utils.BadRequestError("Failed to process is_admin token")
		utils.RespondWithError(w, err)
		return
	}

	if employeeIDAdmin == false {
		err := utils.ForbiddenError("Only administrators can update other employees passwords")
		utils.RespondWithError(w, err)
		return
	}

	request, err := io.ReadAll(r.Body)
	if err != nil {
		err := utils.UnprocessableEntityError("Error getting body data")
		utils.RespondWithError(w, err)
		return
	}

	var pass models.UpdatePassWord

	if err = json.Unmarshal(request, &pass); err != nil {
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

	existingEmployee, err := repository.ListIDRepositoryEmployee(ctx, employeeID)
	if err != nil {
		err := utils.InternalServerError("Error fetching employee")
		utils.RespondWithError(w, err)
		return
	}

	// Verifica se esse funcionário existe
	if existingEmployee.ID == uint64(0) {
		err := utils.NotFoundError("Employee not found")
		utils.RespondWithError(w, err)
		return
	}

	newPassWordHash, err := security.Hash(pass.NewPassWord)
	if err != nil {
		err := utils.BadRequestError(err.Error())
		utils.RespondWithError(w, err)
		return
	}

	if err = repository.UpdatePassWord(ctx, employeeID, string(newPassWordHash)); err != nil {
		err := utils.InternalServerError("Error updating password")
		utils.RespondWithError(w, err)
		return
	}

	utils.RespondWithSuccessMessage(w, "Employee updated passWord successfully")

}
