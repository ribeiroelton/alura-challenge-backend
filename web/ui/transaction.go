package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/api"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/spi"
)

type TransactionsHandlerConfig struct {
	Service api.Transaction
	Log     spi.Logger
}

type TransactionsHandler struct {
	service api.Transaction
	log     spi.Logger
}

func NewTransactionsHandler(c *TransactionsHandlerConfig) *TransactionsHandler {
	return &TransactionsHandler{
		service: c.Service,
		log:     c.Log,
	}
}

func (h *TransactionsHandler) GetUpload(c echo.Context) error {
	imports, err := h.service.ListImports()
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	data["imports"] = imports
	err = c.Render(http.StatusOK, "upload-page.tmpl", data)
	if err != nil {
		return err
	}

	return nil
}

func (h *TransactionsHandler) PostUpload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		data := make(map[string]interface{})
		imports, err := h.service.ListImports()
		if err != nil {
			return err
		}
		data["status"] = "error"
		data["details"] = err.Error()
		data["imports"] = imports
		c.Render(http.StatusSeeOther, "upload-page.tmpl", data)
		return err
	}

	h.log.Infof("Filename: %v, Size in MB: %f", file.Filename, float64(file.Size)/1024/1024)

	m, err := file.Open()
	if err != nil {
		data := make(map[string]interface{})
		imports, err := h.service.ListImports()
		if err != nil {
			return err
		}
		data["status"] = "error"
		data["details"] = err.Error()
		data["imports"] = imports
		c.Render(http.StatusSeeOther, "upload-page.tmpl", data)
		return err
	}

	req := &api.ImportTransactionsFileRequest{
		FileReader:   m,
		Filename:     file.Filename,
		FileSizeInMB: float64(file.Size) / 1024 / 1024,
	}

	res, err := h.service.ImportTransactionsFile(req)
	if err != nil {
		data := make(map[string]interface{})
		imports, err := h.service.ListImports()
		if err != nil {
			return err
		}
		data["status"] = res.Status.String()
		data["details"] = res.Details
		data["imports"] = imports
		c.Render(http.StatusSeeOther, "upload-page.tmpl", data)
		h.log.Info(data)
		return err
	}

	data := make(map[string]interface{})
	data["status"] = res.Status.String()
	imports, err := h.service.ListImports()
	if err != nil {
		return err
	}
	data["success_records"] = res.TotalValidRecords
	data["total_records"] = res.TotalProcessedRecords
	data["imports"] = imports
	c.Render(http.StatusSeeOther, "upload-page.tmpl", data)

	return nil
}
