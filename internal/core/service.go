package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/h-varmazyar/gopet/national_id"
	"github.com/h-varmazyar/gopet/phone"
	drivingLicenceRepo "github.com/h-varmazyar/insurate/internal/core/repository/drivingLicence"
	"github.com/h-varmazyar/insurate/internal/core/repository/drivingOffence"
	personRepo "github.com/h-varmazyar/insurate/internal/core/repository/person"
	plateRepo "github.com/h-varmazyar/insurate/internal/core/repository/plate"
	scoreRepo "github.com/h-varmazyar/insurate/internal/core/repository/score"
	db "github.com/h-varmazyar/insurate/internal/pkg/db/PostgreSQL"
	"github.com/h-varmazyar/insurate/pkg/finnotech"
	platePkg "github.com/h-varmazyar/insurate/pkg/plate"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Service struct {
	drivingLicenceRepository drivingLicenceRepo.Repository
	personRepository         personRepo.Repository
	plateRepository          plateRepo.Repository
	scoreRepository          scoreRepo.Repository
	finnoTechClient          *finnotech.Client
	log                      *log.Logger
}

func NewService(ctx context.Context, log *log.Logger, db *db.DB, finnotechClient *finnotech.Client) (*Service, error) {
	service := &Service{
		finnoTechClient: finnotechClient,
		log:             log,
	}
	var err error

	service.drivingLicenceRepository, err = drivingLicenceRepo.NewRepository(ctx, log, db)
	if err != nil {
		log.WithError(err).Error("driving licence repository failed")
		return nil, err
	}

	service.personRepository, err = personRepo.NewRepository(ctx, log, db)
	if err != nil {
		log.WithError(err).Error("person repository failed")
		return nil, err
	}

	service.plateRepository, err = plateRepo.NewRepository(ctx, log, db)
	if err != nil {
		log.WithError(err).Error("driving licence repository failed")
		return nil, err
	}

	service.personRepository, err = personRepo.NewRepository(ctx, log, db)
	if err != nil {
		log.WithError(err).Error("driving licence repository failed")
		return nil, err
	}

	service.scoreRepository, err = scoreRepo.NewRepository(ctx, log, db)
	if err != nil {
		log.WithError(err).Error("score repository failed")
		return nil, err
	}

	return service, nil
}

type NewScoreReq struct {
	DrivingLicenceNumber string `json:"driving_licence_number"`
	NationalCode         string `json:"national_code"`
	Mobile               string `json:"mobile"`
	PlateAlphabet        string `json:"plate_alphabet"`
	PlateStart           int8   `json:"plate_start"`
	PlateEnd             int16  `json:"plate_end"`
	PlateRegion          int8   `json:"plate_region"`
}

type NewScoreResp struct {
	Score *scoreRepo.Score
}

type DownloadScoreReq struct {
	ScoreID string
}

type DownloadScoreResp struct {
	FilePath string
	Score    *scoreRepo.Score
}

func (s *Service) NewScore(ctx context.Context, req *NewScoreReq) (*NewScoreResp, error) {
	var (
		score *scoreRepo.Score
		err   error
	)

	if !national_id.Validate(req.NationalCode) {
		return nil, errors.New("invalid_national_code")
	}

	score, err = s.scoreRepository.Last(ctx, req.NationalCode)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//todo: read time from config
	if score != nil && score.CreatedAt.After(time.Now().Add(-1*time.Second*20)) {
		return &NewScoreResp{Score: score}, nil
	} else {
		score = new(scoreRepo.Score)
		score.NationalCode = req.NationalCode
	}
	defer func() {
		if err != nil && score != nil {
			_ = s.scoreRepository.UpdateStatus(ctx, score.ID, scoreRepo.Failed)
		}
	}()

	score.Person, err = s.preparePerson(ctx, req.Mobile, req.NationalCode)
	if err != nil {
		return nil, err
	}

	score.Plate, err = s.preparePlate(ctx, score.Person, req)
	if err != nil {
		return nil, err
	}

	score.DrivingLicence, err = s.prepareDrivingLicence(ctx, score.Person, req.DrivingLicenceNumber)
	if err != nil {
		return nil, err
	}

	err = s.scoringAsyncProcesses(ctx, score)
	if err != nil {
		return nil, err
	}

	return &NewScoreResp{Score: score}, nil
}

func (s *Service) DownloadScore(ctx context.Context, req *DownloadScoreReq) (*DownloadScoreResp, error) {
	scoreID, err := uuid.Parse(req.ScoreID)
	if err != nil {
		return nil, err
	}
	score, err := s.scoreRepository.ReturnByID(ctx, scoreID)
	if err != nil {
		return nil, err
	}
	resp := &DownloadScoreResp{Score: score}
	return resp, nil
}

