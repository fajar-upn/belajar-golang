package transaction

import (
	"bwastartup/campaign"
	"errors"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

func NewService(repository Repository, campignRepository campaign.Repository) *service {
	return &service{repository, campignRepository}
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error) {

	campaign, err := s.campaignRepository.FindById(input.User.ID)
	if err != nil {
		return nil, err
	}

	if campaign.UserID != input.User.ID {
		return nil, errors.New("this campaign transaction is not own user")
	}

	transaction, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
