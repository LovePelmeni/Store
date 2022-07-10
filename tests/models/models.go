package test_models

import (
	"testing"

	mocked_models "github.com/LovePelmeni/OnlineStore/StoreService/mocks/models"
	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type BaseModelSuite struct {
	suite.Suite
	Controller       *gomock.Controller
	Models           []models.BaseModel
	MockedController *mocked_models.MockBaseModel
}

func (this *BaseModelSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.Models = []models.BaseModel{&models.product, &models.customer, &models.cart}
	this.MockedController = mocked_models.NewMockBaseModel(this.Controller)
}

func (this *BaseModelSuite) TeardownTest() {
	this.Controller.Finish()
}

func TestBaseModelSuite(t *testing.T) {
	suite.Run(t, new(BaseModelSuite))
}

func (this *BaseModelSuite) CreateModel() {}

func (this *BaseModelSuite) UpdateModel() {}

func (this *BaseModelSuite) DeleteModel() {}
