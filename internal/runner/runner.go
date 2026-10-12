package runner

import (
	"github.com/zx06/xurl/internal/client"
	"github.com/zx06/xurl/internal/config"
	"github.com/zx06/xurl/internal/errors"
	"github.com/zx06/xurl/internal/httpfile"
	"github.com/zx06/xurl/internal/profile"
	"github.com/zx06/xurl/internal/secret"
)

type RequestInfo struct {
	Name   string            `json:"name,omitempty"`
	Method string            `json:"method"`
	URL    string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}

type RunResult struct {
	OK       bool              `json:"ok"`
	Request  RequestInfo       `json:"request"`
	Status   int               `json:"status,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
	Body     string            `json:"body,omitempty"`
	Error    *ErrorObj         `json:"error,omitempty"`
}

type ErrorObj struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type Options struct {
	Profile     config.Profile
	Attrs       map[string]string
	HTTPClient  client.HTTPDoer
	Keyring     secret.KeyringAPI
}

func Run(req httpfile.Request, opts Options) RunResult {
	resolved, xe := profile.ResolveSecrets(opts.Profile, opts.Keyring)
	if xe != nil {
		return failRun(req, xe)
	}

	vars := httpfile.BuildVars(resolved.BaseURL, opts.Profile.Vars, req.Vars, opts.Attrs)
	urlStr := httpfile.Interpolate(req.URL, vars)
	urlStr, xe = profile.ResolveURL(resolved.BaseURL, urlStr)
	if xe != nil {
		return failRun(req, xe)
	}

	headers := map[string]string{}
	for k, v := range resolved.Headers {
		headers[k] = httpfile.Interpolate(v, vars)
	}
	for k, v := range req.Headers {
		headers[k] = httpfile.Interpolate(v, vars)
	}

	auth := resolved.Auth
	if auth.Bearer != "" {
		auth.Bearer = httpfile.Interpolate(auth.Bearer, vars)
	}
	if auth.Basic != nil {
		auth.Basic.Password = httpfile.Interpolate(auth.Basic.Password, vars)
	}
	if auth.Header != nil {
		auth.Header.Value = httpfile.Interpolate(auth.Header.Value, vars)
	}
	profile.ApplyAuth(headers, auth)

	body := httpfile.Interpolate(req.Body, vars)
	method := req.Method

	info := RequestInfo{
		Name:   req.Name,
		Method: method,
		URL:    urlStr,
		Headers: redactHeaders(headers),
	}

	httpRes, xe := client.Do(opts.HTTPClient, client.DoInput{
		Method:  method,
		URL:     urlStr,
		Headers: headers,
		Body:    body,
	})
	if xe != nil {
		return RunResult{
			OK:      false,
			Request: info,
			Error:   toErrorObj(xe),
		}
	}

	return RunResult{
		OK:      httpRes.Status >= 200 && httpRes.Status < 300,
		Request: info,
		Status:  httpRes.Status,
		Headers: httpRes.Headers,
		Body:    httpRes.Body,
	}
}

func failRun(req httpfile.Request, xe *errors.XError) RunResult {
	return RunResult{
		OK: false,
		Request: RequestInfo{
			Name:   req.Name,
			Method: req.Method,
			URL:    req.URL,
		},
		Error: toErrorObj(xe),
	}
}

func toErrorObj(xe *errors.XError) *ErrorObj {
	if xe == nil {
		return nil
	}
	return &ErrorObj{Code: xe.Code, Message: xe.Message, Details: xe.Details}
}

func redactHeaders(h map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range h {
		if k == "Authorization" {
			out[k] = "***"
		} else {
			out[k] = v
		}
	}
	return out
}
