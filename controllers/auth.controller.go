package controllers

import (
	// "errors"
	db "Hdfc_Assignment/models"
	"Hdfc_Assignment/services"
	"Hdfc_Assignment/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary      Register
// @Description  registers a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.RegisterRequest true "Register Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/register [post]
func UserRegister(c *gin.Context) {
	var requestBody utils.RegisterRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// is email in use
	err := services.CheckUserMail(requestBody.Email, db.RoleUser)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// create user record
	requestBody.Name = strings.TrimSpace(requestBody.Name)
	user, err := services.CreateUser(requestBody.Name, requestBody.Email, requestBody.Password)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// generate access tokens
	accessToken, refreshToken, err := services.GenerateTokens(user.Email, user.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{
		"user": user,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}

func AdminRegister(c *gin.Context) {
	var requestBody utils.RegisterRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// is email in use
	err := services.CheckUserMail(requestBody.Email, db.RoleAdmin)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// create user record
	requestBody.Name = strings.TrimSpace(requestBody.Name)
	admin, err := services.CreateAdmin(requestBody.Name, requestBody.Email, requestBody.Password)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// generate access tokens
	accessToken, refreshToken, err := services.GenerateTokens(admin.Email, admin.ID)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{
		"admin": admin,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}

func EmployeeRegister(c *gin.Context) {
	var requestBody utils.RegisterRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// is email in use
	err := services.CheckUserMail(requestBody.Email, db.RoleEmployee)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// create user record
	requestBody.Name = strings.TrimSpace(requestBody.Name)
	user, err := services.CreateUser(requestBody.Name, requestBody.Email, requestBody.Password)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// generate access tokens
	accessToken, refreshToken, err := services.GenerateTokens(user.Email, user.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{
		"user": user,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}

// Login godoc
// @Summary      Login
// @Description  login a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.LoginRequest true "Login Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var requestBody utils.LoginRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var getemail string
	var getPwd string
	var getUserId primitive.ObjectID

	var data interface{}

	if requestBody.Role == db.RoleAdmin {
		res, err := services.FindAdminByEmail(requestBody.Email)
		if err != nil {
			response.StatusCode = http.StatusBadRequest
			response.Message = err.Error()
			response.SendResponse(c)
			return
		}

		getemail = res.Email
		getPwd = res.Password
		getUserId = res.ID

		data = res

	} else if requestBody.Role == db.RoleUser {

		res, err := services.FindUserByEmail(requestBody.Email)
		if err != nil {
			response.StatusCode = http.StatusBadRequest
			response.Message = err.Error()
			response.SendResponse(c)
			return
		}

		getemail = res.Email
		getPwd = res.Password
		getUserId = res.ID

		data = res

	} else if requestBody.Role == db.RoleEmployee {

		res, err := services.FindEmployeeByEmail(requestBody.Email)
		if err != nil {
			response.Message = err.Error()
			response.SendResponse(c)
			return
		}

		getemail = res.Email
		getPwd = res.Password
		getUserId = res.ID

		data = res
	} else {
		response.StatusCode = http.StatusBadRequest
		response.Success = false
		response.Message = "Invalid Data"
		response.SendResponse(c)
		return
	}

	// check hashed password
	err := bcrypt.CompareHashAndPassword([]byte(getPwd), []byte(requestBody.Password))
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "email and password don't match"
		response.SendResponse(c)
		return
	}

	// generate new access tokens
	accessToken, refreshToken, err := services.GenerateTokens(getemail, getUserId)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		requestBody.Role: data,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}

// Refresh godoc
// @Summary      Refresh
// @Description  refreshes a user token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.RefreshRequest true "Refresh Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/refresh [post]
func Refresh(c *gin.Context) {
	var requestBody utils.RefreshRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// check token validity
	token, err := services.VerifyToken(requestBody.Token, db.TokenTypeRefresh)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	var getemail string
	var getUserId primitive.ObjectID

	var data interface{}

	if requestBody.Role == db.RoleAdmin {
		res, err := services.FindAdminById(token.User)
		if err != nil {
			response.Message = err.Error()
			response.SendResponse(c)
			return
		}

		getemail = res.Email
		getUserId = res.ID
		data = res

	} else if requestBody.Role == db.RoleUser {

		res, err := services.FindUserById(token.User)
		if err != nil {
			response.Message = err.Error()
			response.SendResponse(c)
			return
		}

		getemail = res.Email
		getUserId = res.ID
		data = res

	} else if requestBody.Role == db.RoleEmployee {

		res, err := services.FindEmployeeById(token.User)
		if err != nil {
			response.Message = err.Error()
			response.SendResponse(c)
			return
		}

		getemail = res.Email
		getUserId = res.ID
		data = res
	} else {
		response.StatusCode = http.StatusBadRequest
		response.Success = false
		response.Message = "Invalid Data"
		response.SendResponse(c)
		return
	}

	// delete old token
	err = services.DeleteTokenById(token.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	accessToken, refreshToken, err := services.GenerateTokens(getemail, getUserId)
	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"user": data,
		"token": gin.H{
			"access":  accessToken.GetResponseJson(),
			"refresh": refreshToken.GetResponseJson()},
	}
	response.SendResponse(c)
}
