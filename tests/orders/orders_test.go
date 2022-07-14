package test_orders

import (
	"testing"

	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/orders"
	"github.com/LovePelmeni/StoreService/orders"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"errors"
)


type OrderSuite struct {
	suite.Suite 
	Controller *gomock.Controller 
	MockedDatabaseController mock_database.MockDatabaseControllerInterface
	// mocked controller, that is responsible for handling Real Time Database Operations.. 
	MockedOrderCredentialsInterface *mock_orders.MockOrderCredentialsInterface // Mocked Credentials Interface 
	OrderController orders.OrderController 
}

func (this *OrderSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.MockedDatabaseController = mock_database.NewDatabaseController(this.Controller)
	this.MockedOrderCredentialsInterface = mock_orders.NewOrderCredentials(this.Controller)

	this.OrderController = orders.NewOrderController(
	this.MockedDatabaseController, this.MockedOrderCredentialsInterface)
}

func (this *OrderSuite) TeardownTest() {
	this.Controller.Finish()
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

func (this *OrderSuite) TestOrderCreate() {

	orderCredentials := orders.OrderCredentials{}

	this.MockedDatabaseController.InitializeFirebaseDatabase().EXPECT().Return(nil).Times(1)
	this.MockedOrderCredentialsInterface.GetCredentials().EXPECT().Return(orderCredentials).Times(1)

	OrderResponse, Error := this.OrderController.CreateOrder(&orderCredentials)
	ExpectedResponse, Exception := func()(bool, error){return true, nil}()
	assert.NoError(this.T(), Error)
	assert.Equal(this.T(), Error, Exception, "Exception should equals to None..")
	assert.Equal(this.T(), OrderResponse, ExpectedResponse, "Response Should be Equal To True...")
}

func (this *OrderSuite) TestOrderFailCreate() {

	DatabaseError := errors.New("Database Failed to Establish Connection.")
	orderCredentials := orders.OrderCredentials{}

	this.MockedDatabaseController.InitializeFirebaseDatabase().EXPECT().Return(nil, DatabaseError).Times(1)

	Response, error := this.OrderController.CreateOrder(&orderCredentials)
	assert.Equal(this.T(), error, DatabaseError)
	assert.Equal(this.T(), Response, nil)
}

func (this *OrderSuite) TestOrderCancelCreate() {

}

func (this *OrderSuite) TestOrderFailCancelCreate() {
	
}