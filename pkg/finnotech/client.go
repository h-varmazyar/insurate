package finnotech

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	NID       string `yaml:"NID"`
	BaseUrl   string `yaml:"baseUrl"`
	ClientID  string `yaml:"clientID"`
	AuthToken string `yaml:"authToken"`
}

type Client struct {
	configs         *Config
	BearerToken     string
	BearerExpiredAt time.Time
}

func NewClient(ctx context.Context, configs *Config) (*Client, error) {
	client := &Client{
		configs:     configs,
		BearerToken: "",
	}

	err := client.Authorize(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

func (c *Client) basicAuthorize(req *http.Request) {
	encodedToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", c.configs.ClientID, c.configs.AuthToken)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %v", encodedToken))
}

func (c *Client) bearerAuthorize(req *http.Request) {
	if c.BearerExpiredAt.Before(time.Now()) {
		//refresh token
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.BearerToken))
}

func (c *Client) setQueryParams(req *http.Request, values url.Values) {
	if len(values) > 0 {
		req.URL.RawQuery = values.Encode()
	}
}

func (c *Client) doRequest(req *http.Request, res interface{}) (int, error) {
	c.setHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, err
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return response.StatusCode, err
}

//todo: must be handle errors
func handleFinnotechError(_ context.Context, err *Error) error {
	if err != nil {
		if err.Message != "" {
			return errors.New(err.Message)
		} else {
			return errors.New(err.Code)
		}
	}
	return errors.New("invalid status")
}
