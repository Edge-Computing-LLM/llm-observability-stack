package gateway

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func testServer(t *testing.T) (*Server, *httptest.Server) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/generate" {
			_ = json.NewEncoder(w).Encode(map[string]any{"response": "synthetic-response", "done": true})
			return
		}
		w.Header().Set("Content-Type", "application/x-ndjson")
		_, _ = w.Write([]byte("{\"token\":\"hello\"}\n{\"token\":\"world\"}\n"))
	}))
	t.Cleanup(upstream.Close)
	server, err := New(Config{Address: ":0", UpstreamURL: upstream.URL, Model: defaultModel, ServiceName: "test", Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	return server, httptest.NewServer(server.Handler())
}
func TestEndpointsAndInvoke(t *testing.T) {
	_, httpServer := testServer(t)
	defer httpServer.Close()
	for _, path := range []string{"/", "/healthz", "/config", "/metrics"} {
		response, err := http.Get(httpServer.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != 200 {
			t.Fatalf("%s returned %d", path, response.StatusCode)
		}
		_ = response.Body.Close()
	}
	response, err := http.Post(httpServer.URL+"/invoke", "application/json", strings.NewReader(`{"prompt":"hello","system":"brief"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var result map[string]any
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result["response"] != "synthetic-response" {
		t.Fatalf("response: %#v", result)
	}
}
func TestStreamingProxy(t *testing.T) {
	_, httpServer := testServer(t)
	defer httpServer.Close()
	response, err := http.Post(httpServer.URL+"/ollama/api/chat", "application/json", strings.NewReader(`{"stream":true}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "{\"token\":\"hello\"}\n{\"token\":\"world\"}\n" {
		t.Fatalf("content: %q", content)
	}
}
