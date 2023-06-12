package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/h-varmazyar/gopet/national_id"
	"github.com/h-varmazyar/gopet/phone"
	drivingLicenceRepo "github.com/h-varmazyar/insurate/internal/core/repository/drivingLicence"
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
		log.WithError(err).Error("driving licence repository failed")
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

	return service, nil
}

type NewScoreReq struct {
	DrivingLicenceNumber string
	NationalCode         string
	Mobile               string
	PlateAlphabet        string
	PlateStart           int8
	PlateEnd             int8
	PlateRegion          int8
}

type NewScoreResp struct {
	ScoreID uuid.UUID
}

type DownloadScoreReq struct {
	ScoreID string
}

type DownloadScoreResp struct {
	FilePath string
	Score    *scoreRepo.Score
}

func (s *Service) NewScore(ctx context.Context, req *NewScoreReq) (*NewScoreResp, error) {
	if !national_id.Validate(req.NationalCode) {
		return nil, errors.New("invalid_national_code")
	}

	var (
		score                *scoreRepo.Score
		err                  error
		scoreCalculateParams = new(ScoreCalculateParams)
	)
	score, err = s.scoreRepository.Last(ctx, req.NationalCode)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//todo: read time from config
	if score != nil && score.CreatedAt.After(time.Now().Add(-1*time.Hour*24*30)) {
		return &NewScoreResp{ScoreID: score.ID}, nil
	} else {
		score = &scoreRepo.Score{
			NationalCode: req.NationalCode,
			Status:       scoreRepo.Pending,
		}
		err = s.scoreRepository.Create(ctx, score)
		if err != nil {
			return nil, err
		}
	}
	defer func() {
		if err != nil && score != nil {
			_ = s.scoreRepository.UpdateStatus(ctx, score.ID, scoreRepo.Failed)
		}
	}()

	scoreCalculateParams.Person, err = s.preparePerson(ctx, req.Mobile, req.NationalCode)
	if err != nil {
		return nil, err
	}

	scoreCalculateParams.Plate, err = s.preparePlate(ctx, scoreCalculateParams.Person, req)
	if err != nil {
		return nil, err
	}

	scoreCalculateParams.DrivingLicence, err = s.prepareDrivingLicence(ctx, scoreCalculateParams.Person, req.DrivingLicenceNumber)
	if err != nil {
		return nil, err
	}

	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Minute)
	go s.scoringAsyncProcesses(timeoutCtx, score, scoreCalculateParams)

	return &NewScoreResp{ScoreID: score.ID}, nil
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

func (s *Service) scoringAsyncProcesses(ctx context.Context, score *scoreRepo.Score, params *ScoreCalculateParams) {
	var err error
	defer func() {
		if err != nil {
			_ = s.scoreRepository.UpdateStatus(ctx, score.ID, scoreRepo.Failed)
		}
	}()
	_ = s.scoreRepository.UpdateStatus(ctx, score.ID, scoreRepo.PreparingData)
	params.Offences, err = s.getDrivingOffence(ctx, params.Plate, params.Person)
	if err != nil {
		log.WithError(err).Error("failed to get driving offences")
		return
	}

	params.NegativeScore, err = s.getNegativeScore(ctx, params.Person, params.DrivingLicence)
	if err != nil {
		log.WithError(err).Error("failed to get negative score")
		return
	}
	_ = s.scoreRepository.UpdateStatus(ctx, score.ID, scoreRepo.Calculating)
	score.Value = params.CalculateScore(ctx)
	score.Status = scoreRepo.Done
	err = s.scoreRepository.Update(ctx, score)
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
		plate.Alphabet = req.PlateAlphabet
		plate.StartNumber = req.PlateStart
		plate.EndNumber = req.PlateEnd
		plate.RegionCode = req.PlateRegion
		plate.Person = person
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
		licence.Number = licenceNumber
		licence.Person = person
		err = s.drivingLicenceRepository.Create(ctx, licence)
	}
	if err != nil {
		return nil, err
	}
	return licence, nil
}

func (s *Service) getDrivingOffence(ctx context.Context, plate *plateRepo.Plate, person *personRepo.Person) (*finnotech.DrivingOffenceResult, error) {
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

	offences, err := s.finnoTechClient.DrivingOffence(ctx, drivingOffenceReq)
	if err != nil {
		return nil, err
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
