// pickjson is a traefik middleware plugin for picking up specified field of response json
package pickjson

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Field string `json:"field,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Field: "body",
	}
}

// Picker JSON Picker plugin.
type Picker struct {
	next http.Handler
	name string
	cfg  *Config
}

// New created a new Demo plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Picker{
		next: next,
		name: name,
		cfg:  config,
	}, nil
}

func (a *Picker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// make a "response stealer" so we can read response and modify them
	b := make([]byte, 0)
	buf := bytes.NewBuffer(b)
	stealer := newResponseStealer(rw, buf)

	a.next.ServeHTTP(stealer, req)

	resp, err := stealer.Steal()
	if err != nil {
		// do nothing
		// TODO: status code...
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp := map[string]any{}
	if err := json.Unmarshal(resp, &jsonResp); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	field, ok := jsonResp[a.cfg.Field]
	if !ok {
		if err := json.NewEncoder(rw).Encode(jsonResp); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	fieldStr, ok := (field).(string)
	if !ok {
		if err := json.NewEncoder(rw).Encode(jsonResp); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	value := make(map[string]any, 0)
	if err := json.Unmarshal([]byte(fieldStr), &value); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(rw).Encode(value); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
