package vehicle

import plateRepo "github.com/h-varmazyar/insurate/internal/core/repository/plate"

type VehicleType int8

//complete all types
const (
	MotorCycle = iota
	MotorCycle200
	CarPersonal
	CarTaxi
	CarMiddle
)

type Vehicle struct {
	Plate         *plateRepo.Plate
	Type          VehicleType
	AxleCount     int8
	CylinderCount int8
	EngineVolume  int16
}
