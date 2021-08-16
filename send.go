package beemsms

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/techcraftlabs/beemsms/internal"
	"io"
	"net/http"
	"strings"
)

func (c *Client) send(ctx context.Context, rn internal.RequestName, request *internal.Request, v interface{}) error {

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	var (
		req *http.Request
	)
	var res *http.Response

	var reqBodyBytes, resBodyBytes []byte
	defer func(debug bool) {
		if debug {
			req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
			res.Body = io.NopCloser(bytes.NewBuffer(resBodyBytes))
			name := strings.ToUpper(strings.ToUpper(rn.String()))
			c.logOut(name, req, res)
		}
	}(c.debug)
	req, err := request.NewRequestWithContext()

	if err != nil {
		return err
	}

	if req.Body != nil {
		reqBodyBytes, _ = io.ReadAll(req.Body)
	}

	if v == nil {
		return errors.New("v interface can not be empty")
	}

	// restore request body for logging
	req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))

	res, err = c.http.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		resBodyBytes, _ = io.ReadAll(res.Body)
	}

	contentType := res.Header.Get("Content-Type")
	if strings.Contains(contentType, ContentTypeJson) {
		if err := json.NewDecoder(bytes.NewBuffer(resBodyBytes)).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}

	if strings.Contains(contentType, ContentTypeXml) ||
		strings.Contains(contentType, ContentTypeTextXML) {
		if err := xml.NewDecoder(bytes.NewBuffer(resBodyBytes)).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}
	return nil
}
