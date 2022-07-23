package services

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
)

var (
	PROMETHEUS_HOST = os.Getenv("PROMETHEUS_HOST")
	PROMETHEUS_PORT = os.Getenv("PROMETHEUS_PORT")
)

func init() {
	LogFile, Error := os.OpenFile("Metrics_Services.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if Error != nil {
		panic(Error)
	}
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	WarningLogger = log.New(LogFile, "WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
}

type MetricServerInterface interface {
	// Interface for Metrics Server... basically represents communication layer for Prometheus..
	GetUrl() string // returns Metrics Server URl
	PushMetrics(Metric prometheus.Metric, Job string) (bool, error)
}

type MetricServer struct {
	Host string
	Port string
}

func NewMetricServer(Host string, Port string) *MetricServer {
	return &MetricServer{Host: Host, Port: Port}
}

func (this *MetricServer) GetUrl() string {
	return "http://" + this.Host + ":" + this.Port
}

func (this *MetricServer) PushMetrics(Metric prometheus.Metric, Job string) (bool, error) {
	PushedError := push.New(this.GetUrl(), Job).Push()
	if PushedError != nil {
		DebugLogger.Println(
			"Failed to Push Metric To Prometheus")
		return false, PushedError
	} else {
		DebugLogger.Println("Metric Has been Pushed To Prometheus.")
		return true, nil
	}
}

type CustomerOnlineMetricService struct {
	CustomerOnlineMetricCounter prometheus.Counter
	MetricServer                MetricServerInterface
}

func NewCustomerOnlineMetricService() *CustomerOnlineMetricService {
	MetricServer := NewMetricServer(PROMETHEUS_HOST, PROMETHEUS_PORT)
	return &CustomerOnlineMetricService{
		MetricServer: MetricServer,
	}
}
func (this *CustomerOnlineMetricService) ProcessCustomersOnlineMetrics(DataMetric prometheus.Metric) (bool, error) {
	return this.MetricServer.PushMetrics(
		DataMetric, "Customer Online.")
}

type CustomerRegisteredMetricService struct {
	MetricServer MetricServerInterface
}

func (this *CustomerRegisteredMetricService) ProcessCustomerRegisteredMetrics(DataMetric prometheus.Metric) (bool, error) {
	return this.MetricServer.PushMetrics(
		DataMetric, "Customer Registered.")
}

func NewCustomerRegisteredMetricService() *CustomerRegisteredMetricService {
	MetricServer := NewMetricServer(PROMETHEUS_HOST, PROMETHEUS_PORT)
	return &CustomerRegisteredMetricService{
		MetricServer: MetricServer,
	}
}

type RequestsLatencyMetricsSummaryService struct {
	MetricServer MetricServerInterface
}

func NewRequestsLatencyMetricsService() *RequestsLatencyMetricsSummaryService {
	return &RequestsLatencyMetricsSummaryService{
		MetricServer: NewMetricServer(PROMETHEUS_HOST, PROMETHEUS_PORT),
	}
}

func (this *RequestsLatencyMetricsSummaryService) ProcessRequestsLatencySummaryMetrics(DataMetric prometheus.Metric) (bool, error) {

	return this.MetricServer.PushMetrics(
		DataMetric,
		"Requests Time Latency Summary.")
}

type RequestsLatencyMetricsHistogramService struct {
	TotalRequestsTimeProcessingHistogram prometheus.Histogram
	MetricServer                         MetricServerInterface
}

func NewRequestsLatencyMetricsHistogramService() *RequestsLatencyMetricsHistogramService {

	return &RequestsLatencyMetricsHistogramService{
		MetricServer: NewMetricServer(PROMETHEUS_HOST, PROMETHEUS_PORT),
	}
}

func (this *RequestsLatencyMetricsHistogramService) ProcessRequestsLatencyHistogramMetric(DataMetric prometheus.Metric) (bool, error) {

	return this.MetricServer.PushMetrics(
		DataMetric,
		"Requests Time Latency Histogram.")
}
