package payment

import (
	"bwastartup/user"
	"fmt"
	"os"
	"strconv"

	midtrans "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

type servicePayment struct {
}

type ServicePayment interface {
	GetPaymentURL(Transaction, user.User) (string, error)
}

func NewService() *servicePayment {
	return &servicePayment{}
}

func (s *servicePayment) GetPaymentURL(transaction Transaction, user user.User) (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}

	// 1. Set you ServerKey with globally
	midtrans.ServerKey = os.Getenv("SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox

	fmt.Println(midtrans.ServerKey)

	// 2. Initiate Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	// 3. Request create Snap transaction to Midtrans
	snapResp, _ := snap.CreateTransaction(req)

	return snapResp.RedirectURL, nil
}
