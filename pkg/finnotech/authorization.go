package finnotech

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) Authorize(ctx context.Context) error {
	body := map[string]string{
		"grant_type": "client_credentials",
		"nid":        c.configs.NID,
		"scopes":     "billing:cc-negative-score:get,billing:driving-offense-inquiry:get,billing:riding-offense-inquiry:get,billing:riding-offense-inquiry:get",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/dev/v2/oauth2/token", c.configs.BaseUrl), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	c.setHeaders(req)
	c.basicAuthorize(req)

	res := new(AuthResponse)
	code, err := c.doRequest(req, res)
	if err != nil {
		return err
	}
	if res.Status != "DONE" {
		return handleFinnotechError(ctx, res.Error)
	}
	if code == http.StatusOK {
		c.BearerToken = res.Result.Value
		c.BearerExpiredAt = time.Now().Add(res.Result.LifeTime)
		return nil
	}
	return errors.New(http.StatusText(code))
}
