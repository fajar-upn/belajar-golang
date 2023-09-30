package transaction

type GetCampaignTransactionInput struct {
	ID int `uri:"id" binding:"required"`
}
