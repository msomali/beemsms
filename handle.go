package beemsms

import (
	"bytes"
	"context"
	"github.com/techcraftlabs/beemsms/internal"
	"io/ioutil"
	"net/http"
)

func (c *Client) handle(ctx context.Context, payloadType internal.PayloadType, rn internal.RequestName) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		switch rn {
		case internal.CallbackRequest:

			if c.debug {
				request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				c.log(rn.String(), request)
			}

			request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			var callbackRequest CallbackRequest
			var callbackResponse CallbackResponse
			var response *internal.Response
			var statusCode int

			statusCode = 200

			defer func(debug bool) {
				if debug {
					c.logPayload(payloadType, "callback response", &callbackResponse)
					return
				}
				return
			}(c.debug)

			err := internal.Receive(request, payloadType, &callbackRequest)

			if err != nil {
				statusCode = http.StatusInternalServerError
				http.Error(writer, err.Error(), statusCode)
				return
			}

			callbackResponse, err = c.Callback.Handle(ctx, callbackRequest)

			if err != nil {
				statusCode = http.StatusInternalServerError
				http.Error(writer, err.Error(), statusCode)
				return
			}

			var responseOpts []internal.ResponseOption
			headers := internal.WithDefaultJsonHeader()

			responseOpts = append(responseOpts, headers, internal.WithErr(err))
			response = internal.NewResponse(statusCode, callbackResponse, payloadType, responseOpts...)

			internal.Reply(response, writer)

			return

		default:
			msg := "unknown request type could not handle"
			statusCode := http.StatusInternalServerError
			http.Error(writer, msg, statusCode)
			return
		}

	}
}
