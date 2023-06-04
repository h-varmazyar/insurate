package entity

type DrivingOffence struct {
	ID             string
	Type           string
	Description    string
	Code           string
	Price          int32
	City           interface{}
	Location       string
	Date           string
	Serial         string
	DataValue      string
	Barcode        interface{}
	License        string
	BillID         int32
	PaymentID      int32
	NormalizedDate string
	IsPayable      bool
	PolicemanCode  interface{}
	HasImage       bool
}
