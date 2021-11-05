package beemsms

import (
	"io"
)

type ClientOpt func(client *Client)

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
