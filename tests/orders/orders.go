package test_orders

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type OrderSuite struct {
	suite.Suite
	Controller                   *gomock.Controller
	OrderData                    map[string]string
	MockedRabbitmqTransportLayer *mocks.NewRabbitmqTransportInterface
	OrderStructController        *orders.OrderController
}

func (this *OrderSuite) TeardownTest() {
	this.Controller.Finish()
}

func (this *OrderSuite) SetupTest() {
	this.OrderData = map[string]string{}
	this.Controller = gomock.NewController(this.T())
	this.MockedRabbitmqTransportLayer = mocks.NewRabbitmqTransportInterface(this.Controller)
	this.OrderStructController = orders.NewOrderController(
		this.MockedRabbitmqTransportLayer, this.OrderData)
}

func (this *OrderSuite) TestOrderSendEvent() {
}

func (this *OrderSuite) TestOrderSendFailEvent() {
}

func (this *OrderSuite) TestOrderReceiveEvent() {
}

func (this *OrderSuite) TestOrderReceiveFailEvent() {
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}
