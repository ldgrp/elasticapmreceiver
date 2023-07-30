package elasticapmreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	typeStr                 = "elasticapm"
	defaultHTTPEndpoint     = "localhost:8200"
	defaultEventsURLPath    = "/intake/v2/events"
	defaultRUMEventsURLPath = "/intake/v2/rum/events"
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithTraces(createTraces, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		HTTPServerSettings: &confighttp.HTTPServerSettings{
			Endpoint: defaultHTTPEndpoint,
		},
		EventsURLPath:    defaultEventsURLPath,
		RUMEventsUrlPath: defaultEventsURLPath,
	}
}

func createTraces(
	_ context.Context,
	params receiver.CreateSettings,
	baseCfg component.Config,
	nextConsumer consumer.Traces,
) (receiver.Traces, error) {
	cfg := baseCfg.(*Config)
	r, err := newElasticAPMReceiver(cfg, params)

	if err != nil {
		return nil, err
	}

	if err = r.registerTraceConsumer(nextConsumer); err != nil {
		return nil, err
	}

	return r, nil
}
