package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{db}
}

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
}

func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transaction []Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Find(&transaction).Error
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
