package test_payments

import (
	"errors"

	"testing"

	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments"
	grpcControllers "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/proto"
	mock_intent "github.com/LovePelmeni/OnlineStore/StoreService/mocks/intent"
	mock_session "github.com/LovePelmeni/OnlineStore/StoreServices/mocks/session"
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

func TestPaymentIntentSuite(t *testing.T) {
	suite.Run(t, new(PaymentIntentSuite))
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
	Controller                      *gomock.Controller
	MockedPaymentSessionCredentials *mock_session.MockedPaymentSessionCredentialsInterface
	MockedPaymentSessionClient      *mock_session.MockedPaymentSessionClientInterface
	PaymentSessionController        *payments.PaymentSessionController
}

func TestPaymentSessionSuite(t *testing.T) {
	suite.Run(t, new(PaymentSessionSuite))
}

func (this *PaymentSessionSuite) SetupTest() {

	MockedGrpcServerConnection := mock_grpc.NewMockGrpcServerConnection()

	this.Controller = gomock.NewController(this.T())
	this.MockedPaymentSessionClient = mock_session.NewMockPaymentSessionClient(MockedGrpcServerConnection)
	this.PaymentSessionController = payments.NewPaymentSessionController(this.MockedPaymentSessionClient)
}

func (this *PaymentSessionSuite) TestPaymentSessionCreate() {

	PaymentSessionCredentials := payments.NewPaymentSessionCredentials(
		struct{}{})

	this.MockedPaymentSessionCredentials.EXPECT().GetCredentials().Return(PaymentSessionCredentials, nil).Times(1)
	this.MockedPaymentSessionClient.EXPECT().CreatePaymentSession(
		grpcControllers.PaymentSessionParams{
			ProductId:   PaymentSessionCredentials.ProductId,
			PurchaserId: PaymentSessionCredentials.PurchaserId}).Return(true, nil).Times(1)

	Response, Error := this.PaymentSessionController.CreatePaymentSession(
		this.MockedPaymentSessionCredentials)

	assert.Equal(this.T(), Response, true, "Response should be equals to true.")
	assert.Equal(this.T(), Error, nil, "Error Should be equals to None, Because of Success Payment Session Response.")
}

type PaymentRefundSuite struct {
	suite.Suite
	Controller              *gomock.Controller
	MockedRefundCredentials *mock_refund.MockPaymentRefundCredentialsInterface
	MockedRefundClient      *mock_refund.MockPaymentRefundClientInterface
	PaymentRefundController *payments.PaymentRefundController
}

func TestPaymentRefundSuite(t *testing.T) {
	suite.Run(t, new(PaymentRefundSuite))
}

func (this *PaymentRefundSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.MockedRefundCredentials = mock_refund.NewMockPaymentRefundCredentialsInterface(this.Controller)
	this.MockedRefundClient = mock_refund.NewMockPaymentRefundClientInterface(this.Controller)
	this.PaymentRefundController = payments.NewPaymentRefundController(this.MockedrefundCredentials, this.MockedRefundClient)
}

func (this *PaymentRefundSuite) TeardownTest() {
	this.Controller.Finish()
}
func (this *PaymentRefundSuite) TestCreatePaymentRefund() {
	RefundCredentials := payments.PaymentRefundCredentials{}
	RefundResponse := struct{ RefundId string }{RefundId: "1"}

	this.MockedRefundCredentials.EXPECT().Validate().Return(
		RefundCredentials, nil).Times(1)
	this.MockedRefundClient.CreatePaymentRefund(
		gomock.Eq(RefundCredentials)).Return(RefundResponse, nil).Times(1)

	Response, Error := this.PaymentRefundController.CreatePaymentRefund(RefundCredentials)
	assert.Equal(Response, RefundResponse, "Refund Response should be an Struct.")
	assert.NoError(Error, "Error Should Equals To Nil, Because of Positive Response.")
}

func (this *PaymentRefundSuite) TestCreateFailCreatePaymentRefund() {

}
