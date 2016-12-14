package gwc

import (
	"fmt"
	"net/http"

	"go.delic.rs/cliware-middlewares/headers"
	"go.delic.rs/cliware-middlewares/errors"
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
		errors.Errors(),
	)
	// send request
	resp, err := client.Get().URL("https://httpbin.org/get").Send()

	// check errors
	// because of errors middleware included in client, ever status codes
	// 400+ will be turned into errors.
	if err != nil {
		panic(err)
	}

	resp.JSON(respBody)
	fmt.Println(respBody.Headers.UserAgent)
	// output: example-client
}
