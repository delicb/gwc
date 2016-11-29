package gwc

import (
	"fmt"
	"net/http"

	"go.delic.rs/cliware-middlewares/headers"
	"go.delic.rs/cliware-middlewares/responsebody"
)

type HTTPBinResponse struct {
	Headers struct {
		UserAgent string `json:"User-Agent"`
	} `json:"headers"`
}

func Example() {
	respBody := new(HTTPBinResponse)

	client := New(
		http.DefaultClient,
		// Add user agent header, it will be applied to all requests
		headers.Add("User-Agent", "example-client"),
		responsebody.JSON(respBody),
	)
	// send request
	resp, err := client.Get().URL("https://httpbin.org/get").Send()
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("Expected status code 200, got: %s", resp.Status))
	}

	fmt.Println(respBody.Headers.UserAgent)
	// output: example-client
}
