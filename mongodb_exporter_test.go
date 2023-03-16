package mongodbexporter

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"

	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
)

const (
	defaultRequestDataEnvelopeName          = "Microsoft.ApplicationInsights.Request"
	defaultRemoteDependencyDataEnvelopeName = "Microsoft.ApplicationInsights.RemoteDependency"
	defaultMessageDataEnvelopeName          = "Microsoft.ApplicationInsights.Message"
	defaultExceptionDataEnvelopeName        = "Microsoft.ApplicationInsights.Exception"
	defaultServiceName                      = "foo"
	defaultServiceNamespace                 = "ns1"
	defaultServiceInstance                  = "112345"
	defaultScopeName                        = "myinstrumentationlib"
	defaultScopeVersion                     = "1.0"
	defaultHTTPMethod                       = "GET"
	defaultHTTPServerSpanName               = "/bar"
	defaultHTTPClientSpanName               = defaultHTTPMethod
	defaultHTTPStatusCode                   = 200
	defaultRPCSystem                        = "grpc"
	defaultRPCSpanName                      = "com.example.ExampleRmiService/exampleMethod"
	defaultRPCStatusCode                    = 0
	defaultDBSystem                         = "mssql"
	defaultDBName                           = "adventureworks"
	defaultDBSpanName                       = "astoredproc"
	defaultDBStatement                      = "exec astoredproc1"
	defaultDBOperation                      = "exec astoredproc2"
	defaultMessagingSpanName                = "MyQueue"
	defaultMessagingSystem                  = "kafka"
	defaultMessagingDestination             = "MyQueue"
	defaultMessagingURL                     = "https://queue.amazonaws.com/80398EXAMPLE/MyQueue"
	defaultInternalSpanName                 = "MethodX"

	connString       = "mongodb://myUserAdmin:myUserAdmin@172.21.109.159:27017"
	database         = "OtelDB"
	tracesCollection = "Requests"
)

var (
	defaultTraceID                = [16]byte{35, 191, 77, 229, 162, 242, 217, 75, 148, 170, 81, 99, 227, 163, 145, 25}
	defaultTraceIDAsHex           = fmt.Sprintf("%02x", defaultTraceID)
	defaultSpanID                 = [8]byte{35, 191, 77, 229, 162, 242, 217, 76}
	defaultSpanIDAsHex            = fmt.Sprintf("%02x", defaultSpanID)
	defaultParentSpanID           = [8]byte{35, 191, 77, 229, 162, 242, 217, 77}
	defaultParentSpanIDAsHex      = fmt.Sprintf("%02x", defaultParentSpanID)
	defaultSpanStartTime          = pcommon.Timestamp(0)
	defaultSpanEndTme             = pcommon.Timestamp(60000000000)
	defaultSpanEventTime          = pcommon.Timestamp(0)
	defaultHTTPStatusCodeAsString = strconv.FormatInt(defaultHTTPStatusCode, 10)
	defaultRPCStatusCodeAsString  = strconv.FormatInt(defaultRPCStatusCode, 10)

	// Same as RPC codes?
	defaultDatabaseStatusCodeAsString  = strconv.FormatInt(defaultRPCStatusCode, 10)
	defaultMessagingStatusCodeAsString = strconv.FormatInt(defaultRPCStatusCode, 10)

	// Required attribute for any HTTP Span
	requiredHTTPAttributes = map[string]interface{}{
		conventions.AttributeHTTPMethod: defaultHTTPMethod,
	}

	// Required attribute for any RPC Span
	requiredRPCAttributes = map[string]interface{}{
		conventions.AttributeRPCSystem: defaultRPCSystem,
	}

	requiredDatabaseAttributes = map[string]interface{}{
		conventions.AttributeDBSystem: defaultDBSystem,
		conventions.AttributeDBName:   defaultDBName,
	}

	requiredMessagingAttributes = map[string]interface{}{
		conventions.AttributeMessagingSystem:      defaultMessagingSystem,
		conventions.AttributeMessagingDestination: defaultMessagingDestination,
	}

	defaultResource               = getResource()
	defaultInstrumentationLibrary = getScope()
)

func TestMongoDbConnection(t *testing.T) {
	config := &Config{
		ConnectionURI:    connString,
		Database:         database,
		CollectionTraces: tracesCollection,
	}
	ctx := context.TODO()
	mongoDbExporter, err := newMongoDbExporter(ctx, config)
	if err != nil {
		t.Error(err.Error())
	}
	//defer mongoDbExporter.Shutdown(ctx)
	resource := getResource()
	scope := getScope()
	span := getDefaultHTTPServerSpan()
	traces := ptrace.NewTraces()
	rs := traces.ResourceSpans().AppendEmpty()
	r := rs.Resource()
	resource.CopyTo(r)
	ilss := rs.ScopeSpans().AppendEmpty()
	scope.CopyTo(ilss.Scope())
	span.CopyTo(ilss.Spans().AppendEmpty())
	//rs := traces.ResourceSpans().AppendEmpty()
	err = mongoDbExporter.ConsumeTraces(ctx, traces)
	if err != nil {
		t.Error(err.Error())
	}

}

// Returns a default Resource
func getResource() pcommon.Resource {
	r := pcommon.NewResource()
	r.Attributes().PutStr(conventions.AttributeServiceName, defaultServiceName)
	r.Attributes().PutStr(conventions.AttributeServiceNamespace, defaultServiceNamespace)
	r.Attributes().PutStr(conventions.AttributeServiceInstanceID, defaultServiceInstance)
	return r
}

// Returns a default instrumentation library
func getScope() pcommon.InstrumentationScope {
	il := pcommon.NewInstrumentationScope()
	il.SetName(defaultScopeName)
	il.SetVersion(defaultScopeVersion)
	return il
}

func getDefaultHTTPServerSpan() ptrace.Span {
	return getServerSpan(
		defaultHTTPServerSpanName,
		requiredHTTPAttributes)
}

// Returns a default server span
func getServerSpan(spanName string, initialAttributes map[string]interface{}) ptrace.Span {
	return getSpan(spanName, ptrace.SpanKindServer, initialAttributes)
}

/*
The remainder of these methods are for building up test assets
*/
func getSpan(spanName string, spanKind ptrace.SpanKind, initialAttributes map[string]interface{}) ptrace.Span {
	span := ptrace.NewSpan()
	span.SetTraceID(defaultTraceID)
	span.SetSpanID(defaultSpanID)
	span.SetParentSpanID(defaultParentSpanID)
	span.SetName(spanName)
	span.SetKind(spanKind)
	span.SetStartTimestamp(defaultSpanStartTime)
	span.SetEndTimestamp(defaultSpanEndTme)
	//nolint:errcheck
	span.Attributes().FromRaw(initialAttributes)
	return span
}
