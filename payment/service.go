package payment

import (
	"bwastartup/transaction"
	"bwastartup/user"
	"strconv"

	"github.com/veritrans/go-midtrans"
)

type servicePayment struct{}

type ServicePayment interface {
	GetPaymenURL(transaction.Transaction, user.User) (string, error)
}

func NewService() *servicePayment {
	return &servicePayment{}
}

func (s *servicePayment) GetPaymenURL(transaction transaction.Transaction, user user.User) (string, error) {
	midClient := midtrans.NewClient()
	midClient.ServerKey = ""
	midClient.ClientKey = ""

	snapGateway := midtrans.SnapGateway{
		Client: midClient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil

}
