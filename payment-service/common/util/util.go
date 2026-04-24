package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
)

const defaultLimit = 10

type PaginationParam struct {
	Count int64 `json:"count"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Data  any   `json:"data"`
}

type PaginationResult struct {
	TotalPage    int   `json:"total_page"`
	TotalData    int64 `json:"total_data"`
	NextPage     *int  `json:"next_page"`
	PreviousPage *int  `json:"previous_page"`
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	Data         any   `json:"data"`
}

func GeneratePagination(params PaginationParam) PaginationResult {
	if params.Limit <= 0 {
		params.Limit = defaultLimit
	}

	totalPage := (int(params.Count) + params.Limit - 1) / params.Limit

	var (
		nextPage     int
		previousPage int
	)

	if params.Page < totalPage {
		nextPage = params.Page + 1
	}

	if params.Page > 1 {
		previousPage = params.Page - 1
	}

	return PaginationResult{
		TotalPage:    totalPage,
		TotalData:    params.Count,
		NextPage:     &nextPage,
		PreviousPage: &previousPage,
		Page:         params.Page,
		Limit:        params.Limit,
		Data:         params.Data,
	}
}

func GenerateSHA256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}

func RupiahFormat(amount *float64) string {
	stringValue := "0"
	if amount != nil {
		humanizeValue := humanize.CommafWithDigits(*amount, 0)
		stringValue = strings.ReplaceAll(humanizeValue, ",", ".")
	}
	return fmt.Sprintf("Rp %s", stringValue)
}

func add1(x int) int {
	return x + 1
}

func GetValueOrDefault[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

func GeneratePDFFromHTML(htmlTemplate string, data any) ([]byte, error) {
	funcMap := template.FuncMap{
		"add1": add1,
	}

	parsedTemplate, err := template.New("htmlTemplate").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	var filledTemplate bytes.Buffer
	err = parsedTemplate.Execute(&filledTemplate, data)
	if err != nil {
		return nil, err
	}

	htmlContent := filledTemplate.String()

	pdfGenerator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		logrus.Errorf("failed to create PDF generator: %v", err)
		return nil, err
	}

	pdfGenerator.Dpi.Set(600)
	pdfGenerator.NoCollate.Set(false)
	pdfGenerator.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfGenerator.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfGenerator.Grayscale.Set(false)
	pdfGenerator.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlContent)))

	err = pdfGenerator.Create()
	if err != nil {
		logrus.Errorf("failed to create PDF: %v", err)
		return nil, err
	}

	return pdfGenerator.Bytes(), nil
}
