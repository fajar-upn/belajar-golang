package payment

import (
	"bwastartup/user"
	"fmt"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type servicePayment struct{}

type ServicePayment interface {
	GetPaymentURL(Transaction, user.User) (string, error)
}

func NewService() *servicePayment {
	return &servicePayment{}
}

func (s *servicePayment) GetPaymentURL(transaction Transaction, user user.User) (string, error) {

	// 1. Set you ServerKey with globally
	midtrans.ServerKey = ""
	midtrans.Environment = midtrans.Sandbox

	// 2. Initiate Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	// 3. Request create Snap transaction to Midtrans
	snapResp, err := snap.CreateTransaction(req)
	if err != nil {
		fmt.Println("Error here")
		return "", err
	}
	fmt.Println("Response :", snapResp)

	return snapResp.RedirectURL, nil
}