func (s *Service) scoringAsyncProcesses(ctx context.Context, score *scoreRepo.Score) error {
	var err error
	defer func() {
		if err != nil {
			_ = s.scoreRepository.UpdateStatus(ctx, score.ID, scoreRepo.Failed)
		}
	}()
	_ = s.scoreRepository.UpdateStatus(ctx, score.ID, scoreRepo.PreparingData)
	score.DrivingOffences, err = s.getDrivingOffence(ctx, score.Plate, score.Person)
	if err != nil {
		log.WithError(err).Error("failed to get driving offences")
		return err
	}

	score.NegativeScore, err = s.getNegativeScore(ctx, score.Person, score.DrivingLicence)
	if err != nil {
		log.WithError(err).Error("failed to get negative score")
		return err
	}
	_ = s.scoreRepository.UpdateStatus(ctx, score.ID, scoreRepo.Calculating)
	//todo: scoreCalculator must be removed
	params := &ScoreCalculateParams{score}
	score.Value = params.CalculateScore(ctx)
	score.Status = scoreRepo.Done
	err = s.scoreRepository.Create(ctx, score)
	if err != nil {
		s.log.WithError(err).Error("failed to update score")
		return err
	}
	return nil
}

func (s *Service) preparePerson(ctx context.Context, mobileNumber, nationalCode string) (*personRepo.Person, error) {
	mobile, err := phone.GetPhoneNumberDetails(mobileNumber)
	if err != nil {
		return nil, err
	}

	var person *personRepo.Person
	person, err = s.personRepository.Return(ctx, nationalCode)
	if err != nil && err == gorm.ErrRecordNotFound {
		person = &personRepo.Person{
			NationalCode: nationalCode,
			Mobile:       mobile.FullNumber,
		}
		err = s.personRepository.Create(ctx, person)
	}
	if err != nil {
		return nil, err
	}
	return person, nil
}

func (s *Service) preparePlate(ctx context.Context, person *personRepo.Person, req *NewScoreReq) (*plateRepo.Plate, error) {
	var err error
	plateText := ""
	{
		plate := platePkg.Plate{
			Alphabet:    req.PlateAlphabet,
			StartNumber: req.PlateStart,
			EndNumber:   req.PlateEnd,
			RegionCode:  req.PlateRegion,
		}
		plateText, err = plate.Format(platePkg.CodedAlphabet)
		if err != nil {
			return nil, err
		}
	}

	var plate *plateRepo.Plate
	plate, err = s.plateRepository.ReturnByText(ctx, plateText)
	if err != nil && err == gorm.ErrRecordNotFound {
		plate = &plateRepo.Plate{
			Alphabet:    req.PlateAlphabet,
			StartNumber: req.PlateStart,
			EndNumber:   req.PlateEnd,
			RegionCode:  req.PlateRegion,
			Person:      person,
		}

		err = s.plateRepository.Create(ctx, plate)
	}
	if err != nil {
		return nil, err
	}
	return plate, nil
}

func (s *Service) prepareDrivingLicence(ctx context.Context, person *personRepo.Person, number string) (*drivingLicenceRepo.DrivingLicence, error) {
	licenceNumber, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		return nil, err
	}
	licence, err := s.drivingLicenceRepository.ReturnByNumber(ctx, licenceNumber)
	if err != nil && err == gorm.ErrRecordNotFound {
		licence := &drivingLicenceRepo.DrivingLicence{
			Number: licenceNumber,
			Person: person,
		}
		err = s.drivingLicenceRepository.Create(ctx, licence)
	}
	if err != nil {
		return nil, err
	}
	return licence, nil
}

func (s *Service) getDrivingOffence(ctx context.Context, plate *plateRepo.Plate, person *personRepo.Person) ([]*drivingOffence.DrivingOffence, error) {
	drivingOffenceReq := &finnotech.DrivingOffenceReq{
		NationalCode: person.NationalCode,
		Mobile:       person.Mobile,
		Plate: &finnotech.Plate{
			Alphabet:    plate.Alphabet,
			StartNumber: plate.StartNumber,
			EndNumber:   plate.EndNumber,
			RegionCode:  plate.RegionCode,
		},
	}

	finnotechOffences, err := s.finnoTechClient.DrivingOffence(ctx, drivingOffenceReq)
	if err != nil {
		return nil, err
	}

	offences := make([]*drivingOffence.DrivingOffence, 0)
	for _, bill := range finnotechOffences.Bills {
		offence := &drivingOffence.DrivingOffence{
			ID:             bill.ID,
			Type:           bill.Type,
			Description:    bill.Description,
			Code:           bill.Code,
			Price:          bill.Price,
			City:           bill.City,
			Location:       bill.Location,
			Date:           bill.Date,
			PlateCode:      bill.PlateCode,
			DataValue:      bill.DataValue,
			Barcode:        bill.Barcode,
			Plate:          plate,
			BillID:         bill.BillID,
			PaymentID:      bill.PaymentID,
			NormalizedDate: bill.NormalizedDate,
			IsPayable:      bill.IsPayable,
			PolicemanCode:  bill.PolicemanCode,
			HasImage:       bill.HasImage,
		}
		offences = append(offences, offence)
	}
	return offences, nil
}

func (s *Service) getNegativeScore(ctx context.Context, person *personRepo.Person, drivingLicence *drivingLicenceRepo.DrivingLicence) (int8, error) {
	drivingLicenceReq := &finnotech.NegativeScoreReq{
		NationalCode:  person.NationalCode,
		Mobile:        person.Mobile,
		LicenceNumber: fmt.Sprint(drivingLicence.Number),
	}

	negativeScore, err := s.finnoTechClient.NegativeScore(ctx, drivingLicenceReq)
	if err != nil {
		return 0, err
	}
	return negativeScore, nil
}

func (s *Service) fillScoreData(ctx context.Context, score *scoreRepo.Score, params *ScoreCalculateParams) {

	score.Status = scoreRepo.Done
}
