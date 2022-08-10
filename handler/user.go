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
