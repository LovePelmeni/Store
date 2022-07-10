package test_payments

import (
	"testing"

	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type PaymentIntentSuite struct {
	suite.Suite
	Controller                    *gomock.Controller
	MockedPaymentIntentController interface{}
	Payment                       *payments.Payment
}

func (this *PaymentIntentSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.MockedPaymentIntentController = []string{}
	this.Payment = payments.Payment{}
}

func TestPaymentIntentSuite(t *testing.T) {
	suite.Run(t, new(PaymentIntentSuite))
}

func (this *PaymentIntentSuite) TeardownTest() {
	this.Controller.Finish()
}

func (this *PaymentIntentSuite) TestCreatePaymentIntent() {
	this.MockedPaymentIntentController.EXPECT().CreatePaymentIntent().Return(true, nil).Times(1)
}

func (this *PaymentIntentSuite) TestCreatePaymentSession() {
	this.MockedPaymentIntentController.EXPECT().CreatePaymentSession().Return(true, nil).Times(1)
}
