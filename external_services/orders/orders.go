package orders

import (
	"time"

	"github.com/"
	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/go-kit/kit/transport/amqp"
	amqp_server "github.com/streadway/amqp"
)

var (
	RabbitmqHost = os.Getenv("RABBITMQ_HOST")
	RabbitmqPort = os.Getenv("RABBITMQ_PORT")
	RabbitmqUser = os.Getenv("RABBITMQ_USER")
	RabbitmqPassword = os.Getenv("RABBITMQ_PASSWORD")
	RabbitmqVhost = os.Getenv("RABBITMQ_VHOST")
)

type RabbitmqTransportInterface interface{

	Connect(credentials *RabbitmqCredentials) *amqp_server.Channel 
	Disconnect() (bool, error)

	CreateQueue(QueueName string) (bool, error) 
	RemoteQueue(QueueName string) (bool, error)

	PublishEvent(BodyData map[string]interface{}, RoutingKey string) (bool, error)
	ListenEvent(ResponseHandler func(Queue string, 
	Method string, Props map[string]interface{}, Body []byte)) (bool, error)
}
type RabbitmqTransport struct {

}
func (this *RabbitmqTransport) Connect(credentials *RabbitmqCredentials) {
	// Connection Method.. Uses AMQP Protocol for Connection.
	ConnectionURI := url.URL(fmt.Sprintf("amqp://%s:%s@%s:%s/%s", RabbitmqUser, RabbitmqPassword,
    RabbitmqHost, RabbitmqPort, RabbitmqVhost))
	ConnectionChannel, error := amqp_server.Dial(ConnectionURI)

	if error != nil {ErrorLogger.Println("Connection Failed, Check If Server Running.")}
	Channel, err := ConnectionChannel.Channel()
	if err != nil {ErrorLogger.Println("Failed To Initialize Connection Channel.")}
	return Channel 
}

func NewOrder(credentials *OrderStruct) *RabbitmqTransportInterface {
	return RabbitmqTransportInterface{}{}
}


type OrderInterface interface {
	// Interface 
	CreateOrder(credentials map[string]string) (bool, error)
	CancelOrder(OrderId string) (bool, error)
}

type OrderStruct struct {
	OrderInfo struct{OrderName string; PurchaserId string; 
	Goods []models.Product; TotalPrice string}
	CreatedAt time.Time 
}


func ProcessOrder(order *OrderInterface) error {
	// Method That Processing Order Initialization....
}

func ProcessCancelOrder(order *OrderInterface) error {

}