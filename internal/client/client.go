package client

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zx06/xurl/internal/errors"
)

type DoInput struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

type Result struct {
	Status     int               `json:"status"`
	StatusText string            `json:"status_text,omitempty"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

func Do(client HTTPDoer, in DoInput) (Result, *errors.XError) {
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}
	var bodyReader io.Reader
	if in.Body != "" {
		bodyReader = strings.NewReader(in.Body)
	}
	req, err := http.NewRequest(in.Method, in.URL, bodyReader)
	if err != nil {
		return Result{}, errors.Wrap(errors.CodeRequest, "failed to build request", nil, err)
	}
	for k, v := range in.Headers {
		req.Header.Set(k, v)
	}
	if in.Body != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "text/plain")
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{}, errors.Wrap(errors.CodeRequest, "request failed", map[string]any{"url": in.URL}, err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, errors.Wrap(errors.CodeRequest, "failed to read response body", nil, err)
	}

	headers := map[string]string{}
	for k, vals := range resp.Header {
		if len(vals) > 0 {
			headers[k] = strings.Join(vals, ", ")
		}
	}

	return Result{
		Status:     resp.StatusCode,
		StatusText: resp.Status,
		Headers:    headers,
		Body:       string(b),
	}, nil
}
