package test_payments

import (
	"testing"

	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"errors"
	"github.com/stretchr/testify/assert"
)

type PaymentIntentSuite struct {
	suite.Suite
	Controller                    *gomock.Controller
	MockedPaymentGRPCClientController *mock_payments.MockPaymentIntentClient
	Payment                       *payments.PaymentInfoCredentials
	PaymentIntentController 	*payments.Payment 
}

func (this *PaymentIntentSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.MockedPaymentGRPCClientController = mock_payments.NewMockPaymentIntentClient(this.Controller, this.Payment)
	this.Payment = &payments.PaymentInfoCredentials{
		PaymentSessionId: "some-session-id",
		ProductId: "1", 
		CustomerId: "1",
		Currency: "USD",
		Price: "1.00",
	}
}

func TestPaymentIntentSuite(t *testing.T) {
	suite.Run(t, new(PaymentIntentSuite))
}

func (this *PaymentIntentSuite) TeardownTest() {
	this.Controller.Finish()
}

func (this *PaymentIntentSuite) TestCreatePaymentIntent() {

	ProductId := "1"
	CustomerId := "5"
	Currency := "USD"
	Price := "5.00"

	PaymentIntentCredentials := payments.PaymentInfoCredentials{
		ProductId: ProductId,
		Price: Price,
		CustomerId: CustomerId,
		Currency: Currency,

	}

	this.MockedPaymentGRPCClientController.EXPECT().CreatePaymentIntent(
	gomock.Eq(this.Payment, PaymentIntentCredentials)).Return(true, nil).Times(1)

	ExpectedResponse := true 
	ActualResponse, error := this.PaymentIntentController.Pay(PaymentIntentCredentials) 

	assert.Equal(ExpectedResponse, ActualResponse)
	assert.Equal(error, nil)
	assert.NoError(error)
}

func (this *PaymentIntentSuite) TestFailCreatePaymentIntent() {
}


type PaymentSessionSuite struct {
	suite.Suite 
	Controller *gomock.Controller 
	MockedPaymentGRPCClientController *mock_payments.MockPaymentClient 
	PaymentController 				  *payments.Payment 
	payments.PaymentInfoCredentials	  *payments.PaymentInfoCredentials

}

func (this *PaymentSessionSuite) SetupTest() {
	this.Controller = *gomock.NewController(this.T())
	this.PaymentController = 
}
func (this *PaymentIntentSuite) TestCreatePaymentSession() {

	ProductId := "1"
	CustomerId := "2"
	Currency := "USD"
	Price := "5.00"
	
	PaymentSessionCredentials := payments.PaymentInfoCredentials{
		ProductId: ProductId,
		CustomerId: CustomerId,
		Currency: Currency,
		Price: Price,
	}
	this.MockedPaymentGRPCClientController.EXPECT().CreatePaymentSession(gomock.Eq(this.Client, PaymentSessionCredentials)).Return(true, nil).Times(1)
	ActualResponse, error := this.Payment.StartPaymentSession(PaymentSessionCredentials) 
	ExpectedResponse, error := true, nil 
	assert.Equal(ExpectedResponse, ActualResponse)
	assert.NoError(error)
}


func (this *PaymentIntentSuite) TestFailCreatePaymentSession(){

}