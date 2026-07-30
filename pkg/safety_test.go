package main

import (
	"context"
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"google.golang.org/protobuf/proto"
)

func TestSafeErrorProducesValidUTF8(t *testing.T) {
	err := errors.New(string([]byte{0xff, 'x'}))
	got := safeError(err)

	if !utf8.ValidString(got.Error()) {
		t.Fatal("error is not valid UTF-8")
	}
	response := &backend.QueryDataResponse{
		Responses: backend.Responses{
			"A": {Error: got},
		},
	}
	protobufResponse, err := backend.ToProto().QueryDataResponse(backend.DataFrameFormat_ARROW, response)
	if err != nil {
		t.Fatalf("convert query response to protobuf: %v", err)
	}
	if _, err := proto.Marshal(protobufResponse); err != nil {
		t.Fatalf("marshal query response: %v", err)
	}
}

func TestHealthSuccessProducesSerializableResponse(t *testing.T) {
	version := string([]byte{0xff, 'v'})
	result := healthSuccess(version)

	if !utf8.ValidString(result.Message) {
		t.Fatal("successful health message is not valid UTF-8")
	}
	if _, err := proto.Marshal(backend.ToProto().CheckHealthResponse(result)); err != nil {
		t.Fatalf("marshal health response: %v", err)
	}
}

func TestHealthErrorProducesSerializableResponse(t *testing.T) {
	err := errors.New(string([]byte{0xff, 'x'}))
	result := healthError(err)

	if !utf8.ValidString(result.Message) {
		t.Fatal("health message is not valid UTF-8")
	}
	if _, err := proto.Marshal(backend.ToProto().CheckHealthResponse(result)); err != nil {
		t.Fatalf("marshal health response: %v", err)
	}
}

func TestNewDataSourceInstanceRejectsInvalidRestoredPassword(t *testing.T) {
	invalidPassword := string([]byte{0xff, 'p'})
	settings := getSettings(configArgs{})
	settings.DecryptedSecureJSONData = map[string]string{"password": invalidPassword}

	_, err := newDataSourceInstance(context.Background(), settings)

	const want = "datasource password contains invalid UTF-8; re-save the credentials or restore Grafana with the original encryption key"
	if err == nil || err.Error() != want {
		t.Fatalf("newDataSourceInstance() error = %v, want %q", err, want)
	}
	if strings.Contains(err.Error(), invalidPassword) {
		t.Fatal("error exposes invalid password")
	}
}

func TestNewDataSourceInstanceRejectsInvalidRestoredOAuthToken(t *testing.T) {
	invalidToken := string([]byte{0xff, 't'})
	settings := getSettings(configArgs{})
	settings.DecryptedSecureJSONData = map[string]string{"OauthToken": invalidToken}

	_, err := newDataSourceInstance(context.Background(), settings)

	const want = "datasource OAuth token contains invalid UTF-8; re-save the credentials or restore Grafana with the original encryption key"
	if err == nil || err.Error() != want {
		t.Fatalf("newDataSourceInstance() error = %v, want %q", err, want)
	}
	if strings.Contains(err.Error(), invalidToken) {
		t.Fatal("error exposes invalid OAuth token")
	}
}
