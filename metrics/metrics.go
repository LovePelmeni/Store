package metrics

import (
	"flag"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	PROMETHEUS_HOST = os.Getenv("PROMETHEUS_HOST")
	PROMETHEUS_PORT = os.Getenv("PROMETHEUS_PORT")
	addr            = flag.String(fmt.Sprintf("%s", PROMETHEUS_HOST), ":"+PROMETHEUS_PORT,
		"Listens for a Prometheus Metric HTTP Requests...")
)
var (
	RegisteredUsers                     prometheus.Counter
	CustomersOnline                     prometheus.Gauge
	RequestTotalProcessingSummary       prometheus.Summary
	RequestTotalProcessingTimeHistogram prometheus.Histogram
	TotalPostRequests                   prometheus.Histogram
	TotalGetRequests                    prometheus.Histogram
	TotalDeleteRequests                 prometheus.Histogram
)

func main() {

	flag.Parse()

	RegisteredUsers = CustomersRegisteredMetric()
	CustomersOnline = CustomersOnlineMetric()

	RequestTotalProcessingTimeHistogram = RequestProcessingTimeHistogramMetric()
	RequestTotalProcessingSummary = RequestProcessingTimeSummaryMetric()

	TotalPostRequests = PostRequestsMetric()
	TotalDeleteRequests = DeleteRequestsMetric()
	TotalGetRequests = GetRequestsMetric()
}

func CustomersOnlineMetric() prometheus.Gauge {
	customersOnline := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "customers Online"})
	prometheus.MustRegister(customersOnline)
	return customersOnline
}

func CustomersRegisteredMetric() prometheus.Counter {
	customersRegistered := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "customers Registered"})
	prometheus.MustRegister(customersRegistered)
	return customersRegistered
}

func RequestProcessingTimeSummaryMetric() prometheus.Summary {
	summary := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "Request Processing Time Summary",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}})
	prometheus.MustRegister(summary)
	return summary
}

func RequestProcessingTimeHistogramMetric() prometheus.Histogram {
	requestProcessingTimeHistogramMs := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "Request Processing Time Histogram",
			Buckets: prometheus.LinearBuckets(0, 10, 20),
		})
	prometheus.MustRegister(requestProcessingTimeHistogramMs)
	return requestProcessingTimeHistogramMs
}

func PostRequestsMetric() prometheus.Histogram {
	totalPostRequests := prometheus.NewHistogram(
		prometheus.HistogramOpts{Name: "Total Post Requests."})
	prometheus.MustRegister()
	return totalPostRequests
}

func GetRequestsMetric() prometheus.Histogram {
	totalGetRequests := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "Total Get Requests Histogram",
			Buckets: prometheus.LinearBuckets(0, 0, 0)})
	return totalGetRequests
}

func DeleteRequestsMetric() prometheus.Histogram {
	totalDeleteRequests := prometheus.NewHistogram(
		prometheus.HistogramOpts{Name: "Total Delete Requests",
			Buckets: prometheus.LinearBuckets(0, 0, 0)})
	return totalDeleteRequests
}
