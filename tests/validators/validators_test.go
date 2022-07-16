package test_validators

import (
	"errors"
	"testing"

	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ModelValidatorSuite struct {
	suite.Suite
	ModelValidators []models.BaseModelValidator
	TestValidData   map[string]map[string]string
}

func (this *ModelValidatorSuite) SetupTest() {

	newCustomerValidator := models.NewCustomerModelValidator()
	newProductValidator := models.NewProductModelValidator()

	this.ModelValidators = []models.BaseModelValidator{
		newCustomerValidator, newProductValidator}

	this.TestValidData = map[string]map[string]string{

		"CustomerData": {"Username": "Test-Customer", "Email": "Email@gmail.com", "Password": "Test-Password"}, // Valid Data for Customer Model.

		"ProductData": {"ProductName": "ProductName", "ProductDescription": "Some-Product-Description", "OwnerEmail": "OwnerEmail@gmail.com"}, // Valid data for Product Model.

		"CartData": {"Products": "nil"}, // Valid Data for Cart Model.
	}
}

func TestModelValidatorSuite(t *testing.T) {
	suite.Run(t, new(ModelValidatorSuite))
}

func (this *ModelValidatorSuite) TestValidators() {

	Passed := testing.RunTests(func(pat string, s string) (bool, error) {
		return true, nil
	}, []testing.InternalTest{

		{"Customer Validator", func(t *testing.T) {
			Valid, Errors := this.ModelValidators[0].Validate(this.TestValidData["CustomerData"])
			if Valid == nil || len(Errors) != 0 {
				assert.Error(t,
					errors.New("Invalid Data for Customer Validator"), "Customer Validator Error")
			}
		}},

		{"Product Validator", func(t *testing.T) {
			Valid, Errors := this.ModelValidators[1].Validate(this.TestValidData["ProductData"])
			if Valid == nil || len(Errors) != 0 {
				assert.Error(t, errors.New("Product Validator Error"))
			}
		}},

		{"Cart Validator", func(t *testing.T) {
			Valid, Errors := this.ModelValidators[2].Validate(this.TestValidData["CartData"])
			if Valid == nil || len(Errors) != 0 {
				assert.Error(t,
					errors.New("Cart Validator Error"))
			}
		}},
	})
	assert.Equal(this.T(), Passed, true, "Validators Failed.") // Suppose that all validators has passed all the data successfully...
}

func (this *ModelValidatorSuite) TestFailValidators() {

	InvalidModelData := map[string]map[string]string{

		"CustomerData": {"": ""},

		"ProductData": {"": ""},

		"CartData": {"": ""},
	}

	Passed := testing.RunTests(func(pat string, s string) (bool, error) {
		return true, nil
	}, []testing.InternalTest{

		{"Customer Validator", func(t *testing.T) {
			Valid, Errors := this.ModelValidators[0].Validate(InvalidModelData["CustomerData"])
			if Valid == nil || len(Errors) != 0 {
				assert.Error(t,
					errors.New("Invalid Data for Customer Validator"), "Customer Validator Error")
			}
		}},

		{"Product Validator", func(t *testing.T) {
			Valid, Errors := this.ModelValidators[1].Validate(InvalidModelData["CustomerData"])
			if Valid == nil || len(Errors) != 0 {
				assert.Error(t, errors.New("Product Validator Error"))
			}
		}},

		{"Cart Validator", func(t *testing.T) {
			Valid, Errors := this.ModelValidators[2].Validate(InvalidModelData["CustomerData"])
			if Valid == nil || len(Errors) != 0 {
				assert.Error(t,
					errors.New("Cart Validator Error"))
			}
		}},
	})
	assert.Equal(this.T(), Passed, false, "Validators should fail")
}
