package finnotech

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type NegativeScoreReq struct {
	LicenceNumber string
	NationalCode  string
	Mobile        string
}

func (c *Client) NegativeScore(ctx context.Context, negativeScoreReq *NegativeScoreReq) (int8, error) {
	scoreUrl := fmt.Sprintf("%v/billing/v2/clients/%v/negativeScore", c.BaseUrl, c.ID)
	req, err := http.NewRequest(http.MethodGet, scoreUrl, nil)
	if err != nil {
		return 0, err
	}

	queryParams := url.Values{
		"licenseNumber": []string{fmt.Sprint(negativeScoreReq.LicenceNumber)},
		"nationalID":    []string{negativeScoreReq.NationalCode},
		"mobile":        []string{negativeScoreReq.Mobile},
	}

	c.setQueryParams(req, queryParams)

	res := new(NegativeScore)
	code, err := c.doRequest(req, res)
	if err != nil {
		return 0, err
	}
	if res.Status != "DONE" {
		return 0, handleFinnotechError(ctx, res.Error)
	}
	if code == http.StatusOK {
		score, err := strconv.Atoi(res.Result.NegativeScore)
		if err != nil {
			return 0, errors.New("failed to parse negative score")
		}
		return int8(score), nil
	}
	return 0, errors.New(http.StatusText(code))
}
