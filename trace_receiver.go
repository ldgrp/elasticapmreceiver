package elasticapmreceiver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

type elasticapmReceiver struct {
	cfg        *Config
	serverHTTP *http.Server
	httpMux    *http.ServeMux
	shutdownWG sync.WaitGroup

	nextConsumer  consumer.Traces
	traceReceiver *obsreport.Receiver
	settings      receiver.CreateSettings
}

func newElasticAPMReceiver(cfg *Config, settings receiver.CreateSettings) (*elasticapmReceiver, error) {
	r := &elasticapmReceiver{
		cfg:      cfg,
		settings: settings,
		httpMux:  http.NewServeMux(),
	}

	var err error
	r.traceReceiver, err = obsreport.NewReceiver(obsreport.ReceiverSettings{
		ReceiverID:             settings.ID,
		Transport:              "http",
		ReceiverCreateSettings: settings,
	})

	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *elasticapmReceiver) startHTTPServer(cfg *confighttp.HTTPServerSettings, host component.Host) error {
	r.settings.Logger.Info("Starting HTTP server", zap.String("endpoint", cfg.Endpoint))
	var hln net.Listener
	hln, err := cfg.ToListener()
	if err != nil {
		return err
	}

	r.shutdownWG.Add(1)
	go func() {
		defer r.shutdownWG.Done()

		if errHTTP := r.serverHTTP.Serve(hln); errHTTP != http.ErrServerClosed {
			host.ReportFatalError(errHTTP)
		}
	}()
	return nil
}

func (r *elasticapmReceiver) Start(ctx context.Context, host component.Host) error {
	var err error
	r.serverHTTP, err = r.cfg.HTTPServerSettings.ToServer(
		host,
		r.settings.TelemetrySettings,
		r.httpMux,
	)

	if err != nil {
		return err
	}

	err = r.startHTTPServer(r.cfg.HTTPServerSettings, host)
	return err
}

func (r *elasticapmReceiver) Shutdown(ctx context.Context) error {
	var err error

	if r.serverHTTP != nil {
		err = r.serverHTTP.Shutdown(ctx)
	}

	r.shutdownWG.Wait()
	return err
}

func (r *elasticapmReceiver) registerTraceConsumer(nextConsumer consumer.Traces) error {
	if nextConsumer == nil {
		return component.ErrNilNextConsumer
	}

	if r.httpMux != nil {
		r.httpMux.HandleFunc(r.cfg.EventsURLPath, r.handleEvents)
		r.httpMux.HandleFunc(r.cfg.RUMEventsUrlPath, r.handleEvents)
	}

	return nil
}

func (r *elasticapmReceiver) handleEvents(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Only POST requests are supported")
		return
	}

	switch req.Header.Get("Content-Type") {
	// Only parse ndjson
	case "application/x-ndjson":
		handleTraces(w, req, r.traceReceiver)
	default:
		writeError(w, http.StatusUnsupportedMediaType, "Only application/ndjson is supported")
		return
	}
}

func handleTraces(w http.ResponseWriter, req *http.Request, tracesReceiver *obsreport.Receiver) {
	d := json.NewDecoder(req.Body)
	for {
		var v interface{}
		err := d.Decode(&v)
		if err != nil {
			if err != io.EOF {
				writeError(w, http.StatusBadRequest, "Error decoding request body")
				return
			}
			break
		}
		fmt.Println(v)
	}
	w.WriteHeader(http.StatusAccepted)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
