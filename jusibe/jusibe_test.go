package jusibe

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/azeezolaniran2016/jusibe-go/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
	})

	t.Run("New should set fields on Jusibe instance", func(t *testing.T) {
		accessToken, publicKey := "some_access_token", "some_public_key"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey}
		jusibe, err := New(cfg)

		assert.NoError(t, err, "Should not return error when creating Jusibe instance with New function")

		assert.Equal(t, publicKey, jusibe.publicKey, "Should set publicKey")

		assert.Equal(t, accessToken, jusibe.accessToken, "Should set accessToken")

		assert.NotNil(t, jusibe.httpClient, "Should set httpClient to default")

		assert.Equal(t, defaultHTTPClientTimeout, jusibe.httpClient.Timeout, "Should set httpClient timeout to default")
	})

	t.Run("NewWithHTTPClient", func(t *testing.T) {
		accessToken, publicKey := "some_access_token", "some_public_key"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey}

		testTimeout := (10 * time.Second)
		httpClient := &http.Client{Timeout: testTimeout}

		jusibe, err := NewWithHTTPClient(cfg, httpClient)

		assert.NoError(t, err, "Should not return error when creating Jusibe instance with NewWithHTTPClient function")

		assert.NotNil(t, jusibe.httpClient, "Should set httpClient")

		assert.Equal(t, testTimeout, jusibe.httpClient.Timeout, "Should have specified httpClient timeout")

		assert.Equal(t, publicKey, jusibe.publicKey, "Should set publicKey field")

		assert.Equal(t, accessToken, jusibe.accessToken, "Should set accessToken")
	})

	t.Run("SendSMS", func(t *testing.T) {
		accessToken, publicKey := "some_access_token", "some_public_key"
		to, from, message := "09001000101", "test_user", "Hello World!"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey}

		mockController := gomock.NewController(t)
		mockRoundTripper := mocks.NewMockRoundTripper(mockController)

		httpClient := &http.Client{Transport: mockRoundTripper}

		jusibe, err := NewWithHTTPClient(cfg, httpClient)
		assert.NoError(t, err, "Should not return error when creating Jusibe instance with NewWithHTTPClient function")

		mockRoundTripper.EXPECT().RoundTrip(gomock.AssignableToTypeOf(&http.Request{})).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://jusibe.com/smsapi/send_sms?from=test_user&message=Hello+World%21&to=09001000101", req.URL.String())
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
		s, res, err := jusibe.SendSMS(ctx, to, from, message)

		assert.Equal(t, 200, res.StatusCode)
		assert.NoError(t, err)
		assert.Equal(t, string(StatusSMSSent), s.Status)
		assert.Equal(t, "xyz123", s.MessageID)
		assert.Equal(t, 1, s.SMSCreditsUsed)
	})

	t.Run("SendBulkSMS", func(t *testing.T) {
		accessToken, publicKey := "some_access_token", "some_public_key"
		to, from, message := "09001000101,08030000000,09050000000", "test_user", "Hello World!"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey}

		mockController := gomock.NewController(t)
		mockRoundTripper := mocks.NewMockRoundTripper(mockController)

		httpClient := &http.Client{Transport: mockRoundTripper}

		jusibe, err := NewWithHTTPClient(cfg, httpClient)
		assert.NoError(t, err, "Should not return error when creating Jusibe instance with NewWithHTTPClient function")

		mockRoundTripper.EXPECT().RoundTrip(gomock.AssignableToTypeOf(&http.Request{})).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://jusibe.com/smsapi/bulk/send_sms?from=test_user&message=Hello+World%21&to=09001000101%2C08030000000%2C09050000000", req.URL.String())
			res := &http.Response{}
			res.StatusCode = 200
			bodyBytes := []byte(`{
				"status": "Submitted",
				"bulk_message_id": "xeqd6rs3d26"
			}`)
			res.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return res, nil
		})

		ctx := context.Background()
		s, res, err := jusibe.SendBulkSMS(ctx, to, from, message)

		assert.Equal(t, 200, res.StatusCode)
		assert.NoError(t, err)
		assert.Equal(t, string(StatusBulkSMSSubmitted), s.Status)
		assert.Equal(t, "xeqd6rs3d26", s.MessageID)
	})

	t.Run("CheckSMSCredits", func(t *testing.T) {
		accessToken, publicKey := "some_access_token", "some_public_key"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey}

		mockController := gomock.NewController(t)
		mockRoundTripper := mocks.NewMockRoundTripper(mockController)

		httpClient := &http.Client{Transport: mockRoundTripper}

		jusibe, err := NewWithHTTPClient(cfg, httpClient)
		assert.NoError(t, err, "Should not return error when creating Jusibe instance with NewWithHTTPClient function")

		mockRoundTripper.EXPECT().RoundTrip(gomock.AssignableToTypeOf(&http.Request{})).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://jusibe.com/smsapi/get_credits", req.URL.String())
			res := &http.Response{}
			res.StatusCode = 200
			bodyBytes := []byte(`{
				"sms_credits": "100"
			}`)
			res.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return res, nil
		})

		ctx := context.Background()
		sc, res, err := jusibe.CheckSMSCredits(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "100", sc.SMSCredits)
	})

	t.Run("CheckSMSDeliveryStatus", func(t *testing.T) {
		accessToken, publicKey := "some_access_token", "some_public_key"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey}

		mockController := gomock.NewController(t)
		mockRoundTripper := mocks.NewMockRoundTripper(mockController)

		httpClient := &http.Client{Transport: mockRoundTripper}

		jusibe, err := NewWithHTTPClient(cfg, httpClient)
		assert.NoError(t, err, "Should not return error when creating Jusibe instance with NewWithHTTPClient function")

		mockRoundTripper.EXPECT().RoundTrip(gomock.AssignableToTypeOf(&http.Request{})).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://jusibe.com/smsapi/delivery_status?message_id=xyz123", req.URL.String())
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
		ds, res, err := jusibe.CheckSMSDeliveryStatus(ctx, "xyz123")

		assert.Equal(t, 200, res.StatusCode)
		assert.NoError(t, err)
		assert.Equal(t, string(StatusSMSDelivered), ds.Status)
		assert.Equal(t, "xyz123", ds.MessageID)
		assert.Equal(t, "2015-05-19 04:34:48", ds.DateSent)
		assert.Equal(t, "2015-05-19 04:35:05", ds.DateDelivered)
	})

	t.Run("CheckBulkSMSDeliveryStatus", func(t *testing.T) {
		accessToken, publicKey := "some_access_token", "some_public_key"
		cfg := &Config{AccessToken: accessToken, PublicKey: publicKey}

		mockController := gomock.NewController(t)
		mockRoundTripper := mocks.NewMockRoundTripper(mockController)

		httpClient := &http.Client{Transport: mockRoundTripper}

		jusibe, err := NewWithHTTPClient(cfg, httpClient)
		assert.NoError(t, err, "Should not return error when creating Jusibe instance with NewWithHTTPClient function")

		mockRoundTripper.EXPECT().RoundTrip(gomock.AssignableToTypeOf(&http.Request{})).DoAndReturn(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://jusibe.com/smsapi/bulk/status?bulk_message_id=xeqd6rs3d26", req.URL.String())
			res := &http.Response{}
			res.StatusCode = 200
			bodyBytes := []byte(`{
				"bulk_message_id": "xeqd6rs3d26",
				"status": "Completed",
				"created": "2019-04-02 15:23:13",
				"processed": "2019-04-02 15:25:03",
				"total_numbers": "2",
				"total_unique_numbers": "2",
				"total_valid_numbers": "2",
				"total_invalid_numbers": "0"
		}`)
			res.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return res, nil
		})

		ctx := context.Background()
		ds, res, err := jusibe.CheckBulkSMSStatus(ctx, "xeqd6rs3d26")

		assert.Equal(t, 200, res.StatusCode)
		assert.NoError(t, err)
		assert.Equal(t, "xeqd6rs3d26", ds.BulkMessageID)
		assert.Equal(t, "Completed", ds.Status)
		assert.Equal(t, "2019-04-02 15:23:13", ds.Created)
		assert.Equal(t, "2019-04-02 15:25:03", ds.Processed)
		assert.Equal(t, "2", ds.TotalNumbers)
		assert.Equal(t, "2", ds.TotalUniqueNumbers)
		assert.Equal(t, "2", ds.TotalValidNumbers)
		assert.Equal(t, "0", ds.TotalInvalidNumbers)
	})
}
