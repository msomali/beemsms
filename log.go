package beemsms

import (
	"fmt"
	"github.com/techcraftlabs/beemsms/internal"
	"net/http"
	"net/http/httputil"
	"strings"
)

func (c *Client) logPayload(t internal.PayloadType, prefix string, payload interface{}) {
	buf, _ := internal.MarshalPayload(t, payload)
	_, _ = c.logger.Write([]byte(fmt.Sprintf("%s: %s\n\n", prefix, buf.String())))
}

// log is called to print the details of http.Request sent from Tigo during
// callback, namecheck or ussd payment. It is used for debugging purposes
func (c *Client) log(name string, request *http.Request) {

	if request != nil {
		reqDump, _ := httputil.DumpRequest(request, true)
		_, err := fmt.Fprintf(c.logger, "%s REQUEST: %s\n", name, reqDump)
		if err != nil {
			fmt.Printf("error while logging %s request: %v\n",
				strings.ToLower(name), err)
			return
		}
		return
	}
	return
}

// logOut is like log except this is for outgoing client requests:
// http.Request that is supposed to be sent to tigo
func (c *Client) logOut(name string, request *http.Request, response *http.Response) {

	if request != nil {
		reqDump, _ := httputil.DumpRequestOut(request, true)
		_, err := fmt.Fprintf(c.logger, "%s REQUEST: %s\n", name, reqDump)
		if err != nil {
			fmt.Printf("error while logging %s request: %v\n",
				strings.ToLower(name), err)
		}
	}

	if response != nil {
		respDump, _ := httputil.DumpResponse(response, true)
		_, err := fmt.Fprintf(c.logger, "%s RESPONSE: %s\n", name, respDump)
		if err != nil {
			fmt.Printf("error while logging %s response: %v\n",
				strings.ToLower(name), err)
		}
	}

	return
}
