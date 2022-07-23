package test_refunds

import (
	"testing"

	"github.com/LovePelmeni/Store/external_services/payments"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RefundSuite struct {
	suite.Suite
	Controller             *gomock.Controller
	PaymentIntentId        string
	MockedRefundController interface{}
}

func (this *RefundSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.PaymentIntentId = ""
	this.MockedRefundController = []interface{}{}
}

func (this *RefundSuite) TeardownTest() {
	this.Controller.Finish()
}

func TestRefundSuite(t *testing.T) {
	suite.Run(t, new(RefundSuite))
}

func (this *RefundSuite) TestCreateRefund() {
	this.MockedRefundController.EXPECT().CreateRefund().Return(true, nil).Times(1)
	response, error := payments.Payment.Refund()
	assert.Equal(error, nil)
}
