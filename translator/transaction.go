package translator

import (
	"github.com/elastic/apm-data/model/modelpb"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func ConvertTransaction(event *modelpb.APMEvent, dest ptrace.Span) {
	if event == nil {
		return
	}

	attrs := dest.Attributes()
	parseBaseEvent(event, attrs)

	if event.Transaction == nil {
		return
	}

	transaction := event.Transaction

	parseTrace(event.Trace, dest)
	dest.SetSpanID(ConvertSpanId(transaction.Id))
	if event.GetParentId() != "" {
		dest.SetParentSpanID(ConvertSpanId(event.ParentId))
	}
	dest.SetName(transaction.Name)
	dest.SetKind(ConvertSpanKind(transaction.Type))
	start, end := GetStartAndEndTimestamps(event.Timestamp, event.GetEvent().GetDuration())
	if start != nil && end != nil {
		dest.SetStartTimestamp(*start)
		dest.SetEndTimestamp(*end)
	}

	// TODO: transaction.SpanCount
	// TODO: transaction.UserExperience
	// TODO: transaction.Custom
	// TODO: transaction.Marks
	// TODO: transaction.Message
	// TODO: transaction.Result
	// TODO: transaction.DroppedSpansStats
	// TODO: transaction.DurationSummary
	// TODO: transaction.RepresentativeCount
	// TODO: transaction.Sampled
	// TODO: transaction.Root
}
