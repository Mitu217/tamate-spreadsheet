package spreadsheet

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
)

const (
	KeyCredentialFilePath = "TAMATE_SPREADSHEET_CREDENTIAL_FILE_PATH"
	KeyCredentialData     = "TAMATE_SPREADSHEET_CREDENTIAL_DATA"
	KeySheetId            = "TAMATE_SPREADSHEET_SHEET_ID"
	KeySheetName          = "TAMATE_SPREADSHEET_SHEET_NAME"
)

func getClient(ctx context.Context, credentialData []byte) (*http.Client, error) {
	conf, err := google.JWTConfigFromJSON(credentialData, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	return conf.Client(ctx), nil
}

type SpreadsheetService interface {
	GetValues(ctx context.Context, spreadsheetID string, ranges ...string) ([][]interface{}, error)
}

type googleSpreadsheetService struct {
	service *sheets.SpreadsheetsService
}

func newGoogleSpreadsheetService(ctx context.Context) (SpreadsheetService, error) {
	strCredentialData := os.Getenv(KeyCredentialData)
	if strCredentialData != "" {
		return newGoogleSpreadsheetServiceFromData(ctx, []byte(strCredentialData))
	}
	credentialFilePath := os.Getenv(KeyCredentialFilePath)
	if credentialFilePath != "" {
		return newGoogleSpreadsheetServiceFromFile(ctx, credentialFilePath)
	}
	return nil, fmt.Errorf("env: %s or %s not set", KeyCredentialFilePath, KeyCredentialData)
}

func newGoogleSpreadsheetServiceFromFile(ctx context.Context, path string) (SpreadsheetService, error) {
	credentialData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return newGoogleSpreadsheetServiceFromData(ctx, credentialData)
}

func newGoogleSpreadsheetServiceFromData(ctx context.Context, data []byte) (SpreadsheetService, error) {
	client, err := getClient(ctx, data)
	if err != nil {
		return nil, err
	}
	service, err := sheets.New(client)
	if err != nil {
		return nil, err
	}
	return &googleSpreadsheetService{
		service: service.Spreadsheets,
	}, nil
}

func (s *googleSpreadsheetService) GetValues(ctx context.Context, spreadsheetID string, ranges ...string) ([][]interface{}, error) {
	valueRange, err := s.service.Values.BatchGet(spreadsheetID).Ranges(ranges...).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("%#v\n", valueRange)
}
