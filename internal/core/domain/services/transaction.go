package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports"
	"github.com/ribeiroelton/alura-challenge-backend/pkg/logger"
	"github.com/ribeiroelton/alura-challenge-backend/pkg/validator"
	"github.com/shopspring/decimal"
)

type TransactionServiceConfig struct {
	Log logger.Logger
	DB  ports.Repository
}

type TransactionService struct {
	log logger.Logger
	db  ports.Repository
}

func NewTransactionService(c TransactionServiceConfig) *TransactionService {
	return &TransactionService{
		log: c.Log,
		db:  c.DB,
	}
}

func (s *TransactionService) ImportTransactionsFile(r *ports.ImportTransactionsFileRequest) (*ports.ImportTransactionsFileResponse, error) {
	var firstTransactionDate time.Time
	var countSuccess int
	var countTotal int
	var res *ports.ImportTransactionsFileResponse

	csv := csv.NewReader(r.FileReader)

	first := true
	for {
		countTotal++
		record, err := csv.Read()
		if first {
			if err == io.EOF {
				return &ports.ImportTransactionsFileResponse{
					Status:  ports.WARNING,
					Details: "empty file",
				}, errors.New("empty file")
			}

			t, err := newTransaction(record)
			if err != nil {
				return &ports.ImportTransactionsFileResponse{
					Status:  ports.ERROR,
					Details: "first transaction is invalid",
				}, err
			}

			firstTransactionDate = t.TransactionDatetime

			first = false
		}
		if err == io.EOF {
			countTotal--
			res = &ports.ImportTransactionsFileResponse{
				Status:                ports.OK,
				Details:               "end of file",
				TotalProcessedRecords: countTotal,
				TotalValidRecords:     countSuccess,
			}
			break
		}

		t, err := newTransaction(record)
		if err != nil {
			s.log.Error("error while creating Transaction, details: ", err)
			continue
		}

		if t.TransactionDatetime.Format("2006-01-02") != firstTransactionDate.Format("2006-01-02") {
			s.log.Error("invalid date for this record")
			continue
		}

		ok, err := s.db.HasImportByTransactionDate(t.TransactionDatetime)
		if err != nil {
			s.log.Error("error while checking if transaction exists, details", err)
			continue
		}
		if ok {
			return &ports.ImportTransactionsFileResponse{
				Status:  ports.WARNING,
				Details: "transactions already imported",
			}, errors.New("transactions already imported")
		} else {
			err := s.db.SaveTransaction(t)
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

	if err := s.db.SaveImport(us); err != nil {
		return &ports.ImportTransactionsFileResponse{
			Status:  ports.ERROR,
			Details: "error while saving upload status",
		}, err
	}

	return res, nil
}

func (s *TransactionService) ListImports() ([]ports.ListImportsResponse, error) {
	imports, err := s.db.ListImports()
	if err != nil {
		return nil, err
	}

	res := []ports.ListImportsResponse{}

	for _, i := range imports {
		res = append(res, ports.ListImportsResponse{ImportDate: i.ImportDate, TransactionsImportDate: i.TransactionDate})
	}
	return res, nil
}

func newTransaction(record []string) (*model.Transaction, error) {

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

	if m, _ := validator.Validate(t); len(m) != 0 {
		return nil, fmt.Errorf("invalid fiels for this records, details, %+v", m)
	}

	return t, nil
}
