package beemsms

import (
	"context"
	"github.com/techcraftt/beemsms/internal"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	defaultTimeout     = 60 * time.Second
	ContentTypeTextXML = "text/xml"
	ContentTypeXml     = "application/xml"
	ContentTypeJson    = "application/json"
)

var (
	_             CallbackHandler = (*CallbackFunc)(nil)
	_             Service         = (*Client)(nil)
	defaultWriter io.Writer       = os.Stderr
	defaultClient                 = http.DefaultClient
)

type (
	CallbackRequest struct {
		RequestID   string `json:"request_id,omitempty"`
		RecipientID string `json:"recipient_id,omitempty"`
		DestAddr    string `json:"dest_addr,omitempty"`
		Status      string `json:"Status,omitempty"`
	}

	CallbackResponse struct {
		RequestID string `json:"request_id,omitempty"`
		Status    string `json:"Status,omitempty"`
		Success   string `json:"successful,omitempty"`
	}

	ErrResponse struct {
		Code    int64  `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	}

	BalanceResponse struct {
		Data struct {
			CreditBalance int `json:"credit_balance,omitempty"`
		} `json:"data,omitempty"`
		Code    int64  `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	}

	//SendRequest contains details of a send sms request body:
	SendRequest struct {
		//Source (source_addr) is a message source address or sender ID. Limited to 11 characters if text.
		//Or valid mobile number in valid international number format with country code.
		//No leading + sign. Has to be active on the sms portal.
		Source string `json:"source_addr"`
		//ScheduleTime this is optional field that specify the Scheduled time
		//of the message. GMT+0 timezone. Format (yyyy-mm-dd hh:mm)
		ScheduleTime string `json:"schedule_time,omitempty"`
		//Message contains the payload/ content .
		Message string `json:"message"`
		//Encoding (number) is a message encoding type.Default value is 0
		Encoding string `json:"encoding"`
		//Recipients is Array of destination numbers.
		Recipients []Recipient `json:"recipients"`
	}

	SendResponse struct {
		Successful bool   `json:"successful,omitempty"`
		RequestID  int64  `json:"request_id,omitempty"`
		Code       int64  `json:"code,omitempty"`
		Message    string `json:"message,omitempty"`
		Valid      int64  `json:"valid,omitempty"`
		Invalid    int64  `json:"invalid,omitempty"`
		Duplicates int64  `json:"duplicates,omitempty"`
	}

	Recipient struct {
		ID    int64  `json:"recipient_id,omitempty"`
		Phone string `json:"dest_addr,omitempty"`
	}

	Config struct {
		SendSMSURL      string
		CheckBalanceURL string
		CallbackURL     string
		APIKey          string
		SecretKey       string
	}

	CallbackHandler interface {
		Handle(ctx context.Context, request CallbackRequest) (CallbackResponse, error)
	}

	CallbackFunc func(ctx context.Context, request CallbackRequest) (CallbackResponse, error)

	Client struct {
		*Config
		logger   io.Writer
		http     *http.Client
		Callback CallbackHandler
		debug    bool
	}

	Service interface {
		Balance(ctx context.Context) (response BalanceResponse, err error)
		Text(ctx context.Context, req SendRequest) (response SendResponse, err error)
		CallbackHandler(ctx context.Context) http.HandlerFunc
	}
)

func (c *Client) Balance(ctx context.Context) (response BalanceResponse, err error) {
	var opts []internal.RequestOption
	authOpt := internal.WithBasicAuth(c.Config.APIKey, c.Config.SecretKey)
	cTypeHeader := map[string]string{
		"Content-Type": ContentTypeJson + "; charset=utf-8",
	}
	cTypeOpt := internal.WithMoreHeaders(cTypeHeader)
	opts = append(opts, authOpt, cTypeOpt)
	request := internal.NewRequest(ctx, http.MethodGet, c.Config.CheckBalanceURL, internal.JsonPayload, nil, opts...)

	err = c.send(ctx, internal.CheckBalance, request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (c *Client) Text(ctx context.Context, req SendRequest) (response SendResponse, err error) {
	var opts []internal.RequestOption
	authOpt := internal.WithBasicAuth(c.Config.APIKey, c.Config.SecretKey)
	cTypeHeader := map[string]string{
		"Content-Type": ContentTypeJson,
	}
	cTypeOpt := internal.WithMoreHeaders(cTypeHeader)
	opts = append(opts, authOpt, cTypeOpt)
	request := internal.NewRequest(ctx, http.MethodPost, c.Config.SendSMSURL, internal.JsonPayload, req, opts...)
	err = c.send(ctx, internal.SendSMS, request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (c *Client) CallbackHandler(ctx context.Context) http.HandlerFunc {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return c.handle(ctx, internal.JsonPayload, internal.CallbackRequest)

}

func NewClient(config *Config, handler CallbackHandler, opts ...ClientOpt) *Client {
	c := &Client{
		Config:   config,
		logger:   defaultWriter,
		http:     defaultClient,
		Callback: handler,
		debug:    true,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c CallbackFunc) Handle(ctx context.Context, request CallbackRequest) (CallbackResponse, error) {
	return c(ctx, request)
}
