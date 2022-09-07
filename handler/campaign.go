package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/**
step by step insert campaign:
1. receive parameter in handler
2. handler to service
3. service will consider call repository
4. Repository: FindAll, FindByUserID
5. DB
*/
type CampaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *CampaignHandler {
	return &CampaignHandler{service}
}

// get api 'api/v1/campaign'
func (h *CampaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id")) //c.Query = 'api/v1/campaign?user_id'

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		// response error to JSON
		response := helper.APIResponse("Error to get campaigns!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

/**
step by step detail campaign:

2.
3. in handler mapping id which url to struct input, this struct will call service and will format json output
4. service for input struct to receive id from url, this service will call get campaign by id from repository
5. repository for get campaign by id
*/
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Detail Campaign", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

/**
step by step add campaign:

1. get parameter from user to input struct
2. get current user from JWT or handler
3. call service, where parameter input (can create slug automatically)
4. call repository for save new campaign data
*/
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User) //get data from user table
	input.User = currentUser

	newCampign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Create campaign has been success", http.StatusOK, "success", campaign.FormatCampaign(newCampign))
	c.JSON(http.StatusOK, response)
}

/**
step by step update campaign:
1. User insert input
2. mapping from input user (API) and input URI  to input struct (object) (handler)
3. parsing to service layer with parameter input (service)
4. in the service find appropriate campaign id data from repository, after that catch parameter and insert to struct in service (service)
4. repository update the campaign
*/
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID) // ShouldByUri for tag uri from url
	if err != nil {
		response := helper.APIResponse("Failed to update campaign in input id json", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput
	err1 := c.ShouldBindJSON(&inputData)

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	if err1 != nil {
		errors := helper.FormatValidationError(err1)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update Campaign in input data json", http.StatusUnprocessableEntity, "errors", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign in service", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

/**
step by step upload image campaign

1. handler : catch input from repository and change struct input, save campaign image to a folder
2. service : have condition paint 5 (call repository), ruh repository in service
3. repository : create image or save data image in table data images, change is_primary true to false
*/
func (h *CampaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input) //this is form so using ShouldBind()
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to Upload Campaign Image from input JSON", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": file}
		response := helper.APIResponse("Failed to Upload Campaign Image from input JSON File", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err1 := c.SaveUploadedFile(file, path)
	if err1 != nil {
		data := gin.H{"is_uploaded": file}
		response := helper.APIResponse("Failed to Upload Campaign Image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err2 := h.service.SaveCampaignImage(input, path)
	if err2 != nil {
		data := gin.H{"is_uploaded": file}
		response := helper.APIResponse("Failed to Upload Campaign Image in Service", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Success upload image", http.StatusOK, "error", data)
	c.JSON(http.StatusOK, response)
}
