package gwc_test

import (
	"net/http"
	"testing"

	"context"

	"github.com/delicb/cliware"
	"github.com/delicb/gwc"
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
	client := gwc.New(dummyClient())
	req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
	mockMiddleware := &mockMiddleware{}
	req.Use(mockMiddleware)
	_, err := req.Send()
	if err != nil {
		t.Error("Got unexpected error:", err)
	}
	if !mockMiddleware.called {
		t.Error("Middleware not called.")
	}
}

func TestRequest_UseFunc(t *testing.T) {
	client := gwc.New(dummyClient())
	req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
	var called bool
	req.UseFunc(func(next cliware.Handler) cliware.Handler {
		return cliware.HandlerFunc(func(ctx context.Context, req *http.Request) (*http.Response, error) {
			called = true
			return next.Handle(ctx, req)
		})
	})
	_, err := req.Send()
	if err != nil {
		t.Error("Got unexpected error:", err)
	}
	if !called {
		t.Error("Middleware func not called.")
	}
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

func TestRequest_AddPath(t *testing.T) {
	for _, data := range []struct {
		OriginalURL string
		Param       string
		Expected    string
	}{
		{
			OriginalURL: "",
			Param:       "",
			Expected:    "http://",
		},
		{
			OriginalURL: "www.google.com",
			Param:       "",
			Expected:    "http://www.google.com",
		},
		{
			OriginalURL: "www.example.com/path",
			Param:       "",
			Expected:    "http://www.example.com/path",
		},
		{
			OriginalURL: "www.example.com/path",
			Param:       "/additional_path",
			Expected:    "http://www.example.com/path/additional_path",
		},
		{
			OriginalURL: "https://www.example.com/path",
			Param:       "/additional_path",
			Expected:    "https://www.example.com/path/additional_path",
		},
	} {
		client := gwc.New(dummyClient())
		req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())

		resp, err := req.URL(data.OriginalURL).AddPath(data.Param).Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		if resp.Request.URL.String() != data.Expected {
			t.Errorf("Wrong URL. Got: %s, expected: %s.", resp.Request.URL.String(), data.Expected)
		}
	}
}

func TestRequest_Param(t *testing.T) {
	for _, data := range []struct {
		OriginalURL string
		Params      map[string]string
		Expected    string
	}{
		{
			OriginalURL: "",
			Params:      map[string]string{},
			Expected:    "http://",
		},
		{
			OriginalURL: "www.example.com/:param1/keep/:param2",
			Params: map[string]string{
				"param1": "value",
			},
			Expected: "http://www.example.com/value/keep/:param2",
		},
	} {
		client := gwc.New(dummyClient())
		req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
		req.URL(data.OriginalURL)
		for k, v := range data.Params {
			req.Param(k, v)
		}
		resp, err := req.Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		got := resp.Request.URL.String()
		if got != data.Expected {
			t.Errorf("Wrong URL. Got: %s, expected: %s", got, data.Expected)
		}
	}
}

func TestRequest_Params(t *testing.T) {
	for _, data := range []struct {
		OriginalURL string
		Params      map[string]string
		Expected    string
	}{
		{
			OriginalURL: "",
			Params:      map[string]string{},
			Expected:    "http://",
		},
		{
			OriginalURL: "www.example.com/:param1/keep/:param2",
			Params: map[string]string{
				"param1": "value",
			},
			Expected: "http://www.example.com/value/keep/:param2",
		},
	} {
		client := gwc.New(dummyClient())
		req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
		req.URL(data.OriginalURL)
		req.Params(data.Params)

		resp, err := req.Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		got := resp.Request.URL.String()
		if got != data.Expected {
			t.Errorf("Wrong URL. Got: %s, expected: %s", got, data.Expected)
		}
	}
}

func TestRequest_AddQuery(t *testing.T) {
	for _, data := range []struct {
		OriginalURL string
		Params      map[string]string
		Expected    string
	}{
		{
			OriginalURL: "",
			Params:      map[string]string{},
			Expected:    "http://",
		},
		{
			OriginalURL: "",
			Params: map[string]string{
				"a": "b",
			},
			Expected: "http://?a=b",
		},
		{
			OriginalURL: "www.example.com?a=b",
			Params: map[string]string{
				"a": "c",
			},
			Expected: "http://www.example.com?a=b&a=c",
		},
		{
			OriginalURL: "www.example.com/",
			Params: map[string]string{
				"param1": "value",
			},
			Expected: "http://www.example.com/?param1=value",
		},
	} {
		client := gwc.New(dummyClient())
		req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
		req.URL(data.OriginalURL)
		for k, v := range data.Params {
			req.AddQuery(k, v)
		}

		resp, err := req.Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		got := resp.Request.URL.String()
		if got != data.Expected {
			t.Errorf("Wrong URL. Got: %s, expected: %s", got, data.Expected)
		}
	}
}

func TestRequest_SetQuery(t *testing.T) {
	for _, data := range []struct {
		OriginalURL string
		Params      map[string]string
		Expected    string
	}{
		{
			OriginalURL: "",
			Params:      map[string]string{},
			Expected:    "http://",
		},
		{
			OriginalURL: "",
			Params: map[string]string{
				"a": "b",
			},
			Expected: "http://?a=b",
		},
		{
			OriginalURL: "www.example.com?a=b",
			Params: map[string]string{
				"a": "c",
			},
			Expected: "http://www.example.com?a=c",
		},
		{
			OriginalURL: "www.example.com/",
			Params: map[string]string{
				"param1": "value",
			},
			Expected: "http://www.example.com/?param1=value",
		},
	} {
		client := gwc.New(dummyClient())
		req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
		req.URL(data.OriginalURL)
		for k, v := range data.Params {
			req.SetQuery(k, v)
		}

		resp, err := req.Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		got := resp.Request.URL.String()
		if got != data.Expected {
			t.Errorf("Wrong URL. Got: %s, expected: %s", got, data.Expected)
		}
	}
}

func TestRequest_SetQueryParams(t *testing.T) {
	for _, data := range []struct {
		OriginalURL string
		Params      map[string]string
		Expected    string
	}{
		{
			OriginalURL: "",
			Params:      map[string]string{},
			Expected:    "http://",
		},
		{
			OriginalURL: "",
			Params: map[string]string{
				"a": "b",
			},
			Expected: "http://?a=b",
		},
		{
			OriginalURL: "www.example.com?a=b",
			Params: map[string]string{
				"a": "c",
			},
			Expected: "http://www.example.com?a=c",
		},
		{
			OriginalURL: "www.example.com/",
			Params: map[string]string{
				"param1": "value",
			},
			Expected: "http://www.example.com/?param1=value",
		},
	} {
		client := gwc.New(dummyClient())
		req := gwc.NewRequest(client, cliware.NewChain(), cliware.NewChain())
		req.URL(data.OriginalURL)
		req.SetQueryParams(data.Params)

		resp, err := req.Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}
		got := resp.Request.URL.String()
		if got != data.Expected {
			t.Errorf("Wrong URL. Got: %s, expected: %s", got, data.Expected)
		}
	}
}
