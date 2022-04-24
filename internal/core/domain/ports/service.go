package ports

import (
	"io"
	"time"
)

type ImportTransactionsFileStatus int

const (
	WARNING ImportTransactionsFileStatus = iota
	ERROR   ImportTransactionsFileStatus = iota
	OK      ImportTransactionsFileStatus = iota
)

func (s ImportTransactionsFileStatus) String() string {
	switch s {
	case WARNING:
		return "warning"
	case ERROR:
		return "error"
	case OK:
		return "ok"
	}
	return "unknown"
}

type ImportTransactionsFileRequest struct {
	FileReader   io.Reader
	Filename     string
	FileSizeInMB float64
}

type ImportTransactionsFileResponse struct {
	Status                ImportTransactionsFileStatus
	Details               string
	TotalProcessedRecords int
	TotalValidRecords     int
}

type ListImportsResponse struct {
	ImportDate             time.Time
	TransactionsImportDate time.Time
}

type GetUserResponse struct {
	Name  string
	Email string
}

type Service interface {
	ImportTransactionsFile(*ImportTransactionsFileRequest) (*ImportTransactionsFileResponse, error)
	ListImports() ([]ListImportsResponse, error)
}
