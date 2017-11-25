package main

import (
	"context"
	"errors"
	"net/http"
	"path"
	"time"

	"github.com/grafov/m3u8"
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

	playlist, listtype, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil || listtype != m3u8.MASTER {
		return "", err
	}

	plist := playlist.(*m3u8.MasterPlaylist)
	if plist == nil || len(plist.Variants) != 1 || plist.Variants[0] == nil {
		return "", errors.New("invalid m3u8 format")
	}
	return plist.Variants[0].URI, nil
}

// list　の生成
func Getlist(uri string) ([]string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	playlist, listtype, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil || listtype != m3u8.MEDIA {
		return nil, err
	}
	pl := playlist.(*m3u8.MediaPlaylist)

	var list []string
	for _, val := range pl.Segments {
		if val != nil {
			list = append(list, val.URI)
		}
	}
	return list, nil
}
