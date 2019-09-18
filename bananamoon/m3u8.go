package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"

	m3u8 "github.com/poccariswet/m3u8-decoder"
)

func (c *Client) CreateM3U8Playlist(ctx context.Context, stationID string, start time.Time) (string, error) {
	program, err := c.GetStartTime(ctx, stationID, start)
	if err != nil {
		return "", err
	}

	m3u8url := path.Join("v2", "api", "ts/playlist.m3u8")
	u := *c.URL
	u.Path = path.Join(c.URL.Path, m3u8url)
	query := u.Query()
	query.Set("station_id", stationID)
	query.Set("ft", program.Ft)
	query.Set("to", program.To)
	query.Set("l", "15")
	u.RawQuery = query.Encode()

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("pragma", "no-cache")
	req.Header.Set(radikoAuthTokenHeader, c.authTokenHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	playlist, err := m3u8.DecodeFrom(resp.Body)
	fmt.Println(playlist)
	if err != nil || !playlist.Master() {
		return "", err
	}

	plist := playlist.Segments[0].(*m3u8.VariantSegment)
	fmt.Println(plist)
	if plist == nil {
		return "", errors.New("invalid m3u8 format")
	}

	return plist.URI, nil
}

// list　の生成
func Getlist(uri string) ([]string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	playlist, err := m3u8.DecodeFrom(resp.Body)
	if err != nil || playlist.Master() {
		return nil, err
	}

	var list []string
	for _, val := range playlist.Segments {
		v, ok := val.(*m3u8.InfSegment)
		if v != nil && ok {
			list = append(list, v.URI)
		}
	}
	return list, nil
}
