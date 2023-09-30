package transaction

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error) {
	transaction, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
