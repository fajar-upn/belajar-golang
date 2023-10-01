package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{db}
}

type Repository interface {
	GetByCampaignID(campaignID int, userID int) ([]Transaction, error)
}

func (r *repository) GetByCampaignID(campaignID int, userID int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
