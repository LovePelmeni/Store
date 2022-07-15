package test_payments

import (
	"errors"

	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments"
	grpcControllers "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PaymentIntentSuite struct {
	suite.Suite
	Controller                     *gomock.Controller
	MockedPaymentIntentClient      *mock_intent.MockPaymentIntentClientInterface
	MockedPaymentIntentCredentials *mock_intent.MockPaymentIntentCredentialsInterface
	PaymentIntentController        *payments.PaymentIntentController
}

func (this *PaymentIntentSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.MockedPaymentIntentClient = mock_intent.NewMockPaymentIntentClientInterface(this.Controller)
	this.PaymentIntentController = payments.NewPaymentIntentController(this.MockedPaymentIntentClient)
}

func (this *PaymentIntentSuite) TeardownTest() {
	this.Controller.Finish()
}

func (this *PaymentIntentSuite) TestPaymentIntentCreate() {

	TestPaymentIntentId := "test-payment-intent-id"
	ExpectedgRPCResponse := grpcControllers.PaymentIntentResponse{PaymentIntentId: TestPaymentIntentId}
	ExpectedResponse := map[string]string{"PaymentIntentId": TestPaymentIntentId}
	RequestParams := payments.PaymentIntentCredentials{}

	this.MockedPaymentIntentCredentials.EXPECT().GetCredentials().Return(RequestParams, nil).Times(1)
	this.MockedPaymentIntentClient.EXPECT().CreatePaymentIntent().Return(ExpectedgRPCResponse, nil).Times(1)

	Response, Error := this.PaymentIntentController.CreatePaymentIntent(RequestParams)
	assert.Equal(this.T(), Response, ExpectedResponse, "Response need to be equal to dictionary with `PaymentIntentId` Specified.")
	assert.Equal(this.T(), Error, nil, "Error Need to be Equal to Null")
}

func (this *PaymentIntentSuite) TestPaymentIntentFailCreate() {

	ExpectedgRPCResponse := grpcControllers.PaymentIntentResponse{PaymentIntentId: "-"}
	ExpectedResponseError := errors.New("Payment Intent Failure")
	RequestedParams := payments.PaymentIntentCredentials{}

	this.MockedPaymentIntentCredentials.EXPECT().GetCredentials().Return(RequestedParams, nil).Times(1)
	this.MockedPaymentIntentClient.EXPECT().Return(ExpectedgRPCResponse, nil).Times(1)

	Response, Error := this.PaymentIntentController.CreatePaymentIntent(RequestedParams)
	assert.Equal(this.T(), Response, nil, "Response need to be equals to false, due to Error.")
	assert.Equal(this.T(), Error, ExpectedResponseError, "Error need to be not Null.")
}

type PaymentSessionSuite struct {
	suite.Suite
	Controller *gomock.Controller
}

type PaymentCheckoutSuite struct {
	suite.Suite
	Controller *gomock.Controller
}
