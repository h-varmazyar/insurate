package finnotech

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	BaseUrl string
	ID      string
	Token   string
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.Token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
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
		_ = req.Body.Close()
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
func handleFinnotechError(ctx context.Context, err Error) error {
	return errors.New(err.Message)
}
