package test_models

import (
	"fmt"
	"os"
	"testing"

	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/stretchr/testify/assert"
)

type TestDatabaseConnection struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
}

var (
	TEST_POSTGRES_HOST     = os.Getenv("TEST_POSTGRES_HOST")
	TEST_POSTGRES_PORT     = os.Getenv("TEST_POSTGRES_PORT")
	TEST_POSTGRES_USER     = os.Getenv("TEST_POSTGRES_USER")
	TEST_POSTGRES_PASSWORD = os.Getenv("TEST_POSTGRES_PASSWORD")
	TEST_POSTGRES_DATABASE = os.Getenv("TEST_POSTGRES_DATABASE")
)

func (this *TestDatabaseConnection) StartSession() (*gorm.DB) {
	Database, error := gorm.Open(postgres.New(
		postgres.Config{
			DSN: fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
				TEST_POSTGRES_HOST, TEST_POSTGRES_PORT, TEST_POSTGRES_DATABASE, TEST_POSTGRES_USER, TEST_POSTGRES_PASSWORD),
			PreferSimpleProtocol: true,
		},
	))
	if error != nil {
		panic("Testing Database is not running...!")
	}
	return Database, nil 
}

type ModelSuite struct {
	suite.Suite

	Controller           *gomock.Controller
	MockedModelInterface mock_models.NewMockBaseModel
	ListModels           []models.BaseModel

	ModelsTestCreationCasesData []struct{} // struct that represents Model Data Payload for Create Method
	ModelsTestUpdateCasesData   []struct{} // struct that represents Model Data Payload for Update Method
	ModelsTestDeleteCasesData   []struct{} // struct that represents Model Data Payload for Delete Method

	TestDatabaseConnection *TestDatabaseConnection
}

func (this *ModelSuite) SetupTest() {
	this.Controller = gomock.NewController(this.T())
	this.MockedModelInterface = mock_models.NewMockBaseModel(this.Controller)
	this.ListModels = []models.BaseModel{&models.customer, &models.cart, &models.product}

	this.ModelsTestCreationCasesData = []struct{}{}

	this.ModelsTestUpdateCasesData = []struct{}{}

	this.ModelsTestDeleteCasesData = []struct{}{}
}

func (this *ModelSuite) TeardownTest() {
	this.Controller.Finish()
}

func TestModelSuite(t *testing.T) {
	suite.Run(t, new(ModelSuite))
}

func (this *ModelSuite) TestModelsCreate(t *testing.T) {

	databaseConnection := this.TestDatabaseConnection.StartSession()


	// Creation Models Test Cases Goes There...


	// Products 
	testProductCase := func(t *testing.T) {newProduct := models.Product{}; 
	Saved := databaseConnection.Table("products").Save(&newProduct); assert.Equal(t, Saved.Error, nil)}

	// Customers 
	testCustomerCase := func(t *testing.T) {newCustomer := models.Customer{};
    Saved := databaseConnection.Table("customers").Save(&newCustomer); assert.Equal(t, Saved.Error, nil)}

	// Carts
	testCartCase := func(t *testing.T) {newCart := models.Cart{};
    Saved := databaseConnection.Table("carts").Save(&newCart); assert.Equal(t, Saved.Error, nil)}


	testing.RunTests(func(pt string, str string) (bool, error) {
		return true, nil
	}, []testing.InternalTest{
		{"Test Model Customer", testCustomerCase},
		{"Test Model Cart", testCartCase},
		{"Test Model Product", testProductCase},
	})
}

func (this *ModelSuite) TestModelUpdate(t *testing.T) {

	testCustomerCase := func(t *testing.T) {}

	testProductCase := func(t *testing.T) {}

	testCartCase := func(t *testing.T) {}

	testing.RunTests(func(pat string, str string) (bool, error) {
		return true, nil
	},
		[]testing.InternalTest{
			{"Test Update Customer", testCustomerCase},
			{"Test Update Product", testProductCase},
			{"Test Update Cart", testCartCase},
		})
}

func (this *ModelSuite) TestModelsDelete(t *testing.T) {

	testCustomerCase := func(t *testing.T) {}
	testProductCase := func(t *testing.T) {}
	testCartCase := func(t *testing.T) {}

	testing.RunTests(func(pat string, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{"Test Delete Customer", testCustomerCase},
			{"Test Delete Product", testProductCase},
			{"Test Delete Cart", testCartCase},
		})
}
