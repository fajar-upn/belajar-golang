package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

// this code call in main.go
func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	/**
	1. receive input from user
	2. map input from user to RegisterUserInput struct
	3. struct NewHandlerUser will pass as a parameter service
	*/

	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	// error because validation
	if err != nil {

		//-----this code for error formatting
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"error": errors} //gin.H is map[string]interface{}
		// -----

		// response error to JSON
		response := helper.APIResponse("Register account has been failed!", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register account has been failed!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//token, err := h.jwtService.GenerateToken()

	formatter := user.FormatUser(newUser, "JWTTOKENNOTACTIVATED")

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	/**
	step by step login
	1. user input email & password
	2. input will receive by handler (handler.user.go)
	3. mapping from input user to struct (user.input.go)
	4. input struct passing to service (user.service.go)
	4.1. in the service will find appropriate email & password
	4.2. in the service we need match email and password
	*/
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		// response error to JSON
		response := helper.APIResponse("Login Failed!", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		// response error to JSON
		response := helper.APIResponse("Login Failed!", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, "JWTTOKENNOTACTIVATEDLOGIN")

	response := helper.APIResponse("Login success", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}
