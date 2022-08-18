package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

// this code call in main.go
func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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

	token, err := h.authService.GenerateToken(newUser.ID) //jwt service
	if err != nil {
		response := helper.APIResponse("Register account has been failed!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Registration account success", http.StatusOK, "success", formatter)

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

	token, err := h.authService.GenerateToken(loggedInUser.ID) //jwt service
	if err != nil {
		response := helper.APIResponse("Login Failed!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)

	response := helper.APIResponse("Login success", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

// error handling when email already exist
func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	/**
	step by step check available email
	1. check email input from user form (user.input.go)
	2. input email will be mapping to struct input (handler.user.go)
	3. input struct will passing (handler.user.go) to service (user.service.go)
	4. service (user.service.go) will be call repository (user.repository)
	5. call database in repository (user.repository)
	*/
	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email has been available!", http.StatusUnauthorized, "errors", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}

		response := helper.APIResponse("Email check error", http.StatusUnauthorized, "errors", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	// organize status to API
	var metaMessage string
	if !isEmailAvailable {
		metaMessage = "Email has been registered"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {

	/**
	note: we don't need mapping image
	this is step by step for upload image
	1. Receive input from user
	2. save image in "images/" folder
	3. in the service (user.service.go) we call repository (user.service.go)
	3.1 for call access user with JWT (for a while user has been login/ hardcode to login because using jwt with appropriate ID)
	3.2 call repository with appropriate user id, ex: user_id=1
	3.2 with repository user data will update and save in file location
	*/
	file, err := c.FormFile("avatar") //'avatar' is key form data in postman

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	// userID := 9 //suppose userId get from JWT, because JWT not available right now
	/**
	fmt.printf is concat string.
	file.filename for get filename
	(images/<ID>/<IMAGE_PATH>)
	%d will be override with "userID"
	%s will be override with "file.Filename"
	*/
	path := fmt.Sprintf("images/%d - %s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path) //save file(image)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image!", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image!", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Success to upload avatar image", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}
