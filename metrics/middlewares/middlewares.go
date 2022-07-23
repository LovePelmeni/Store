package metrics_middlewares

import (
	"net/http/httptrace"

	metrics "github.com/LovePelmeni/Store/metrics"
	metrics_services "github.com/LovePelmeni/Store/metrics/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

var (
	RequestSummaryService     = metrics_services.RequestsLatencyMetricsSummaryService{}
	RequestHistogramService   = metrics_services.RequestsLatencyMetricsHistogramService{}
	CustomerOnlineService     = metrics_services.CustomerOnlineMetricService{}
	CustomerRegisteredService = metrics_services.CustomerRegisteredMetricService{}
)

type RequestLatencyParser struct{}

func NewRequestLatencyParser() *RequestLatencyParser {
	return &RequestLatencyParser{}
}

func (this *RequestLatencyParser) TraceHttpContext(
	RequestContext context.Context) (context.Context, *httptrace.ClientTrace) {
	// Logic Of calculating Request Latency...

	trace := httptrace.ContextClientTrace(RequestContext)
	traceContext := httptrace.WithClientTrace(RequestContext, trace)
	return traceContext, trace
}

func RequestProcessingTimeSummaryMiddleware(Context *gin.Context) gin.HandlerFunc {
	return gin.HandlerFunc(func(context *gin.Context) {

		RequestMetaParser := NewRequestLatencyParser()
		TracedContext, LatencyResult := RequestMetaParser.TraceHttpContext(context)
		context.Request.WithContext(TracedContext)

		DataMetric := metrics.RequestProcessingTimeSummaryMetric()
		DataMetric.Observe(LatencyResult)
		go RequestSummaryService.ProcessRequestsLatencySummaryMetrics(DataMetric)
		context.Next()
	})
}

func RequestTimeProcessingMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(context *gin.Context) {

		RequestMetaParser := NewRequestLatencyParser()
		DataMetric := metrics.RequestProcessingTimeHistogramMetric()
		TracedContext, LatencyResult := RequestMetaParser.TraceHttpContext(context)
		context.Request.WithContext(TracedContext)

		DataMetric.Observe(LatencyResult)
		go RequestHistogramService.ProcessRequestsLatencyHistogramMetric(DataMetric)
		context.Next()
	})
}

func CustomerRegisteredMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(context *gin.Context) {

		DataMetric := metrics.CustomersRegisteredMetric()
		DataMetric.Inc()
		go CustomerRegisteredService.ProcessCustomerRegisteredMetrics(DataMetric)
		context.Next()
	})
}

func CustomerOnlineMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(context *gin.Context) {

		DataMetric := metrics.CustomersOnlineMetric()
		DataMetric.Inc()
		go CustomerOnlineService.ProcessCustomersOnlineMetrics(DataMetric)
		context.Next()
	})
}
