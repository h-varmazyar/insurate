package finnotech

import (
	"context"
	"errors"
	"fmt"
	"github.com/h-varmazyar/insurate/entity"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) NegativeScore(ctx context.Context, licence *entity.DrivingLicence) (int8, error) {
	if licence.Person == nil {
		return 0, errors.New("nil person not acceptable in licence")
	}

	scoreUrl := fmt.Sprintf("%v/billing/v2/clients/%v/negativeScore", c.BaseUrl, c.ID)
	req, err := http.NewRequest(http.MethodGet, scoreUrl, nil)
	if err != nil {
		return 0, err
	}

	queryParams := url.Values{
		"licenseNumber": []string{fmt.Sprint(licence.Number)},
		"nationalID":    []string{licence.Person.NationalCode},
		"mobile":        []string{licence.Person.Mobile},
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
		licence.NegativeScore = int8(score)
		return int8(score), nil
	}
	return 0, errors.New(http.StatusText(code))
}
