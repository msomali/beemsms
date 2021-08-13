package beemsms

import (
	"io"
	"net/http"
)

type ClientOpt func(client *Client)

func WithHTTPClient(client *http.Client) ClientOpt {
	return func(c *Client) {
		c.http = client
	}
}

func WithWriter(writer io.Writer) ClientOpt {
	return func(client *Client) {
		client.logger = writer
	}
}

func WithDebugMode(debug bool) ClientOpt {
	return func(client *Client) {
		client.debug = debug
	}
}
