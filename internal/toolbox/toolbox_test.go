package toolbox

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDNS(t *testing.T) {
	var output bytes.Buffer
	if err := dns([]string{"localhost"}, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "localhost") {
		t.Fatalf("output: %s", output.String())
	}
}
func TestHTTPCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) }))
	defer server.Close()
	var output bytes.Buffer
	if err := httpCheck(context.Background(), []string{server.URL}, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "status=200") {
		t.Fatalf("output: %s", output.String())
	}
}
