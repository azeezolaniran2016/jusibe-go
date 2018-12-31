package jusibe

import (
	"bytes"
	"context"
	"fmt"
	"github.com/azeezolaniran2016/jusibe-go/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestJusibe(t *testing.T) {
	t.Run("New should fail when required config fields are empty", func(t *testing.T) {
		_, err := New(&Config{PublicKey: "some_public_key"})
		assert.Error(t, err, "should return error when Config.AccessToken is empty")

		_, err = New(&Config{AccessToken: "some_access_token"})
		assert.Error(t, err, "should return error when Config.PublicKey is empty")

		_, err = New(&Config{})
		assert.Error(t, err, "should return error when Config.AccessToken and Config.PublicKey are empty")
	})

	t.Run("New should set non-required Config fields to their default value", func(t *testing.T) {
		accessToken, publicKey := "some_access_token", "some_public_key"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey}
		jusibe, err := New(cfg)

		assert.NoError(t, err, "should not return error when creating New instance")

		assert.NotNil(t, jusibe.httpClient, "should set default http Client when none is specified")

		assert.Equal(t, defaultHTTPClientTimeout, jusibe.httpClient.Timeout, "should have default http Client timeout")

		assert.Equal(t, defaultAPIBaseURL, jusibe.apiBaseURL, "should default APIBaseURL")
	})

	t.Run("New should set fields on Jusibe instance", func(t *testing.T) {
		accessToken, publicKey, baseURL := "some_access_token", "some_public_key", "https://jusibe.com/"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey, APIBaseURL: baseURL}
		jusibe, err := New(cfg)

		assert.NoError(t, err, "Should not return error when creating Jusibe instance with New function")

		assert.Equal(t, baseURL, jusibe.apiBaseURL, "Should set apiBaseURL")

		assert.Equal(t, publicKey, jusibe.publicKey, "Should set publicKey")

		assert.Equal(t, accessToken, jusibe.accessToken, "Should set accessToken")

		assert.NotNil(t, jusibe.httpClient, "Should set httpClient to default")

		assert.Equal(t, defaultHTTPClientTimeout, jusibe.httpClient.Timeout, "Should set httpClient timeout to default")
	})

	t.Run("NewWithHTTPClient", func(t *testing.T) {
		accessToken, publicKey, baseURL := "some_access_token", "some_public_key", "https://jusibe.com/"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey, APIBaseURL: baseURL}

		testTimeout := (10 * time.Second)
		httpClient := &http.Client{Timeout: testTimeout}

		jusibe, err := NewWithHTTPClient(cfg, httpClient)

		assert.NoError(t, err, "Should not return error when creating Jusibe instance with NewWithHTTPClient function")

		assert.NotNil(t, jusibe.httpClient, "Should set httpClient")

		assert.Equal(t, testTimeout, jusibe.httpClient.Timeout, "Should have specified httpClient timeout")

		assert.Equal(t, baseURL, jusibe.apiBaseURL, "Should set apiBaseURL")

		assert.Equal(t, publicKey, jusibe.publicKey, "Should set publicKey field")

		assert.Equal(t, accessToken, jusibe.accessToken, "Should set accessToken")
	})

	t.Run("SendSMS", func(t *testing.T) {
		accessToken, publicKey, baseURL := "some_access_token", "some_public_key", "https://jusibe.com/"
		to, from, message := "09001000101", "test_user", "Hello World!"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey, APIBaseURL: baseURL}

		mockController := gomock.NewController(t)
		mockRoundTripper := mocks.NewMockRoundTripper(mockController)

		httpClient := &http.Client{Transport: mockRoundTripper}

		jusibe, err := NewWithHTTPClient(cfg, httpClient)
		assert.NoError(t, err, "Should not return error when creating Jusibe instance with NewWithHTTPClient function")

		mockRoundTripper.EXPECT().RoundTrip(gomock.AssignableToTypeOf(&http.Request{})).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, fmt.Sprintf("%ssend_sms?to=%s&from=%s&message=%s", baseURL, to, from, message), req.URL.String())
			res := &http.Response{}
			res.StatusCode = 200
			bodyBytes := []byte(`{
				"status": "Sent",
				"message_id": "xyz123",
				"sms_credits_used": 1
			}`)
			res.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return res, nil
		})

		ctx := context.Background()
		res, code, err := jusibe.SendSMS(ctx, to, from, message)

		assert.Equal(t, 200, code)
		assert.NoError(t, err)
		assert.Equal(t, string(StatusSMSSent), res.Status)
		assert.Equal(t, "xyz123", res.MessageID)
		assert.Equal(t, 1, res.SMSCreditsUsed)
	})

	t.Run("CheckSMSCredits", func(t *testing.T) {
		accessToken, publicKey, baseURL := "some_access_token", "some_public_key", "https://jusibe.com/"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey, APIBaseURL: baseURL}

		mockController := gomock.NewController(t)
		mockRoundTripper := mocks.NewMockRoundTripper(mockController)

		httpClient := &http.Client{Transport: mockRoundTripper}

		jusibe, err := NewWithHTTPClient(cfg, httpClient)
		assert.NoError(t, err, "Should not return error when creating Jusibe instance with NewWithHTTPClient function")

		mockRoundTripper.EXPECT().RoundTrip(gomock.AssignableToTypeOf(&http.Request{})).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, fmt.Sprintf("%sget_credits", baseURL), req.URL.String())
			res := &http.Response{}
			res.StatusCode = 200
			bodyBytes := []byte(`{
				"sms_credits": "100"
			}`)
			res.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return res, nil
		})

		ctx := context.Background()
		res, code, err := jusibe.CheckSMSCredits(ctx)

		assert.Equal(t, 200, code)
		assert.NoError(t, err)
		assert.Equal(t, "100", res.SMSCredits)
	})

	t.Run("CheckSMSCredits", func(t *testing.T) {
		accessToken, publicKey, baseURL := "some_access_token", "some_public_key", "https://jusibe.com/"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey, APIBaseURL: baseURL}

		mockController := gomock.NewController(t)
		mockRoundTripper := mocks.NewMockRoundTripper(mockController)

		httpClient := &http.Client{Transport: mockRoundTripper}

		jusibe, err := NewWithHTTPClient(cfg, httpClient)
		assert.NoError(t, err, "Should not return error when creating Jusibe instance with NewWithHTTPClient function")

		mockRoundTripper.EXPECT().RoundTrip(gomock.AssignableToTypeOf(&http.Request{})).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, fmt.Sprintf("%sdelivery_status?message_id=xyz123", baseURL), req.URL.String())
			res := &http.Response{}
			res.StatusCode = 200
			bodyBytes := []byte(`{
				"message_id": "xyz123",
				"status": "Delivered",
				"date_sent": "2015-05-19 04:34:48",
				"date_delivered": "2015-05-19 04:35:05"
			}`)
			res.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return res, nil
		})

		ctx := context.Background()
		res, code, err := jusibe.CheckSMSDeliveryStatus(ctx, "xyz123")

		assert.Equal(t, 200, code)
		assert.NoError(t, err)
		assert.Equal(t, string(StatusSMSDelivered), res.Status)
		assert.Equal(t, "xyz123", res.MessageID)
		assert.Equal(t, "2015-05-19 04:34:48", res.DateSent)
		assert.Equal(t, "2015-05-19 04:35:05", res.DateDelivered)
	})
}
