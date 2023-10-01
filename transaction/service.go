package transaction

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionInput, userID int) ([]Transaction, error)
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionInput, userID int) ([]Transaction, error) {
	transaction, err := s.repository.GetByCampaignID(input.ID, userID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
