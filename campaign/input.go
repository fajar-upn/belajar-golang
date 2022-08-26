package campaign

type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"` //URI : example:"localhost/api/v1/campaigns/1" number 1 in least url is uri
}
