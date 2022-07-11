package test_models

import (
	"testing"

	mocked_models "github.com/LovePelmeni/OnlineStore/StoreService/mocks/models"
	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)


func TestModelsCreate(t *testing.T) {
	TestModelsData := []struct {
		
		ModelTestFunc func(t *testing.T) 
	}{
		
	func(t *testing.T){models.Product},
	func(t *testing.T){models.Cart}, 
	func (t *testing.T) {models.Customer},
	}

	for _, test := range TestModelsData{
		suite.Run(t, test.ModelTestFunc)
	}
}

func TestModelUpdate(t *testing.T) {
	TestModelsData := []struct {
		TestFunc func(t *testing.T)
		UpdatedData map[string]interface{}
	}{}
}

func TestModelsDelete(t *testing.T){
	TestModelsData := []struct {
		TestFunc func(t *testing.T) 
		ObjId string
	}
}

