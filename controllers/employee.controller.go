package controllers

import (
	db "Hdfc_Assignment/models"
	"Hdfc_Assignment/services"
	"Hdfc_Assignment/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateNewEmployee godoc
// @Summary      Create Employee
// @Description  creates a new employee
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        req  body      models.EmployeeRequest true "Employee Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /employees [post]
// @Security     ApiKeyAuth
func CreateNewEmployee(c *gin.Context) {
	var requestBody utils.EmployeeRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId, exists := c.Get("userId")
	if !exists {
		response.StatusCode = http.StatusBadRequest
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	res, err := services.FindAdminById(userId.(primitive.ObjectID))
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	if res.Role != db.RoleAdmin {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Only Admin Can Callable"
		response.SendResponse(c)
		return
	}

	employee, err := services.CreateEmployee(
		requestBody.Name,
		requestBody.Email,
		requestBody.Password,
		requestBody.Mobile,
		requestBody.Address,
	)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"employee": employee}
	response.SendResponse(c)
}

// GetEmployees godoc
// @Summary      Get Employees
// @Description  gets user employees with pagination
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        page  query    string  false  "Switch page by 'page'"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /employees [get]
// @Security     ApiKeyAuth
func GetEmployees(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	pageQuery := c.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageQuery)
	limit := 5

	employees, _ := services.GetEmployees(userId.(primitive.ObjectID), page, limit)
	hasPrev := page > 0
	hasNext := len(employees) > limit

	if hasNext {
		employees = employees[:limit]
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"employees": employees, "prev": hasPrev, "next": hasNext}
	response.SendResponse(c)
}

// GetOneEmployee godoc
// @Summary      Get a employee
// @Description  get employee by id
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Employee ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /employees/{id} [get]
// @Security     ApiKeyAuth
func GetOneEmployee(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	employeeId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	//employee, err := services.GetEmployeeFromCache(userId.(primitive.ObjectID), employeeId)
	// if err == nil {
	// 	utils.SendResponseData(c, gin.H{"employee": employee, "cache": true})
	// 	return
	// }

	employee, err := services.GetEmployeeById(userId.(primitive.ObjectID), employeeId)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	//	go services.CacheOneEmployee(userId.(primitive.ObjectID), employee)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"employee": employee}
	response.SendResponse(c)
}

// UpdateEmployee godoc
// @Summary      Update a employee
// @Description  updates a employee by id
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id     path    string  true  "Employee ID"
// @Param        req    body    models.EmployeeRequest true "Employee Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /employees/{id} [put]
// @Security     ApiKeyAuth
func UpdateEmployee(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	employeeId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	var employeeRequest utils.EmployeeRequest
	_ = c.ShouldBindBodyWith(&employeeRequest, binding.JSON)

	err := services.UpdateEmployee(userId.(primitive.ObjectID), employeeId, &employeeRequest)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}

// DeleteEmployee godoc
// @Summary      Delete a employee
// @Description  deletes employee by id
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Employee ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /employees/{id} [delete]
// @Security     ApiKeyAuth
func DeleteEmployee(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	employeeId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	err := services.DeleteEmployee(userId.(primitive.ObjectID), employeeId)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}
