package campaign

import "strings"

/**
this code for format campaigns to appropriate with result json
*/
type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageUrl         string `json:"image_url"`
	Slug             string `json:"slug"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}

	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount

	campaignFormatter.ImageUrl = ""
	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageUrl = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	campaignsFormatter := []CampaignFormatter{} //this is to handle campainFormatter when campain is null

	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}

/**
this code for format detail campaigns to appropriate with result json
*/
type CampaignDetailFormatter struct {
	ID               int                       `json:"id"`
	Name             string                    `json:"name"`
	ShortDescription string                    `json:"short_description"`
	Description      string                    `json:"description"`
	ImagesUrl        string                    `json:"images_url"`
	GoalAmout        int                       `json:"goal_amount"`
	CurrentAmount    int                       `json:"current_amount"`
	UserID           int                       `json:"user_id"`
	Slug             string                    `json:"slug"`
	Perks            []string                  `json:"perks"`
	User             CampaignUserFormatter     `json:"user"`
	Images           []CampaignImagesFormatter `json:"images"`
}

type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type CampaignImagesFormatter struct {
	ImageUrl  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {

	/**
	this code will return campaign json
	*/
	campaignDetailFormatter := CampaignDetailFormatter{}
	campaignDetailFormatter.ID = campaign.ID
	campaignDetailFormatter.Name = campaign.Name
	campaignDetailFormatter.ShortDescription = campaign.ShortDescription
	campaignDetailFormatter.Description = campaign.Description
	campaignDetailFormatter.GoalAmout = campaign.GoalAmount
	campaignDetailFormatter.CurrentAmount = campaign.CurrentAmount
	campaignDetailFormatter.UserID = campaign.UserID
	campaignDetailFormatter.Slug = campaign.Slug

	campaignDetailFormatter.ImagesUrl = ""
	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImagesUrl = campaign.CampaignImages[0].FileName
	}

	var perks []string
	for _, perk := range strings.Split(campaign.Perks, ",") { //strings.Split : use for split string appropriate with parameter (example: ',' = comma )
		perks = append(perks, strings.TrimSpace(perk)) //string.TrimSpace : use for remove space from sentence
	}
	campaignDetailFormatter.Perks = perks

	/**
	this code will return user object inside campaign json
	*/
	campaignUserFormatter := CampaignUserFormatter{}
	user := campaign.User //call struct user
	campaignUserFormatter.Name = user.Name
	campaignUserFormatter.ImageUrl = user.AvatarFileName

	campaignDetailFormatter.User = campaignUserFormatter //this code for save user struct inside campaign json

	/**
	this code will return images object inside campaign json
	*/
	images := []CampaignImagesFormatter{}
	for _, image := range campaign.CampaignImages {
		campaignImagesFormatter := CampaignImagesFormatter{}
		campaignImagesFormatter.ImageUrl = image.FileName
		isPrimary := false
		if image.IsPrimary == 1 {
			isPrimary = true
		}
		campaignImagesFormatter.IsPrimary = isPrimary

		images = append(images, campaignImagesFormatter) //this code for add images struct to array 'images'
	}
	campaignDetailFormatter.Images = images //this code for save images struct inside campaign json

	return campaignDetailFormatter
}
