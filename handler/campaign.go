package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
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
