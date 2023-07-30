package elasticapmreceiver

import (
	"go.opentelemetry.io/collector/config/confighttp"
)

type Config struct {
	*confighttp.HTTPServerSettings `mapstructure:",squash"`
	EventsURLPath                  string `mapstructure:"events_url_path"`
	RUMEventsUrlPath               string `mapstructure:"rum_events_url_path"`
}

func (cfg *Config) Validate() error {
	return nil
}
