package gwc_test

import (
	"testing"

	"context"

	"go.delic.rs/cliware"
	"go.delic.rs/gwc"
)

func TestNewRequest(t *testing.T) {
	client := gwc.New(dummyClient())
	req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
	if req.Client != client {
		t.Error("Wrong instance of client for request.")
	}
}

func TestRequest_Context(t *testing.T) {
	client := gwc.New(dummyClient())
	req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
	if req.Context() != nil {
		t.Error("Got non-nil context for fresh request.")
	}
	ctx := context.Background()
	req.SetContext(ctx)
	if req.Context() != ctx {
		t.Error("Got wrong context from request.")
	}
}

func TestRequest_Use(t *testing.T) {

}

func TestRequest_UseFunc(t *testing.T) {

}

func TestRequest_Method(t *testing.T) {
	for _, method := range []string{
		"GET", "POST", "PUT", "DELETE", "OPTIONS", "CUSTOM",
	} {
		client := gwc.New(dummyClient())
		req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
		resp, err := req.Method(method).Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		if resp.Request.Method != method {
			t.Errorf("Wrong request method. Got: %s, expected: %s", resp.Request.Method, method)
		}
	}
}

func TestRequest_URL(t *testing.T) {
	for toSet, expected := range map[string]string{
		"http://example.com":           "http://example.com",
		"https://example.com/path":     "https://example.com/path",
		"https://example.com/path?q=v": "https://example.com/path?q=v",
	} {
		client := gwc.New(dummyClient())
		req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
		resp, err := req.URL(toSet).Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		got := resp.Request.URL.String()
		if got != expected {
			t.Errorf("Wrong request URL. Got: %s, expected: %s", got, expected)
		}
	}
}

func TestRequest_BaseURL(t *testing.T) {
	for toSet, expected := range map[string]string{
		"http://example.com":           "http://example.com",
		"https://example.com/path":     "https://example.com",
		"https://example.com/path?q=v": "https://example.com",
	} {
		client := gwc.New(dummyClient())
		req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
		resp, err := req.BaseURL(toSet).Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		got := resp.Request.URL.String()
		if got != expected {
			t.Errorf("Wrong request URL. Got: %s, expected: %s", got, expected)
		}
	}
}

func TestRequest_Path(t *testing.T) {
	for toSet, expected := range map[string]string{
		"":      "",
		"/path": "/path",
	} {
		client := gwc.New(dummyClient())
		req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
		resp, err := req.Path(toSet).Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		got := resp.Request.URL.String()
		if got != expected {
			t.Errorf("Wrong request URL. Got: %s, expected: %s", got, expected)
		}
	}
}
