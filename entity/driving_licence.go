package entity

import "time"

type LicenceType int8
type LicenceStatus int8

const (
	LicenceMotorCycle LicenceType = iota
	LicenceMotorCycle200
	LicenceNormalCar
	LicenceMiddleCar
	LicenceLargeCar
)

const (
	LicenceStatusAllowed LicenceStatus = iota
	LicenceStatusNotAllowed
)

type DrivingLicence struct {
	Person         *Person
	Number         int32
	ExpirationTime time.Time
	AllowedVehicle []AllowedLicence
	NegativeScore  int8
	OffenceCount   int16
	Rule           string
	IssuedDate     time.Time
}

type AllowedLicence struct {
	Type       LicenceType
	Status     LicenceStatus
	IssuedDate time.Time
}
