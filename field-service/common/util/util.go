package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
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

func Float64(v float64) *float64 {
	return &v
}

func RupiahFormat(amount *float64) string {
	stringValue := "0"
	if amount != nil {
		humanizeValue := humanize.CommafWithDigits(*amount, 0)
		stringValue = strings.ReplaceAll(humanizeValue, ",", ".")
	}
	return fmt.Sprintf("Rp %s", stringValue)
}
