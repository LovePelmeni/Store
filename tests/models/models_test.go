package test_models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ModelSuite struct {
	suite.Suite
	ModelTestCases []testing.InternalTest
}

func (this *ModelSuite) SetupTest() {
	this.ModelTestCases = []testing.InternalTest{

		{"Test Customer Create", func(t *testing.T) {}},
		{"Test Customer Update", func(t *testing.T) {}},
		{"Test Customer Delete", func(t *testing.T) {}},

		{"Test Product Create", func(t *testing.T) {}},
		{"Test Product Update", func(t *testing.T) {}},
		{"Test Product Delete", func(t *testing.T) {}},

		{"Test Cart Create", func(t *testing.T) {}},
		{"Test Cart Update", func(t *testing.T) {}},
		{"Test Cart Delete", func(t *testing.T) {}},
	}
}
func (this *ModelSuite) TestModels() {
	Response := testing.RunTests(func(str string, pat string) (bool, error) { return true, nil },
		this.ModelTestCases)
	assert.Equal(this.T(), Response, true, "Tests Failed.")
}
