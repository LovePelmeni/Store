package test_rest_customers

import (
	"net/http"
	"net/url"
	"testing"

	"fmt"
	"os"

	"github.com/gorilla/csrf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestDatabaseConnection struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
}

func NewTestDatabaseConnection() *TestDatabaseConnection {
	return &TestDatabaseConnection{Host: os.Getenv("TEST_POSTGRES_HOST"),
		Port: os.Getenv("TEST_POSTGRES_PORT"), User: os.Getenv("TEST_POSTGRES_USER"),
		Password: os.Getenv("TEST_POSTGRES_PASSWORD")}
}

func (this *TestDatabaseConnection) GetConnection() *gorm.DB {
	Database, Error := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s password=%s user=%s dbname= %s",
			this.Host, this.Port, this.Password, this.User, this.DatabaseName),
	}))
	if Error != nil {
		panic("Failed to Connect To Test Database..")
	}
	return Database
}

type CustomerRestControllerSuite struct {
	suite.Suite
	TestDatabaseConnection *gorm.DB
	RequestBodyPayload     map[string]string
	RequestQueryParams     map[string]string
	Client                 http.Client
}

func (this *CustomerRestControllerSuite) SetupTest() {

	this.Client = http.Client{}
	this.RequestBodyPayload = map[string]string{
		"Username": "Test-Username",
		"Email":    "Test-Email",
		"Password": "Test-Password",
	}
	this.RequestQueryParams = nil
	this.TestDatabaseConnection = NewTestDatabaseConnection().GetConnection()
}

func TestCustomerSuite(t *testing.T) {
	suite.Run(t, new(CustomerRestControllerSuite))
}

func (this *CustomerRestControllerSuite) TestCreateCustomer() {

	RequestUrl := url.URL{Path: "http://localhost:8000/create/customer/"}
	newHttpRequest, Error := http.NewRequest("POST", RequestUrl.String(), nil)

	if Error != nil {
		assert.Error(this.T(), Error, "Failed to Setup Http Request.")
	}

	newHttpRequest.Header.Set("Access-Control-Allow-Origin", "localhost")
	newHttpRequest.Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	newHttpRequest.Header.Set("Access-Control-Allow-Credentials", "true")
	newHttpRequest.Header.Set("X-CSRF-TOKEN", csrf.Token(newHttpRequest))

	Response, HttpError := this.Client.Do(newHttpRequest)
	assert.Equal(this.T(), Response.StatusCode, 201, "Server Responded with not positive Code.")
	assert.NoError(this.T(), HttpError, "Error Should be Equal to None")

	var count int64
	CountedTransactions := this.TestDatabaseConnection.Table("customers").Count(&count)
	if CountedTransactions.Error != nil {
		assert.Error(this.T(), CountedTransactions.Error, "Failed to Count Transactions.")
	}

	assert.Equal(this.T(), count, 1, "Transactions Quantity Should be Equals to 1.")
}

func (this *CustomerRestControllerSuite) TestUpdateCustomer() {

	TestCustomerId := "1" // Test Customer ID
	RequestUrl := url.URL{Path: "http://localhost:8000/update/customer/"}
	newHttpRequest, Error := http.NewRequest("PUT", RequestUrl.String(), nil)
	newHttpRequest.URL.Query().Add("CustomerId", TestCustomerId)

	if Error != nil {
		assert.Error(this.T(), Error, "Failed to Setup Http Request.")
	}

	newHttpRequest.Header.Set("Access-Control-Allow-Origin", "localhost")
	newHttpRequest.Header.Set("Access-Control-Allow-Credentials", "true")
	newHttpRequest.Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	newHttpRequest.Header.Set("X-CSRF-TOKEN", csrf.Token(newHttpRequest))

	Response, HttpError := this.Client.Do(newHttpRequest)
	assert.Equal(this.T(), Response.StatusCode, 201, "Server Responded with not positive Code.")
	assert.NoError(this.T(), HttpError, "Error Should be Equal to None")
}

func (this *CustomerRestControllerSuite) TestGetProfile() {

	TestCustomerId := "1" // Test Customer ID
	RequestUrl := url.URL{Path: "http://localhost:8000/get/customer/"}
	newHttpRequest, Error := http.NewRequest("POST", RequestUrl.String(), nil)
	newHttpRequest.URL.Query().Add("CustomerId", TestCustomerId)

	if Error != nil {
		assert.Error(this.T(), Error, "Failed to Setup Http Request.")
	}

	newHttpRequest.Header.Set("Access-Control-Allow-Origin", "localhost")
	newHttpRequest.Header.Set("Access-Control-Allow-Credentials", "true")
	newHttpRequest.Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	newHttpRequest.Header.Set("X-CSRF-TOKEN", csrf.Token(newHttpRequest))

	Response, HttpError := this.Client.Do(newHttpRequest)
	assert.Equal(this.T(), Response.StatusCode, 201, "Server Responded with not positive Code.")
	assert.NoError(this.T(), HttpError, "Error Should be Equal to None")
}

func (this *CustomerRestControllerSuite) TestDeleteCustomer() {

	TestCustomerId := "1" // Test Customer ID
	RequestUrl := url.URL{Path: "http://localhost:8000/delete/customer/"}
	newHttpRequest, Error := http.NewRequest("POST", RequestUrl.String(), nil)
	newHttpRequest.URL.Query().Add("CustomerId", TestCustomerId)

	if Error != nil {
		assert.Error(this.T(), Error, "Failed to Setup Http Request.")
	}

	newHttpRequest.Header.Set("Access-Control-Allow-Origin", "localhost")
	newHttpRequest.Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	newHttpRequest.Header.Set("Access-Control-Allow-Credentials", "true")
	newHttpRequest.Header.Set("X-CSRF-TOKEN", csrf.Token(newHttpRequest))

	Response, HttpError := this.Client.Do(newHttpRequest)
	assert.Equal(this.T(), Response.StatusCode, 201, "Server Responded with not positive Code.")
	assert.NoError(this.T(), HttpError, "Error Should be Equal to None")

	var count int64
	CountedTransactions := this.TestDatabaseConnection.Table("customers").Count(&count)
	if CountedTransactions.Error != nil {
		assert.Error(this.T(), CountedTransactions.Error, "Failed to Count Transactions.")
	}

	assert.Equal(this.T(), count, 0, "Transactions Quantity Should be Equals to 0.")
}

// Product Rest Controllers

type ProductRestControllerSuite struct {
	suite.Suite
}

func (this *ProductRestControllerSuite) SetupTest() {}

func TestProductRestControllerSuite(t *testing.T) {
	suite.Run(t, new(ProductRestControllerSuite))
}

func (this *ProductRestControllerSuite) TestCreateProduct() {}

func (this *ProductRestControllerSuite) TestUpdateProduct() {}

func (this *ProductRestControllerSuite) TestDeleteProduct() {}
