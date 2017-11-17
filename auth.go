package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
)

const (
	playerURL = "http://radiko.jp/apps/js/flash/myplayer-release.swf"

	targetID     = 12
	targetCode   = 87
	headerCWS    = 8
	headerRect   = 5
	rectNum      = 4
	headerRest   = 2 + 2
	binaryOffset = 6
)

type Params struct {
	body         io.Reader
	query        map[string]string
	header       map[string]string
	setAuthToken bool
}

func (c *Client) newRequest(ctx context.Context, verb, apiEndpoint string, params *Params) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, apiEndpoint)
	uquery := u.Query()
	for key, val := range params.query {
		uquery.Set(key, val)
	}
	u.RawQuery = uquery.Encode()

	req, err := http.NewRequest(verb, u.String(), params.body)
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		return nil, errors.New("Context is nil")
	}
	req = req.WithContext(ctx)

	for key, val := range params.header {
		req.Header.Set(key, val)
	}
	req.Header.Set("pragma", "no-cache")
	if params.setAuthToken {
		req.Header.Set(radikoAuthTokenHeader, c.authTokenHeader)
	}

	return req, nil
}

func (c *Client) AuthorizeToken(ctx context.Context) (string, error) {
	bin, err := downloadBinary()
	if err != nil {
		return "", err
	}

	f := bytes.NewReader(bin)

	authToken, length, offset, err := c.Auth1Fms(ctx)
	b := make([]byte, length)
	io.CopyN(ioutil.Discard, f, offset)
	if _, err = f.Read(b); err != nil {
		return "", err
	}
	partialKey := base64.StdEncoding.EncodeToString(b)

	slc, err := c.Auth2Fms(ctx, authToken, partialKey)
	if err != nil {
		return "", err
	}
	if err := verifyAuth2FmsResponse(slc); err != nil {
		return "", err
	}

	c.setAuthTokenHeader(authToken)
	return authToken, nil
}

func downloadBinary() ([]byte, error) {
	resp, err := http.Get(playerURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return swfExtract(resp.Body)
}

func swfExtract(rBody io.Reader) ([]byte, error) {
	io.CopyN(ioutil.Discard, rBody, headerCWS)
	zf, err := zlib.NewReader(rBody)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(zf)
	if err != nil {
		return nil, err
	}

	offset := 0

	rectSize := int(buf[offset] >> 3)
	rectOffset := (headerRect + rectNum*rectSize + 7) / 8

	offset += rectOffset

	offset += headerRest

	for i := 0; ; i++ {
		code := int(buf[offset+1])<<2 + int(buf[offset])>>6

		len := int(buf[offset] & 0x3f)

		offset += 2

		if len == 0x3f {
			len = int(buf[offset])
			len += int(buf[offset+1]) << 8
			len += int(buf[offset+2]) << 16
			len += int(buf[offset+3]) << 24

			offset += 4
		}

		if code == 0 {
			return nil, errors.New("swf extract failed")
		}
		id := int(buf[offset]) + int(buf[offset+1])<<8

		if code == targetCode && id == targetID {
			return buf[offset+binaryOffset : offset+len], nil
		}

		offset += len
	}
}

func (c *Client) Auth1Fms(ctx context.Context) (string, int64, int64, error) {
	apiEndpoint := path.Join("v2", "api", "auth1_fms")

	req, err := c.newRequest(ctx, "POST", apiEndpoint, &Params{
		header: map[string]string{
			radikoAppHeader:        radikoApp,
			radikoAppVersionHeader: radikoAppVersion,
			radikoUserHeader:       radikoUser,
			radikoDeviceHeader:     radikoDevice,
		},
	})
	if err != nil {
		return "", 0, 0, err
	}

	resp, err := c.httpClient.Do(req)
	defer resp.Body.Close()

	authToken := resp.Header.Get(radikoAuthTokenHeader)
	keyLength := resp.Header.Get(radikoKeyLentghHeader)
	keyOffset := resp.Header.Get(radikoKeyOffsetHeader)

	length, err := strconv.ParseInt(keyLength, 10, 64)
	if err != nil {
		return "", 0, 0, err
	}
	offset, err := strconv.ParseInt(keyOffset, 10, 64)
	if err != nil {
		return "", 0, 0, err
	}

	return authToken, length, offset, err
}

func (c *Client) Auth2Fms(ctx context.Context, authToken, partialKey string) ([]string, error) {
	apiEndpoint := path.Join("v2", "api", "auth2_fms")

	req, err := c.newRequest(ctx, "POST", apiEndpoint, &Params{
		header: map[string]string{
			radikoAppHeader:        radikoApp,
			radikoAppVersionHeader: radikoAppVersion,
			radikoUserHeader:       radikoUser,
			radikoDeviceHeader:     radikoDevice,
			radikoAuthTokenHeader:  authToken,
			radikoPartialKeyHeader: partialKey,
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := strings.Split(string(b), ",")
	return s, nil
}

func verifyAuth2FmsResponse(slc []string) error {
	if len(slc) == 0 {
		return errors.New("missing token")
	}
	s := strings.TrimSpace(slc[0])
	if !strings.HasPrefix(s, "JP") {
		return fmt.Errorf("invalid token: %s", s)
	}

	return nil
}
