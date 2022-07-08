package customers 


import (
	"github.com/LovePelmeni/OnlineStore/StoreService/customers"
	"os"
	"fmt"
)

//go:generate mockgen -destination=mocks/customer.go --build_flags=--mod=mod . CustomerInterface
type CustomerInterface interface {
	// Interface for Managing Customer Model 
	CreateCustomer(customerData map[string]interface{}) (bool, error)
	UpdateCustomer(customerId string, UpdatedData ...map[string]interface{}) (bool, error)
	DeleteCustomer(customerId string) (bool, error)
}

type CustomerStruct struct {}