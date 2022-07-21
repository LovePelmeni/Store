package test_rest_products

// import (
// 	"fmt"
// 	"net/url"
// 	"testing"

// 	"net/http"

// 	"github.com/LovePelmeni/OnlineStore/StoreService/rest/products"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"log"
// 	"os"
// 	"github.com/LovePelmeni/OnlineStore/StoreService/models"
// )

// var (
// 	DebugLogger *log.Logger
// 	InfoLogger *log.Logger
// 	ErrorLogger *log.Logger
// )

// var (
// 	ServerHost = os.Getenv("APPLICATION_HOST")
// 	ServerPort = os.Getenv("APPLICATION_PORT")
// )

// func InitializeLoggers() (bool, error) {
// 	return true, nil
// }

// func init() {
// 	Initialized, Error := InitializeLoggers()
// 	if !Initialized || Error != nil {panic(Error)}
// }
// type RestProductControllerSuite struct {
// 	suite.Suite
// 	RequestHeaders map[string]string // list of headers that should be applied for every request...
// 	ProductCreateData *models.Product
// 	ProductUpdateData struct{ProductName string;
// 	ProductDescription string; ProductPrice string; Currency string}
// }

// func TestRestProductControllerSuite(t *testing.T) {
// 	suite.Run(t, new(RestProductControllerSuite))
// }

// func (this *RestProductControllerSuite) SetupTest() {
// 	this.RequestHeaders = map[string]string{
// 		"Access-Control-Allow-Origin": "*",
// 		"Access-Control-Allow-Methods": "GET,POST,PUT,OPTIONS",
// 		"Access-Control-Allow-Credentials": "true",
// 		"Access-Control-Allow-Headers": "*",
// 	}
// }

// func (this *RestProductControllerSuite) TestProductCreateController() {
// 	RequestUrl := url.URL{Path: fmt.Sprintf("http://%s:%s/create/product/", ServerHost, ServerPort)}
// 	Client := http.Client()
// 	Request, Error := http.NewRequest("POST", RequestUrl.String(), nil)
// 	assert.Equal(this.T(), Error, nil)
// 	Response, Error := Client.Do(Request)

// 	assert.Equal(this.T(), Response.StatusCode, 200)
// }

// func (this *RestProductControllerSuite) TestProductUpdateController() {
// 	RequestUrl := url.URL(fmt.Sprintf("http://%s:%s/update/product/", ServerHost, ServerPort))
// 	Client := http.Client()
// 	Request, Error := http.NewRequest("PUT", RequestUrl.String())
// 	Response, Error := Client.Do()
// 	assert.Equal(this.T(), Response.StatusCode, 201)
// 	Response, Error := Client.Do(Request)
// 	if Response.StatusCode != 201 {assert.Error(this.T(), )}
// }

// func (this *RestProductControllerSuite) TestProductDeleteController() {
// 	RequestUrl := url.URL(fmt.Sprintf("http://%s:%s/delete/product/", ServerHost, ServerPort))
// 	Client := http.Client()
// 	Request, Error := http.NewRequest("DELETE", RequestUrl.String())
// 	if Error != nil {assert.Error(this.T(), Error, "Product Error.")}

// }
