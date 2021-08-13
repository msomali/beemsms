package internal

import (
	"context"
	"encoding/base64"
	"net/http"
)

var (
	defaultRequestHeaders = map[string]string{
		"Content-Type":  "application/json",
		"Cache-Control": "no-cache",
	}
)

const (
	SendSMS RequestName = iota
	CheckBalance
	CallbackRequest
)

func (rn RequestName) String() string {
	states := [...]string{
		"send SMS",
		"Check Balance",
		"Callback Inner Beem",
	}
	if len(states) < int(rn) {
		return ""
	}

	return states[rn]
}

type (
	//RequestName is used to identify the type of request being saved
	// important in debugging or switch cases where a number of different
	// requests can be served.
	RequestName int

	// Request encapsulate details of a request to be sent to beem.
	Request struct {
		Name        RequestName
		Context     context.Context
		Method      string
		URL         string
		PayloadType PayloadType
		Payload     interface{}
		Headers     map[string]string
		QueryParams map[string]string
	}

	RequestOption func(request *Request)
)

func NewRequest(ctx context.Context, method, url string, payloadType PayloadType, payload interface{}, opts ...RequestOption) *Request {
	request := &Request{
		Context:     ctx,
		Method:      method,
		URL:         url,
		PayloadType: payloadType,
		Payload:     payload,
		Headers:     defaultRequestHeaders,
	}

	for _, opt := range opts {
		opt(request)
	}

	return request
}

func WithRequestContext(ctx context.Context) RequestOption {
	return func(request *Request) {
		request.Context = ctx
	}
}

func WithQueryParams(params map[string]string) RequestOption {
	return func(request *Request) {
		request.QueryParams = params
	}
}

func WithRequestHeaders(headers map[string]string) RequestOption {
	return func(request *Request) {
		request.Headers = headers
	}
}

func WithMoreHeaders(headers map[string]string) RequestOption {
	return func(request *Request) {
		for key, value := range headers {
			request.Headers[key] = value
		}
	}
}

// See 2 (end of page 4) https://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//WithBasicAuth add password and username to request headers
func WithBasicAuth(username, password string) RequestOption {
	return func(request *Request) {
		request.Headers["Authorization"] = "Basic " + basicAuth(username, password)
	}
}

func (request *Request) AddHeader(key, value string) {
	request.Headers[key] = value
}

//NewRequestWithContext takes a *Request and transform into *http.Request with a context
func (request *Request) NewRequestWithContext() (req *http.Request, err error) {

	if request.Payload == nil {
		req, err = http.NewRequestWithContext(request.Context, request.Method, request.URL, nil)
		if err != nil {
			return nil, err
		}
	} else {
		buffer, err := MarshalPayload(request.PayloadType, request.Payload)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequestWithContext(request.Context, request.Method, request.URL, buffer)
		if err != nil {
			return nil, err
		}
	}

	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	for name, value := range request.QueryParams {
		values := req.URL.Query()
		values.Add(name, value)
		req.URL.RawQuery = values.Encode()
	}

	return req, nil
}
