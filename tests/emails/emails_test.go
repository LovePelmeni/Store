package test_emails

import (
	"errors"
	"testing"

	"github.com/LovePelmeni/OnlineStore/EmailService/emails/proto/grpcControllers"
	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/emails"
	mocked_emails "github.com/LovePelmeni/OnlineStore/StoreService/mocks/emails"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EmailIntegrationSuite struct {
	// Integration Test Case for grpc Controller
	//  that provides comminication with Email Service.
	suite.Suite
	Controller                *gomock.Controller
	EmailSenderController     *emails.GrpcEmailSender
	MockedEmailClientInteface *mocked_emails.MockNewEmailClient
}

func (this *EmailIntegrationSuite) SetupTest() {

	this.Controller = gomock.NewController(this.T())

	this.MockedEmailClientInteface = mocked_emails.NewMockNewEmailClient(this.Controller)

	this.EmailSenderController = &emails.GrpcEmailSender{
		Client: this.MockedEmailClientInteface}
}

func (this *EmailIntegrationSuite) TeardownTest() {
	this.Controller.Finish()
}

func (this *EmailIntegrationSuite) TestSendEmail() {

	customerEmail := "kirklimishin@gmail.com"
	message := "Test Message"

	MockedResponse := grpcControllers.EmailResponse{Delivered: true}

	this.MockedEmailClientInteface.EXPECT().SendEmail().Return(MockedResponse, nil).Times(1)
	this.NoError(errors.New("Unexpected Error"))

	ActualResponse, error := this.EmailSenderController.SendDefaultEmail(customerEmail, message)
	ExpectedResponse, error := true, nil
	assert.Equal(this.T(), ActualResponse, ExpectedResponse)
	assert.Equal(this.T(), error, nil)
}

func (this *EmailIntegrationSuite) TestFailSendEmail() {

	Exception := errors.New("Failure Error")
	customerEmail := "test@gmail.com"
	message := "Test Message"

	this.MockedEmailClientInteface.EXPECT().SendOrderEmail().Return(nil, Exception).Times(1)
	ActualResponse, error := this.EmailSenderController.SendDefaultEmail(customerEmail, message)

	assert.Equal(this.T(), ActualResponse, nil, "Response Should Equal To None.")
	assert.Equal(this.T(), error, Exception, "Error Should Be Not None.")
}

func TestEmailInterfaceSuite(t *testing.T) {
	suite.Run(t, new(EmailIntegrationSuite))
}
