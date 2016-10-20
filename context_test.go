package gwc_test

import (
	"context"
	"net/http"
	"testing"

	m "github.com/delicb/gwc2/middleware"
)

func TestClientContextSuccess(t *testing.T) {
	ctx := context.Background()
	cl := new(http.Client)
	clientContext := m.SetClient(ctx, cl)

	returnedClient := m.GetClient(clientContext)
	if cl != returnedClient {
		t.Fatal("Client set to context did not match.")
	}
}

func TestClientContextNoClient(t *testing.T) {
	ctx := context.Background()
	cl := m.GetClient(ctx)
	if cl != nil {
		t.Fatal("Expected nil client, got: ", cl)
	}
}
