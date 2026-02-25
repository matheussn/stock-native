package domain

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var quantityPattern = regexp.MustCompile(`^\d+(\.\d{1,3})?$`)
var allowedMeasureUnits = map[MeasureUnit]struct{}{
	MeasureUnitKg:      {},
	MeasureUnitG:       {},
	MeasureUnitL:       {},
	MeasureUnitMl:      {},
	MeasureUnitUnidade: {},
}

func ValidateName(name string, field string) *AppError {
	if strings.TrimSpace(name) == "" {
		return NewAppError(ErrorValidation, field+" is required", map[string]string{"field": field})
	}
	return nil
}

func ValidateQuantityString(quantity string) (float64, *AppError) {
	trimmed := strings.TrimSpace(quantity)
	if !quantityPattern.MatchString(trimmed) {
		return 0, NewAppError(ErrorValidation, "quantity must be a positive number with up to 3 decimal places", map[string]string{"field": "quantity"})
	}

	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil || value <= 0 {
		return 0, NewAppError(ErrorValidation, "quantity must be greater than zero", map[string]string{"field": "quantity"})
	}
	return Round3(value), nil
}

func ValidateDateNotFuture(t time.Time) *AppError {
	if t.After(time.Now()) {
		return NewAppError(ErrorValidation, "movement date cannot be in the future", map[string]string{"field": "movementDate"})
	}
	return nil
}

func ValidateMeasureUnit(unit string) *AppError {
	normalized := MeasureUnit(strings.ToLower(strings.TrimSpace(unit)))
	if _, ok := allowedMeasureUnits[normalized]; !ok {
		return NewAppError(
			ErrorValidation,
			"measureUnit is invalid",
			map[string]any{
				"field":   "measureUnit",
				"allowed": []string{"kg", "g", "l", "ml", "unidade"},
			},
		)
	}
	return nil
}

func Round3(v float64) float64 {
	return math.Round(v*1000) / 1000
}
