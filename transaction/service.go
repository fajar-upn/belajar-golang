package transaction

import (
	"bwastartup/campaign"
	"bwastartup/payment"
	"errors"
	"fmt"
	"math/rand"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.ServicePayment
}

func NewService(repository Repository, campignRepository campaign.Repository, paymentService payment.ServicePayment) *service {
	return &service{repository, campignRepository, paymentService}
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
	CreateTransactionService(input CreateTransactionInput) (*Transaction, error)
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

func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error) {
	transaction, err := s.repository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *service) CreateTransactionService(input CreateTransactionInput) (*Transaction, error) {

	RandomInteger := rand.Int()
	RandomCode := fmt.Sprintf("ORDER-%d", RandomInteger)

	transaction := Transaction{
		CampaignID: input.CampaignID,
		Amount:     input.Amount,
		UserID:     input.User.ID,
		Status:     "Pending",
		Code:       RandomCode,
	}

	newTransaction, err := s.repository.CreateTransactionRepository(transaction)
	if err != nil {
		return nil, err
	}

	return newTransaction, nil
}
