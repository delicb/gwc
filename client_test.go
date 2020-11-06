package gwc_test

import (
	"context"
	"net/http"
	"testing"

	"reflect"

	"github.com/delicb/cliware"
	"github.com/delicb/gwc"
)

type mockTransport struct {
	called       bool
	responseCode int
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.called = true
	return &http.Response{
		StatusCode: 200,
		Request:    req,
	}, nil
}

func newMockTransport(responseCode int) *mockTransport {
	return &mockTransport{responseCode: responseCode}
}

type mockMiddleware struct {
	called bool
	count  int
}

func (m *mockMiddleware) Exec(next cliware.Handler) cliware.Handler {
	return cliware.HandlerFunc(func(req *http.Request) (*http.Response, error) {
		m.called = true
		m.count += 1
		return next.Handle(req)
	})
}

func dummyClient() *http.Client {
	return &http.Client{
		Transport: newMockTransport(200),
	}
}

func TestNew(t *testing.T) {
	mockTransport := newMockTransport(200)
	httpClient := &http.Client{
		Transport: mockTransport,
	}
	mockMiddleware := &mockMiddleware{}
	client := gwc.New(httpClient, mockMiddleware)
	client.Get().Send()
	if !mockTransport.called {
		t.Error("Request not send - RoundTripper not called.")
	}
	if !mockMiddleware.called {
		t.Error("Middleware not called.")
	}
}

func TestUse(t *testing.T) {
	client := gwc.New(dummyClient())
	mockMiddleware := &mockMiddleware{}
	client.Use(mockMiddleware)
	client.Get().Send()
	if !mockMiddleware.called {
		t.Error("Middleware set to client not called.")
	}
}

func TestUseFunc(t *testing.T) {
	client := gwc.New(dummyClient())
	var called bool
	client.UseFunc(func(next cliware.Handler) cliware.Handler {
		return cliware.HandlerFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			return next.Handle(req)
		})
	})
	client.Get().Send()
	if !called {
		t.Error("Middleware that was added as function not called.")
	}
}

func TestUsePost(t *testing.T) {
	client := gwc.New(dummyClient())
	mockMiddleware := &mockMiddleware{}
	client.UsePost(mockMiddleware)
	client.Get().Send()
	if !mockMiddleware.called {
		t.Error("Middleware added by UsePost not called.")
	}
}

func TestUsePostFunc(t *testing.T) {
	client := gwc.New(dummyClient())
	var called bool
	client.UsePostFunc(func(next cliware.Handler) cliware.Handler {
		return cliware.HandlerFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			return next.Handle(req)
		})
	})
	client.Get().Send()
	if !called {
		t.Error("Middleware added by UsePostFunc not called.")
	}
}

func TestMiddlewareOrder(t *testing.T) {
	client := gwc.New(dummyClient())
	order := []string{}
	client.UseFunc(func(next cliware.Handler) cliware.Handler {
		return cliware.HandlerFunc(func(req *http.Request) (*http.Response, error) {
			order = append(order, "pre")
			return next.Handle(req)
		})
	})
	client.UsePostFunc(func(next cliware.Handler) cliware.Handler {
		return cliware.HandlerFunc(func(req *http.Request) (*http.Response, error) {
			order = append(order, "post")
			return next.Handle(req)
		})
	})
	client.Get().Send()
	if !reflect.DeepEqual(order, []string{"pre", "post"}) {
		t.Errorf("Wrong order of calling middlewares. Got: %v", order)
	}
}

func TestRequest(t *testing.T) {
	client := gwc.New(dummyClient())
	req := client.Request()
	// just check if we got non-nil request, most of the stuff are either private
	// or tested with other tests
	if req == nil {
		t.Error("Got nil request.")
	}
}

func TestClientMethods(t *testing.T) {
	for _, data := range []struct {
		Function func(*gwc.Client) *gwc.Request
		Method   string
	}{
		{
			(*gwc.Client).Get,
			"GET",
		},
		{
			(*gwc.Client).Post,
			"POST",
		},
		{
			(*gwc.Client).Put,
			"PUT",
		},
		{
			(*gwc.Client).Delete,
			"DELETE",
		},
		{
			(*gwc.Client).Patch,
			"PATCH",
		},
		{
			(*gwc.Client).Head,
			"HEAD",
		},
		{
			(*gwc.Client).Options,
			"OPTIONS",
		},
	} {
		c := gwc.New(dummyClient())
		req := data.Function(c)
		resp, err := req.Send()
		if err != nil {
			t.Error("Got unexpected error: ", err)
		}

		if resp.Request.Method != data.Method {
			t.Errorf("Wrong request method. Got: %s, expected: %s", resp.Request.Method, data.Method)
		}
	}
}

func TestDo(t *testing.T) {
	client := gwc.New(dummyClient())
	mockMiddleware := &mockMiddleware{}
	_, err := client.Do(mockMiddleware)
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}
	if !mockMiddleware.called {
		t.Error("Middleware passed to DoCtx not called.")
	}
}

func TestDoCtx(t *testing.T) {
	client := gwc.New(dummyClient())
	mockMiddleware := &mockMiddleware{}
	ctx := context.WithValue(context.Background(), "key", "value")
	resp, err := client.DoCtx(ctx, mockMiddleware)
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}
	contextValue := resp.Request.Context().Value("key")
	if contextValue != "value" {
		t.Errorf("Wrong context value. Got: %s, expected: value", contextValue)
	}
	if !mockMiddleware.called {
		t.Error("Middleware passed to DoCtx not called.")
	}
}

func TestClientOnContext(t *testing.T) {
	httpClient := dummyClient()
	client := gwc.New(httpClient)
	resp, err := client.Get().Send()
	if err != nil {
		t.Error("Got unexpected error: ", err)
	}
	cl := gwc.ClientFromContext(resp.Request.Context())
	if cl != httpClient {
		t.Error("Got wrong instance of http.Client.")
	}
}

func TestClient_AfterCalledOnce(t *testing.T) {
	client := gwc.New(dummyClient())
	countingMockMiddleware := &mockMiddleware{}
	client.UsePost(countingMockMiddleware)
	workingMockMiddleware := &mockMiddleware{}
	_, err := client.Do(workingMockMiddleware)
	if err != nil {
		t.Error("Got unexpected error:", err)
	}
	if !workingMockMiddleware.called {
		t.Error("Working middleware not called.")
	}
	if countingMockMiddleware.count != 1 {
		t.Error("Middleware not called only once.")
	}
}
