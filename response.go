package gwc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"encoding/xml"
)

// Response is thin wrapper around http.Response that provides some
// convenient behavior.
type Response struct {
	*http.Response
	Error error
}

// BuildResponse creates new instance of response based on provided raw HTTP response.
func BuildResponse(rawResponse *http.Response, err error) *Response {
	return &Response{
		Response: rawResponse,
		Error:    err,
	}
}

// SaveToFile writes response content to file with provided path.
func (r *Response) SaveToFile(filename string) error {
	if r.Error != nil {
		return r.Error
	}

	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	defer r.Body.Close()

	_, err = io.Copy(fd, r.Body)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// JSON decodes response body to provided structure from JSON format.
func (r *Response) JSON(userStruct interface{}) error {
	if r.Error != nil {
		return r.Error
	}

	defer r.Body.Close()
	jsonDecoder := json.NewDecoder(r.Body)

	err := jsonDecoder.Decode(&userStruct)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// XML decodes response body to provided structure from XML format.
func (r *Response) XML(userStruct interface{}) error {
	if r.Error != nil {
		return r.Error
	}
	defer r.Body.Close()

	xmlDecoder := xml.NewDecoder(r.Body)
	err := xmlDecoder.Decode(r.Body)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// Bytes returns raw bytes read from response body.
func (r *Response) Bytes() ([]byte, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	defer r.Body.Close()
	buff := bytes.NewBuffer([]byte{})

	// if we got Content-Length set buffer size to it
	if r.ContentLength > 0 {
		buff.Grow(int(r.ContentLength))
	}

	_, err := io.Copy(buff, r.Body)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buff.Bytes(), nil
}

// String returns response body in string format.
func (r *Response) String() (string, error) {
	bytes, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
