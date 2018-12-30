package jusibe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultAPIBaseURL        = "https://jusibe.com/smsapi/"
	defaultHTTPClientTimeout = (time.Second * 10)
)

// Config ...
type Config struct {
	AccessToken string
	PublicKey   string
	APIBaseURL  string
}

// Jusibe ...
type Jusibe struct {
	httpClient  *http.Client
	apiBaseURL  string
	publicKey   string
	accessToken string
}

// createHTTPRequest
func (j *Jusibe) createHTTPRequest(ctx context.Context, method, url string) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, nil)

	if err == nil {
		req.SetBasicAuth(j.publicKey, j.accessToken)
	}

	return
}

// doHTTPRequest ...
func (j *Jusibe) doHTTPRequest(req *http.Request, response interface{}) (statusCode int, err error) {
	res, err := j.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	statusCode = res.StatusCode

	if res.StatusCode > 299 || res.StatusCode < 200 {
		err = fmt.Errorf("unexpected %d http response code", res.StatusCode)
		return
	}

	err = json.NewDecoder(res.Body).Decode(response)

	return
}

// SendSMS ...
func (j *Jusibe) SendSMS(ctx context.Context, to, from, message string) (ssr *SMSResponse, statusCode int, err error) {
	// This check is defined in Jusibe API docs
	if len(from) > 11 {
		err = errors.New("from (SenderID) allows maximum of eleven (11) characters. See API docs https://jusibe.com/docs/")
		return
	}

	url := fmt.Sprintf("%ssend_sms?to=%s&from=%s&message=%s", j.apiBaseURL, to, from, message)

	req, err := j.createHTTPRequest(ctx, http.MethodPost, url)
	if err != nil {
		return
	}

	ssr = &SMSResponse{}
	statusCode, err = j.doHTTPRequest(req, ssr)

	return
}

// CheckSMSCredits ...
func (j *Jusibe) CheckSMSCredits(ctx context.Context) (scr *SMSCreditsResponse, statusCode int, err error) {
	url := j.apiBaseURL + "get_credits"
	req, err := j.createHTTPRequest(ctx, http.MethodGet, url)

	if err != nil {
		return
	}

	scr = &SMSCreditsResponse{}
	statusCode, err = j.doHTTPRequest(req, scr)

	return
}

// CheckSMSDeliveryStatus ...
func (j *Jusibe) CheckSMSDeliveryStatus(ctx context.Context, messageID string) (sds *SMSDeliveryResponse, statusCode int, err error) {
	url := fmt.Sprintf("%sdelivery_status?message_id=%s", j.apiBaseURL, messageID)
	req, err := j.createHTTPRequest(ctx, http.MethodGet, url)

	if err != nil {
		return
	}

	sds = &SMSDeliveryResponse{}
	statusCode, err = j.doHTTPRequest(req, sds)

	return
}

// New returns new Jusibe client configured using the Configer
func New(cfg *Config) (j *Jusibe, err error) {
	httpClient := &http.Client{Timeout: defaultHTTPClientTimeout}

	j, err = NewWithHTTPClient(cfg, httpClient)

	return
}

// NewWithHTTPClient ...
func NewWithHTTPClient(cfg *Config, httpClient *http.Client) (j *Jusibe, err error) {
	if cfg.AccessToken == "" || cfg.PublicKey == "" {
		err = errors.New("Failed to create New Jusibe client. accessToken and publicKey are required")
		return
	}

	if cfg.APIBaseURL == "" {
		cfg.APIBaseURL = defaultAPIBaseURL
	}

	j = &Jusibe{
		httpClient:  httpClient,
		accessToken: cfg.AccessToken,
		publicKey:   cfg.PublicKey,
		apiBaseURL:  cfg.APIBaseURL,
	}

	return
}
