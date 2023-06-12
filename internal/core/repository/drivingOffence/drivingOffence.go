package drivingOffence

import plateRepo "github.com/h-varmazyar/insurate/internal/core/repository/plate"
import "time"

type DrivingOffence struct {
	ID             string
	Type           string
	Description    string
	Code           string
	Price          int32
	City           interface{}
	Location       string
	Date           string
	PlateCode      string
	DataValue      string
	Barcode        interface{}
	Plate          *plateRepo.Plate
	BillID         int32
	PaymentID      int32
	NormalizedDate time.Time
	IsPayable      bool
	PolicemanCode  interface{}
	HasImage       bool
}
