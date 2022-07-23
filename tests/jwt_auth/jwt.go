package test_jwt

import (
	"testing"

	// "github.com/LovePelmeni/OnlineStore/StoreService/authentication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JwtAuthenticationSuite struct {
	suite.Suite
	CustomerInfo map[string]string
}

func TestJwtAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(JwtAuthenticationSuite))
}

func (this *JwtAuthenticationSuite) SetupTest() {
	this.CustomerInfo = map[string]string{"Username": "", "Email": ""}
}

func (this *JwtAuthenticationSuite) TestJwtGenerate() {
	jwtToken := authentication.CreateJwtToken(this.CustomerInfo["Username"],
		this.CustomerInfo["Email"])
	assert.NotEmpty(this.T(), jwtToken)
}
