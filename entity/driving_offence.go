package entity

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
	Plate          *Plate
	BillID         int32
	PaymentID      int32
	NormalizedDate time.Time
	IsPayable      bool
	PolicemanCode  interface{}
	HasImage       bool
}
