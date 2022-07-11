package payments 



import (
	paymentClient "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/proto"
	"fmt"
	"context"
	"time"
	"errors"
	"reflection"
	"strconv"
	"sync"
)


func getClient() *grpc.Client {}



type PaymentInterface interface {
	Pay(PaymentData map[string]string) (bool, error)
	Refund(PaymentId string) (*Refund, error)
}


type PaymentInfoCredentials struct {

	mutex sync.RWMutex 

	PaymentSessionId string `json:"PaymentSessionId"`		
	ProductId string 	`json:"ProductId"`
	CustomerId string 	`json:"CustomerId"`
	Currency string 	`json:"Currency"`
	Price string 		`json:"Price"`
}

func (this *PaymentInfoCredentials) Validate() (bool, error){
	return true, nil 
}

type PaymentCheckoutInfo struct {
	sync.RWMutex 
	
	PaymentId string `json:"PaymentId"`
	PurchaserId string `json:"PurchaserId"`
	Products []models.Products `json:"Products"`
	TotalPrice string `json:"TotalPrice"`
	Currency string `json:"Currency"`
	CreatedAt time.Time `json:"CreatedAt"`
}


func (this *Payment) Pay(Client *grpc.Client, PaymentCredentials PaymentInfoCredentials)  error {


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

	go func (group *sync.WaitGroup, client *grpc.Client, channel chan struct{}, PaymentValidatedData PaymentInfoCredentials) {

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




		// Processing Payment based on Received Payment Intent 
		go func (group *sync.WaitGroup, client *grpc.Client){

			group.Add(1)
			select {

				case data := <- channel:
					// Receives Payment Intent and Processing Payment.
					data := <- channel 
					var decodedData struct{PaymentId string}
					decodedDataError := json.Unmarshal(data, &decodedData)

					if decodedDataError != nil {ErrorLogger.Println("Invalid Data Format.")}


					checkoutParams := CheckoutRequestParams{PaymentId: decodedData.PaymentId}
					RequestContext, CancelError := context.WithTimeout(context.Background(), 10 * time.Second)
					CheckoutResponseconst, error := client.GetPaymentCheckout(checkoutParams)


					if error != nil {ErrorLogger.Println(
					fmt.Sprintf("Failed to Receive Checkout Info about the Payment with ID: %s",
				    &checkoutParams.PaymentId))} else  {


								// Processing Sending Email... About the Purchased Order.
						go func (group *sync.WaitGroup, client *grpc.Client, channel chan struct{}, PaymentCredentials *Payment) {
								group.Add(1)
								
						}()
					}
					

					group.Done()
					defer CancelError()
					
				default:
					InfoLogger.Println("Operation Failed, Aborting Silently...")
					group.Abort()
			}
		}()

	group.Wait()
	return true, nil 

	}()
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



