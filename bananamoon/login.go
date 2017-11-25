package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// ステータスはログインによって異なるため、インターフェースを使って適宜返す
type State interface {
	StatusCode() string
}

// メソッドを使ってステータスだけを返せるようにするために構造体を定義
type LoginState struct {
	Status string `json:"status"`
}

func (l *LoginState) StatusCode() string {
	return l.Status
}

type Success struct {
	Areafree string `json:"areafree"`
	*LoginState
	User_Key    string `json:"user_key"`
	Paid_member string `json:"paid_member"`
}

type NotSuccess struct {
	Message string `json:"message"`
	*LoginState
	Cause string `json:"cause"`
}

func (c *Client) Login(ctx context.Context, mail, pass string) (State, error) {
	loginEndpoint := "ap/member/login/login"
	v := url.Values{}
	v.Set("mail", mail)
	v.Set("pass", pass)

	login := *c.URL
	login.Path = path.Join(c.URL.Path, loginEndpoint)
	req, err := http.NewRequest("POST", login.String(), strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	_, _ = ioutil.ReadAll(resp.Body)

	return c.check(ctx)
}

func (c *Client) check(ctx context.Context) (State, error) {
	loginCheckpoint := "ap/member/webapi/member/login/check"

	url := *c.URL
	url.Path = path.Join(c.URL.Path, loginCheckpoint)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("pragma", "no-cache")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		status := Success{}
		if err := json.Unmarshal(respBody, &status); err != nil {
			return nil, err
		}
		return status, nil
	}

	status := NotSuccess{}
	if err := json.Unmarshal(respBody, &status); err != nil {
		return nil, err
	}
	return status, nil
}
