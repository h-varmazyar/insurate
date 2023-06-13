package finnotech

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type DrivingOffenceReq struct {
	NationalCode string
	Mobile       string
	Plate        *Plate
}

func (c *Client) DrivingOffence(ctx context.Context, drivingOffenceReq *DrivingOffenceReq) (*DrivingOffenceResult, error) {
	if drivingOffenceReq.Plate == nil {
		return nil, errors.New("nil plate not acceptable in driving drivingOffence")
	}

	scoreUrl := fmt.Sprintf("%v/billing/v2/clients/%v/drivingOffense", c.configs.BaseUrl, c.configs.ClientID)
	req, err := http.NewRequest(http.MethodGet, scoreUrl, nil)
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{
		"parameter":   []string{"m16006203"},
		"plateNumber": []string{generatePlateCode(drivingOffenceReq.Plate)},
		"nationalID":  []string{drivingOffenceReq.NationalCode},
		"mobile":      []string{drivingOffenceReq.Mobile},
	}

	c.setQueryParams(req, queryParams)
	c.bearerAuthorize(req)

	res := new(DrivingOffence)
	code, err := c.doRequest(req, res)
	if err != nil {
		return nil, err
	}
	if res.Status != "DONE" {
		return nil, handleFinnotechError(ctx, res.Error)
	}
	if code == http.StatusOK {
		return res.Result, nil
	}
	return nil, errors.New(http.StatusText(code))
}
