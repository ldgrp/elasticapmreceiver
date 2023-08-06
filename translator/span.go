package translator

import (
	"github.com/elastic/apm-data/model/modelpb"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func ConvertSpan(event *modelpb.APMEvent, dest ptrace.Span) {
	if event == nil {
		return
	}

	attrs := dest.Attributes()
	parseBaseEvent(event, attrs)

	if event.Span == nil {
		return
	}

	span := event.Span

	parseTrace(event.Trace, dest)
	dest.SetSpanID(ConvertSpanId(span.Id))
	if event.GetParentId() != "" {
		dest.SetParentSpanID(ConvertSpanId(event.ParentId))
	}
	dest.SetName(span.Name)
	dest.SetKind(ConvertSpanKind(span.Type))
	start, end := GetStartAndEndTimestamps(event.Timestamp, event.Event.Duration)
	if start != nil && end != nil {
		dest.SetStartTimestamp(*start)
		dest.SetEndTimestamp(*end)
	}

	// TODO: span.Message
	// TODO: span.Composite
	// TODO: span.DestinationService
	// TODO: span.Db
	// TODO: span.Sync
	// TODO: span.Kind
	// TODO: span.Action
	// TODO: span.Subtype
	// TODO: span.Stacktrace
	// TODO: span.Links
	// TODO: span.SelfTime
	// TODO: span.RepresentativeCount
}
