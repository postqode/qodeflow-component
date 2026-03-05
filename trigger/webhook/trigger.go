package webhook

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/postqode/qodeflow-core/data/metadata"
	"github.com/postqode/qodeflow-core/support/log"
	"github.com/postqode/qodeflow-core/trigger"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}

func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{id: config.Id, settings: s}, nil
}

type Trigger struct {
	server   *http.Server
	settings *Settings
	id       string
	logger   log.Logger
}

func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	t.logger = ctx.Logger()
	router := mux.NewRouter()

	for _, handler := range ctx.GetHandlers() {
		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}

		method := strings.ToUpper(s.Method)
		path := s.Path

		t.logger.Debugf("Registering webhook handler [%s: %s]", method, path)
		router.HandleFunc(path, t.newHTTPHandler(handler)).Methods(method)
	}

	addr := ":" + strconv.Itoa(t.settings.Port)
	t.server = &http.Server{Addr: addr, Handler: router}

	return nil
}

func (t *Trigger) Start() error {
	t.logger.Infof("Webhook trigger starting on %s", t.server.Addr)
	go func() {
		if err := t.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.logger.Errorf("Webhook trigger failed: %v", err)
		}
	}()
	return nil
}

func (t *Trigger) Stop() error {
	return t.server.Shutdown(context.Background())
}

func (t *Trigger) newHTTPHandler(handler trigger.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out := &Output{
			Method:      r.Method,
			PathParams:  make(map[string]string),
			QueryParams: make(map[string]string),
			Headers:     make(map[string]string),
		}

		vars := mux.Vars(r)
		for key, value := range vars {
			out.PathParams[key] = value
		}

		for key, values := range r.URL.Query() {
			out.QueryParams[key] = strings.Join(values, ",")
		}

		for key, values := range r.Header {
			out.Headers[key] = strings.Join(values, ",")
		}

		if r.ContentLength > 0 {
			var content any
			err := json.NewDecoder(r.Body).Decode(&content)
			if err != nil && err != io.EOF {
				t.logger.Errorf("Error decoding webhook body: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			out.Content = content
		}

		results, err := handler.Handle(context.Background(), out)
		if err != nil {
			t.logger.Errorf("Error handling webhook: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reply := &Reply{}
		if err := reply.FromMap(results); err != nil {
			t.logger.Errorf("Error mapping webhook results: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if reply.Code == 0 {
			reply.Code = http.StatusOK
		}

		if reply.Data != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(reply.Code)
			json.NewEncoder(w).Encode(reply.Data)
		} else {
			w.WriteHeader(reply.Code)
		}
	}
}
