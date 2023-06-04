package finnotech

import (
	"context"
	"errors"
	"fmt"
	"github.com/h-varmazyar/insurate/entity"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) DrivingOffence(ctx context.Context, person *entity.Person, plate *entity.Plate) ([]*entity.DrivingOffence, error) {
	if person == nil {
		return nil, errors.New("nil person not acceptable in driving offence")
	}
	if plate == nil {
		return nil, errors.New("nil plate not acceptable in driving offence")
	}

	scoreUrl := fmt.Sprintf("%v/billing/v2/clients/%v/drivingOffense", c.BaseUrl, c.ID)
	req, err := http.NewRequest(http.MethodGet, scoreUrl, nil)
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{
		"version":     []string{"2"},
		"plateNumber": []string{generatePlateCode(plate)},
		"nationalID":  []string{person.NationalCode},
		"mobile":      []string{person.Mobile},
	}

	c.setQueryParams(req, queryParams)

	res := new(DrivingOffence)
	code, err := c.doRequest(req, res)
	if err != nil {
		return nil, err
	}
	if res.Status != "DONE" {
		return nil, handleFinnotechError(ctx, res.Error)
	}
	if code == http.StatusOK {
		offences := generateOffenceResponse(res.Result, plate)
		return offences, nil
	}
	return nil, errors.New(http.StatusText(code))
}

func generateOffenceResponse(finnotechOffence *DrivingOffenceResult, plate *entity.Plate) []*entity.DrivingOffence {
	offences := make([]*entity.DrivingOffence, 0)
	for _, bill := range finnotechOffence.Bills {
		normalizedDate, err := time.Parse("2006-01-02 15:04:05", bill.NormalizedDate)
		if err != nil {
			normalizedDate = time.Unix(0, 0)
		}
		offence := &entity.DrivingOffence{
			ID:             bill.ID,
			Type:           bill.Type,
			Description:    bill.Description,
			Code:           bill.Code,
			Price:          bill.Price,
			City:           bill.City,
			Location:       bill.Location,
			Date:           bill.Date,
			PlateCode:      bill.Serial,
			DataValue:      bill.DataValue,
			Barcode:        bill.Barcode,
			Plate:          plate,
			BillID:         bill.BillID,
			PaymentID:      bill.PaymentID,
			NormalizedDate: normalizedDate,
			IsPayable:      bill.IsPayable,
			PolicemanCode:  bill.PolicemanCode,
			HasImage:       bill.HasImage,
		}
		offences = append(offences, offence)
	}
	return offences
}
