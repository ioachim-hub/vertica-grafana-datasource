package main

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func safeText(value string) string {
	if utf8.ValidString(value) {
		return value
	}

	return strings.ToValidUTF8(value, "\uFFFD")
}

func safeError(err error) error {
	if err == nil {
		return nil
	}

	return errors.New(safeText(err.Error()))
}

func healthError(err error) *backend.CheckHealthResult {
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusError,
		Message: safeText(err.Error()),
	}
}

func healthSuccess(version string) *backend.CheckHealthResult {
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Successfully connected to " + safeText(version),
	}
}
