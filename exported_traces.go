package mongodbexporter

import (
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type isAnyValue_Value interface {
	isAnyValue_Value()
	MarshalTo([]byte) (int, error)
	Size() int
}

type AnyValue struct {
	Value isAnyValue_Value
}

type KeyValue struct {
	Key   string
	Value AnyValue
}

type Span_Event struct {
	TimeUnixNano           uint64
	Name                   string
	Attributes             map[string]any
	DroppedAttributesCount uint32
}

type Span_Link struct {
	TraceId                [16]byte
	SpanId                 [8]byte
	TraceState             string
	Attributes             map[string]any
	DroppedAttributesCount uint32
}

type ExportedSpan struct {
	TraceId                [16]byte
	SpanId                 [8]byte
	TraceState             string
	ParentSpanId           [8]byte
	Name                   string
	Kind                   string
	StartTimeUnixNano      uint64
	EndTimeUnixNano        uint64
	Attributes             map[string]any
	DroppedAttributesCount uint32
	Events                 []*Span_Event
	DroppedEventsCount     uint32
	Links                  []*Span_Link
	DroppedLinksCount      uint32
	Status                 Status
}

type Status struct {
	Message string
	Code    int32
}

type ExportedResource struct {
	Attributes             map[string]any
	DroppedAttributesCount uint32
	SchemaUrl              string
}

type ExportedTraces struct {
	Resource      ExportedResource
	ExportedSpans []ExportedSpan
}

type TraceDoc struct {
	ExportedTraces []ExportedTraces
}

func newTraceDoc(trace ptrace.Traces) TraceDoc {
	return TraceDoc{
		ExportedTraces: newExportedTraces(trace),
	}
}

func newExportedTraces(trace ptrace.Traces) []ExportedTraces {
	resourceSpans := trace.ResourceSpans()
	exportedTraces := make([]ExportedTraces, resourceSpans.Len())
	for i := 0; i < resourceSpans.Len(); i++ {
		rs := resourceSpans.At(i)
		resource := rs.Resource()
		scopeSpansSlice := rs.ScopeSpans()
		exportedTraces[i] = ExportedTraces{
			Resource: ExportedResource{
				Attributes:             resource.Attributes().AsRaw(),
				DroppedAttributesCount: resource.DroppedAttributesCount(),
			},
		}
		for j := 0; j < scopeSpansSlice.Len(); j++ {
			scopeSpans := scopeSpansSlice.At(j)
			//scope := scopeSpans.Scope()
			spansSlice := scopeSpans.Spans()
			if spansSlice.Len() == 0 {
				continue
			}
			exportedSpans := make([]ExportedSpan, spansSlice.Len())
			for k := 0; k < spansSlice.Len(); k++ {
				exportedSpans[k] = newExportedSpan(spansSlice.At(k))
			}
			exportedTraces[i].ExportedSpans = exportedSpans
		}
	}
	return exportedTraces
}

func newExportedSpan(span ptrace.Span) ExportedSpan {
	return ExportedSpan{
		TraceId:                span.TraceID(),
		SpanId:                 span.SpanID(),
		TraceState:             span.TraceState().AsRaw(),
		ParentSpanId:           span.ParentSpanID(),
		Name:                   span.Name(),
		Kind:                   span.Kind().String(),
		StartTimeUnixNano:      uint64(span.StartTimestamp().AsTime().Unix()),
		EndTimeUnixNano:        uint64(span.EndTimestamp().AsTime().Unix()),
		Attributes:             span.Attributes().AsRaw(),
		DroppedAttributesCount: span.DroppedAttributesCount(),
		DroppedEventsCount:     span.DroppedEventsCount(),
		DroppedLinksCount:      span.DroppedLinksCount(),
		Events:                 newExportedEvents(span),
		Links:                  newExportedLinks(span),
		Status: Status{
			Message: span.Status().Message(),
			Code:    int32(span.Status().Code()),
		},
	}
}

func newExportedEvents(span ptrace.Span) []*Span_Event {
	spanEventsSlice := span.Events()
	spanExportedEvents := make([]*Span_Event, spanEventsSlice.Len())
	for i := 0; i < spanEventsSlice.Len(); i++ {
		spanExportedEvents[i] = &Span_Event{
			TimeUnixNano:           uint64(spanEventsSlice.At(i).Timestamp().AsTime().Unix()),
			Name:                   spanEventsSlice.At(i).Name(),
			Attributes:             spanEventsSlice.At(i).Attributes().AsRaw(),
			DroppedAttributesCount: spanEventsSlice.At(i).DroppedAttributesCount(),
		}
	}
	return spanExportedEvents
}

func newExportedLinks(span ptrace.Span) []*Span_Link {
	spanLinksSlice := span.Links()
	spanExportedLinks := make([]*Span_Link, spanLinksSlice.Len())
	for i := 0; i < spanLinksSlice.Len(); i++ {
		spanExportedLinks[i] = &Span_Link{
			TraceId:                spanLinksSlice.At(i).TraceID(),
			SpanId:                 spanLinksSlice.At(i).SpanID(),
			TraceState:             spanLinksSlice.At(i).TraceState().AsRaw(),
			Attributes:             spanLinksSlice.At(i).Attributes().AsRaw(),
			DroppedAttributesCount: spanLinksSlice.At(i).DroppedAttributesCount(),
		}
	}
	return spanExportedLinks
}
