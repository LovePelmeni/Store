package payments 



import (
	paymentClient "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/proto"
	"fmt"
	"context"
	"time"
	"errors"
	"reflection"
	"strconv"
)


func getClient() *grpc.Client {}



type PaymentInterface interface {
	Pay(PaymentData map[string]interface{}) (bool, error)
	Refund(PaymentId string) (*Refund, error)
}


type PaymentCredentialsInterface interface {
	GetCredentials()
	Validate() error 
}


type Payment struct {}


type PaymentCredentials struct {}




func StartPaymentProcess() {

}



func (this *Payment) Pay(PaymentCredentials PaymentCredentialsInterface) error {


	if valid := PaymentCredentials.Validate(); valid != nil {
		DebugLogger.Println(
		fmt.Sprintf("Invalid Payment Credentials, Exception: %s"))
	}
	grpcClient := getClient()
	channel := make(chan struct{}, 10000)


	// Checking Validation of the Payment Form.

	for paymentElement, paymentValue := range reflection.Items(PaymentValidatedData){

		if paymentElement == "Price" && len(strings.Split(paymentValue, ".")) != 2{
			 // checks If Price has Appropriate Format "whatever_value.0000"
		
			paymentIntent += fmt.Sprintf(".%s", "0")
		}

		if hasPrefix := strings.HasPrefix(paymentElement, "P"); hasPrefix == true &&
		value, error := strconv.Atoi(paymentValue); error != nil {
			return exceptions.InvalidPaymentCredentials(
			[]string{paymentElement})
		}
	}


	group := sync.WaitGroup{}



	// Payment Intent Obtaining...

	go func (group *sync.WaitGroup, client *grpc.Client, channel chan struct{},
		PaymentValidatedData struct{ProductId string; PurchaserId string; Currency string; Price string}) {

			
		group.Add(1)

		PaymentRequestParams := paymentClient.PaymentIntentParams{
			ProductId: PaymentValidatedData.ProductId,
			PurchaserId: PaymentValidatedData.PurchaserId, 
			Currency: PaymentValidatedData.Currency,
			Price: PaymentValidatedData.Price,
		}

		RequestContext, CancelError := context.WithTimeout(context.Background(), 10 * time.Second)
		response, ResponseError := client.CreatePaymentIntent(RequestContext, PaymentRequestParams)

		if errors.Is(ResponseError, error){
		ErrorLogger.Println(fmt.Sprintf(
		"Payment Failed. Reason: %s", ResponseError), group.Abort())

		}else{
		// Sending Channel 
		channel <- struct{PaymentIntentId}{
		PaymentIntentId: response.PaymentIntentId}
		}
		
		defer CancelError()
		group.Done()
	}()




	// Processing Payment based on Received Payment Intent 
	go func (group *sync.WaitGroup, client *grpc.Client, channel chan struct{}){

		group.Add(1)
		select {

			case data := <- channel:
				// Receives Payment Intent and Processing Payment.
				data := <- channel 
				var decodedData struct{PaymentIntentId string}
				decodedDataError := json.Unmarshal(data, &decodedData)

				if decodedDataError != nil {ErrorLogger.Println("Invalid Data Format.")}

				CheckoutRequestParams := 
				RequestContext, CancelError := context.WithTimeout(context.Background(), 10 * time.Second)
				CheckoutContext := client.GetPaymentCheckout()


				group.Done()
				
			default:
				InfoLogger.Println("Operation Failed, Aborting Silently...")
				group.Abort()
		}
	}()



	// Processing Checkout ...
	go func (group *sync.WaitGroup, client *grpc.Client, channel chan struct{}) {

		group.Add(1)
		select {
			case data := <- channel != nil:
				var decodedData struct{PaymentIntentId string}
				decodedDataError := json.Unmarshal(data, &decodedData)
				group.Done()
			
			default:
				InfoLogger.Println("Operation Aborting...")
				group.Abort()
		}
	}()


	group.Wait()
	return true, nil 
}




func (this *Payment) Refund(CustomerId string, PaymentId string) {


	customer := models.Database.Table( // Receiving Purchaser Object Info.
	"customers").Where("id = ?", CustomerId).First(&models.Customer)


	if customer.Error != nil { // Check if it's exists...
	InfoLogger.Println("Customer Does Not Exist."),
    return false, exceptions.CustomerDoesNotExist()}


	// Initialize Wait Group Add Communication Channel For Goroutines...
	group := sync.WaitGroup{}
	channel := make(chan struct{})


	// Making Refund...
	go func (group *sync.WaitGroup, channel chan , refundParams string) {

		group.Add(1)

		// Making Refund Request...
		client := getClient()
		response, error := client.CreateRefund(RefundParams)
		group.Done()
	}()


	// Removing Purchased Products from Customer's Cart.
	go func (group *sync.WaitGroup, channel chan, customerId string) {

		group.Add(1)
		select {

			// Processing Success Refund
			case data := <- channel:
				Products := models.Database.Table(
				"carts").Where("owner.id = ?", customer)
				Products []

			default:
				// Processing Default Behaviourk, basically equivalents to the Failure...
		}
		group.Done()
	}
}





type RefundInteface interface {
	GetRefundInfo() (map[string]interface{}, error)
}


type CheckoutInterface interface {
	GetCheckoutInfo() (map[string]interface{}, error)
	GetRenderedCheckoutImage() ([]byte, error)
}




type Refund struct {}

type Checkout struct {}



