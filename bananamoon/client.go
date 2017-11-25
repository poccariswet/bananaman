package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

const (
	EndPoint               = "https://radiko.jp"
	defaultHTTPTimeout     = 100 * time.Second
	radikoAppHeader        = "X-Radiko-App"
	radikoAppVersionHeader = "X-Radiko-App-Version"
	radikoUserHeader       = "X-Radiko-User"
	radikoDeviceHeader     = "X-Radiko-Device"
	radikoAuthTokenHeader  = "X-Radiko-AuthToken"
	radikoKeyLentghHeader  = "X-Radiko-KeyLength"
	radikoKeyOffsetHeader  = "X-Radiko-KeyOffset"
	radikoPartialKeyHeader = "X-Radiko-Partialkey"

	radikoApp        = "pc_ts"
	radikoAppVersion = "4.0.0"
	radikoUser       = "test-stream"
	radikoDevice     = "pc"

	radikoMail = "RADIKO_MAIL"
	radikoPass = "RADIKO_PASS"
)

var (
	httpClient = &http.Client{Timeout: defaultHTTPTimeout}
)

type Client struct {
	URL             *url.URL
	httpClient      *http.Client
	authTokenHeader string
}

func GetClient(ctx context.Context, areaID string) (*Client, error) {
	if httpClient == nil {
		return nil, errors.New("httpClient error : nil")
	}

	jar, err := cookiejar.New(nil) //ないとログイン状態を保てない
	if err != nil {
		return nil, err
	}
	httpClient.Jar = jar

	urlparse, err := url.Parse(EndPoint)
	if err != nil {
		return nil, err
	}

	client := &Client{
		URL:             urlparse,
		httpClient:      httpClient,
		authTokenHeader: "",
	}

	mail := os.Getenv(radikoMail)
	password := os.Getenv(radikoPass)

	login, err := client.Login(ctx, mail, password)
	switch {
	case err != nil:
		return nil, err
	case login.StatusCode() != "200":
		return nil, fmt.Errorf(
			"invalid login status code: %s", login.StatusCode())
	default:
	}

	return client, nil
}

func (c *Client) setAuthTokenHeader(authToken string) {
	c.authTokenHeader = authToken
}
