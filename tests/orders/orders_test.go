package test_orders

import (
	"testing"

	"errors"

	// "github.com/LovePelmeni/OnlineStore/StoreService/external_services/orders"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

type OrderSuite struct {
	suite.Suite
	Controller             *gomock.Controller
	MockedFirebaseManager  *mock_firebase.MockFirebaseDatabaseOrderManagerInterface
	MockedOrderCredentials *mock_orders.MockOrderCredentialsInterface
	OrderController        *orders.OrderController
}

func (this *OrderSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.MockedFirebaseManager = mock_firebase.NewMockFirebaseDatabaseOrderManagerInterface(this.Controller)
	this.OrderController = orders.NewOrderController(this.MockedFirebaseManager)
	this.MockedOrderCredentials = mock_orders.NewMockOrderCredentialsInteface(this.Controller)
}

func (this *OrderSuite) TeardownTest() {
	this.Controller.Finish()
}

func (this *OrderSuite) TestCreateOrder() {

	TestOrderCredentials := struct {
	}{}

	this.MockedOrderCredentials.EXPECT().GetCredentials().Return(TestOrderCredentials, nil).Times(1)
	this.MockedFirebaseManager.EXPECT().CreateOrder(gomock.Eq(TestOrderCredentials)).Return(true, nil).Times(1)

	ActualResponse, Error := this.OrderController.CreateOrder(this.MockedOrderCredentials)
	assert.Equal(this.T(), ActualResponse, true)
	assert.NoError(this.T(), Error[0])
}

func (this *OrderSuite) TestFailCreateOrder() {

	TestOrderCredentials := struct {
	}{}

	FailedToSaveError := errors.New("Failed to Create Order to Firebase Database.")
	this.MockedOrderCredentials.EXPECT().GetCredentials().Return(TestOrderCredentials, nil).Times(1)
	this.MockedFirebaseManager.EXPECT().CreateOrder(gomock.Eq(TestOrderCredentials)).Return(false, FailedToSaveError).Times(1)

	Response, Exception := this.OrderController.CreateOrder(this.MockedOrderCredentials)
	assert.Equal(this.T(), Response, false, "Response Should be Equal To False, Because of Saving Order Failure.")
	assert.Equal(this.T(), Exception, FailedToSaveError, "Error should be equal to the specific Exception.")
}

func (this *OrderSuite) TestCancelOrder() {
	OrderId := "some-order-identifier"

	this.MockedFirebaseManager.EXPECT().CancelOrder(gomock.Eq(OrderId)).Return(true, nil).Times(1)
	Response, Exception := this.OrderController.CancelOrder(OrderId)
	assert.Equal(this.T(), Response, true, "Response Should be Equal to True, because of Positive Behaviour of Firebase.")
	assert.NoError(this.T(), Exception)
}

func (this *OrderSuite) TestFailCancelOrder() {
	OrderId := "some-order-identifier"
	FailedToCancelError := errors.New("Failed to Cancel Order.")
	this.MockedFirebaseManager.EXPECT().CancelOrder(gomock.Eq(OrderId)).Return(false, FailedToCancelError).Times(1)
	Response, Exception := this.OrderController.CancelOrder(OrderId)

	assert.Equal(this.T(), Response, false, "Response should be equals to False, Because of Canceling Order Failure.")
	assert.Equal(this.T(), Exception, FailedToCancelError, "Error Should be equal to the `FailedToCancelError`")
}
