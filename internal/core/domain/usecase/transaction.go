package usecase

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/api"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/spi"
	"github.com/shopspring/decimal"
)

//TransactionServiceConfig config struct used as param for NewTransactionService.
type TransactionServiceConfig struct {
	Log           spi.Logger
	TransactionDB spi.TransactionRepository
	ImportDB      spi.ImportRepository
	Validator     spi.Validator
}

//TransactionService struct that implements all Transaction api port.
type TransactionService struct {
	log       spi.Logger
	tdb       spi.TransactionRepository
	idb       spi.ImportRepository
	validator spi.Validator
}

//NewTransactionService creates a new TransactionService
func NewTransactionService(c TransactionServiceConfig) *TransactionService {
	return &TransactionService{
		log:       c.Log,
		tdb:       c.TransactionDB,
		idb:       c.ImportDB,
		validator: c.Validator,
	}
}

//ImportTransactionsFile imports a .csv file with one or more transactions of the same type.
func (s *TransactionService) ImportTransactionsFile(r *api.ImportTransactionsFileRequest) (*api.ImportTransactionsFileResponse, error) {
	var firstTransactionDate time.Time
	var countSuccess int
	var countTotal int
	var res *api.ImportTransactionsFileResponse

	csv := csv.NewReader(r.FileReader)

	first := true
	for {
		countTotal++
		record, err := csv.Read()
		if first {
			if err == io.EOF {
				return &api.ImportTransactionsFileResponse{
					Status:  api.WARNING,
					Details: "empty file",
				}, errors.New("empty file")
			}

			t, err := s.newTransaction(record)
			if err != nil {
				return &api.ImportTransactionsFileResponse{
					Status:  api.ERROR,
					Details: "first transaction is invalid",
				}, err
			}

			firstTransactionDate = t.TransactionDatetime

			first = false
		}
		if err == io.EOF {
			countTotal--
			res = &api.ImportTransactionsFileResponse{
				Status:                api.OK,
				Details:               "end of file",
				TotalProcessedRecords: countTotal,
				TotalValidRecords:     countSuccess,
			}
			break
		}

		t, err := s.newTransaction(record)
		if err != nil {
			s.log.Error("error while creating Transaction, details: ", err)
			continue
		}

		if t.TransactionDatetime.Format("2006-01-02") != firstTransactionDate.Format("2006-01-02") {
			s.log.Error("invalid date for this record")
			continue
		}

		ok, err := s.idb.HasImportByTransactionDate(t.TransactionDatetime)
		if err != nil {
			s.log.Error("error while checking if transaction exists, details", err)
			continue
		}
		if ok {
			return &api.ImportTransactionsFileResponse{
				Status:  api.WARNING,
				Details: "transactions already imported",
			}, errors.New("transactions already imported")
		} else {
			err := s.tdb.SaveTransaction(t)
			if err != nil {
				s.log.Error("error while saving transaction, details", err)
				continue
			}
			countSuccess++
			s.log.Info("saved transaction", t)
		}

	}

	us := &model.Import{
		TransactionDate: firstTransactionDate,
		ImportDate:      time.Now(),
		Filename:        r.Filename,
		FileSizeInMB:    r.FileSizeInMB,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	ok, err := s.idb.HasImportByTransactionDate(us.TransactionDate)
	if err != nil {
		return &api.ImportTransactionsFileResponse{
			Status:  api.ERROR,
			Details: "error while checking upload status",
		}, err
	}
	if !ok {
		if err := s.idb.SaveImport(us); err != nil {
			return &api.ImportTransactionsFileResponse{
				Status:  api.ERROR,
				Details: "error while saving upload status",
			}, err
		}
	}

	return res, nil
}

//ListImports lists all file imports
func (s *TransactionService) ListImports() ([]api.ListImportsResponse, error) {
	imports, err := s.idb.ListImports()
	if err != nil {
		return nil, err
	}

	res := []api.ListImportsResponse{}

	for _, i := range imports {
		res = append(res, api.ListImportsResponse{ImportDate: i.ImportDate, TransactionsImportDate: i.TransactionDate})
	}
	return res, nil
}

//newTransaction creates a new transaction struct with validaded required fields.
func (s *TransactionService) newTransaction(record []string) (*model.Transaction, error) {

	if len(record) != 8 {
		return nil, errors.New("invalid record, length mismatch")
	}

	amount, err := decimal.NewFromString(record[6])
	if err != nil {
		return nil, errors.New("error while converting amount")
	}

	transactionDate, err := time.Parse("2006-01-02T15:04:05", record[7])
	if err != nil {
		return nil, errors.New("error while converting transaction date")
	}

	t := &model.Transaction{
		SourceBank:          record[0],
		SourceAgency:        record[1],
		SourceAccount:       record[2],
		TargetBank:          record[3],
		TargetAgency:        record[4],
		TargetAccount:       record[5],
		TransactionAmount:   amount,
		TransactionDatetime: transactionDate,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if m, _ := s.validator.Validate(t); len(m) != 0 {
		return nil, fmt.Errorf("invalid fiels for this records, details, %+v", m)
	}

	return t, nil
}
