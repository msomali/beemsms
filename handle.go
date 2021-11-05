package beemsms

import (
	"context"
	"github.com/techcraftlabs/base"
	"net/http"
	"time"
)


func (c *Client)handleCallback(writer http.ResponseWriter, request *http.Request){

	ctx, cancel := context.WithTimeout(context.Background(),time.Minute)
	defer cancel()
	var (
		callbackRequest CallbackRequest
		callbackResponse CallbackResponse
		statusCode int
	)
	cTypeHeader := map[string]string{
		"Content-Type": ContentTypeJson,
	}

	_, err := c.rv.Receive(ctx, "delivery report callback",request,&callbackRequest)
	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(writer, err.Error(), statusCode)
		return
	}

	callbackResponse, err = c.Callback.Handle(ctx,callbackRequest)
	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(writer, err.Error(), statusCode)
		return
	}

	res := base.NewResponseBuilder().Headers(cTypeHeader).Payload(callbackResponse).StatusCode(http.StatusOK).Build()

	c.rp.Reply(writer,res)

}