package gwc_test

import (
	"context"
	"net/http"
	"testing"

	"go.delic.rs/gwc"
)

func TestClientContextSuccess(t *testing.T) {
	ctx := context.Background()
	cl := new(http.Client)
	clientContext := gwc.SetClient(ctx, cl)

	returnedClient := gwc.GetClient(clientContext)
	if cl != returnedClient {
		t.Fatal("Client set to context did not match.")
	}
}

func TestClientContextNoClient(t *testing.T) {
	ctx := context.Background()
	cl := gwc.GetClient(ctx)
	if cl != nil {
		t.Fatal("Expected nil client, got: ", cl)
	}
}
