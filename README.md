# jusibe-go

[![Build Status](https://travis-ci.org/azeezolaniran2016/jusibe-go.svg?branch=master)](https://travis-ci.org/azeezolaniran2016/jusibe-go)
[![codecov](https://codecov.io/gh/azeezolaniran2016/jusibe-go/branch/master/graph/badge.svg)](https://codecov.io/gh/azeezolaniran2016/jusibe-go)

> Jusibe API Go package

A Go package which wraps [Jusibe](http://jusibe.com) API

## Usage

```go
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

// Send Bulk SMS
to, from, message := "08000000000000,08050000000,08090000000", "Azeez", "Hello World"
bulkSMSResponse, _, err := j.SendBulkSMS(context.Background(), to, from, message)
if err != nil {
  log.Fatal(err)
}

// Check Bulk SMS Status
status, _, err := j.CheckBulkSMSStatus(context.Background(), bulkSMSResponse.BulkMessageID)
if err != nil {
  log.Fatal(err)
}
fmt.Printf("%+v\n", status)

// Get SMS credits
creditsResponse, _, err := j.CheckSMSCredits(context.Background())
if err != nil {
  log.Fatal("err")
}
fmt.Printf("%+v\n", creditsResponse)
```

## Contributing

To contribute to this work:

1. Fork it [here](https://github.com/azeezolaniran2016/jusibe-go)
2. Create your feature branch `git checkout -b my-new-feature`
3. Commit your changes `git commit -am 'Add some feature'`
4. Push to the branch `git push origin my-new-feature`
5. Create a Pull Request

or [create an issue](https://github.com/azeezolaniran2016/jusibe-go/issues)


## License

The package is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).