/*
Package jusibe provides a Jusibe API client.

Refer to https://jusibe.com/docs/ for more information about Jusibe.

Example Usage:

// Create the Jusibe Configuration. Note that your AccessToken and PublicKey are required
cfg := &jusibe.Config{
	PublicKey: os.Getenv("JUSIBE_PUBLIC_KEY"),
	AccessToken: os.Getenv("JUSIBE_ACCESS_TOKEN"),
}

// Create the client
j, err := jusibe.New(cfg)
if err != nil {
		log.Fatal(err)
}

// Send SMS
to, from, message := "08000000000000", "Azeez", "Hello World"
smsResponse, _, err := j.SendSMS(context.Background(), to, from, message)
if err != nil {
	log.Fatal(err)
}

// Check Delivery Status
deliveryResponse, _, err := j.CheckSMSDeliveryStatus(context.Background(), smsResponse.MessageID)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("%+v\n", deliveryResponse)

// Get SMS credits
creditsResponse, _, err := j.CheckSMSCredits(context.Background())
if err != nil {
	log.Fatal("err")
}
fmt.Printf("%+v\n", creditsResponse)
*/
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

// Config is Jusibe client configuration
// AccessToken and PublicKey are required fields
// APIBaseURL defaults to https://jusibe.com/smsapi/
type Config struct {
	AccessToken string
	PublicKey   string
	APIBaseURL  string
}

// Jusibe is Jusibe API client
type Jusibe struct {
	httpClient  *http.Client
	apiBaseURL  string
	publicKey   string
	accessToken string
}

// createHTTPRequest is a helper method for creating *http.Request used in external API calls
// It returns a *http.Request which has Basic Auth and Context set
func (j *Jusibe) createHTTPRequest(ctx context.Context, method, url string) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, nil)

	if err == nil {
		req.SetBasicAuth(j.publicKey, j.accessToken)
		req = req.WithContext(ctx)
	}

	return
}

// doHTTPRequest performs http requests
// It writes the response body into the body parameter before closing the response body
// It returns the *http.Response for convinience to its caller
func (j *Jusibe) doHTTPRequest(req *http.Request, body interface{}) (res *http.Response, err error) {
	res, err = j.httpClient.Do(req)
	if err != nil {
		return
	}

	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			err = fmt.Errorf("%s, %s", err, closeErr)
		}
	}()

	if res.StatusCode > 299 || res.StatusCode < 200 {
		err = fmt.Errorf("unexpected %d http response code", res.StatusCode)
		return
	}

	err = json.NewDecoder(res.Body).Decode(body)

	return
}

// SendSMS sends SMS to the /send_sms endpoint
// It also returns a *http.Response for convinience to its caller, along with a *SMSResponse and error
func (j *Jusibe) SendSMS(ctx context.Context, to, from, message string) (ssr *SMSResponse, res *http.Response, err error) {
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
	res, err = j.doHTTPRequest(req, ssr)

	return
}

// CheckSMSCredits checks SMS credits using the /get_credits endpoint
// It also returns a *http.Response for convinience to its caller, along with a *SMSCreditsReponse and error
func (j *Jusibe) CheckSMSCredits(ctx context.Context) (scr *SMSCreditsResponse, res *http.Response, err error) {
	url := j.apiBaseURL + "get_credits"
	req, err := j.createHTTPRequest(ctx, http.MethodGet, url)

	if err != nil {
		return
	}

	scr = &SMSCreditsResponse{}
	res, err = j.doHTTPRequest(req, scr)

	return
}

// CheckSMSDeliveryStatus checks a sent SMS (specified by a message id) delivery status using the /delivery_status endpoint
// It also returns a *http.Response for convinience to its caller, along with a *SMSDeliveryResponse and error
func (j *Jusibe) CheckSMSDeliveryStatus(ctx context.Context, messageID string) (sds *SMSDeliveryResponse, res *http.Response, err error) {
	url := fmt.Sprintf("%sdelivery_status?message_id=%s", j.apiBaseURL, messageID)
	req, err := j.createHTTPRequest(ctx, http.MethodGet, url)

	if err != nil {
		return
	}

	sds = &SMSDeliveryResponse{}
	res, err = j.doHTTPRequest(req, sds)

	return
}

// New creates a new Jusibe client configured using the *jusibe.Config parameter
// It uses the default net/http Client with a timeout of 10 seconds
// If you need more control over the http Client, you should use the NewWithHTTPClient function instead
func New(cfg *Config) (j *Jusibe, err error) {
	httpClient := &http.Client{Timeout: defaultHTTPClientTimeout}

	j, err = NewWithHTTPClient(cfg, httpClient)

	return
}

// NewWithHTTPClient creates a new Jusibe client configured using the *jusibe.Config and *http.Client paramerter
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
