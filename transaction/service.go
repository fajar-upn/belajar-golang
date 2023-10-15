package transaction

import (
	"bwastartup/campaign"
	"errors"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionInput, userID int) ([]Transaction, error)
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionInput, userID int) ([]Transaction, error) {

	campaign, err := s.campaignRepository.FindById(input.ID)
	if err != nil {
		return nil, err
	}

	if campaign.UserID != input.User.ID {
		return nil, errors.New("This user not owner of the campign")
	}

	transaction, err := s.repository.GetByCampaignID(input.ID, userID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
