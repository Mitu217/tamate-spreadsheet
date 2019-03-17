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
)

func getClient(ctx context.Context, credentialData []byte) (*http.Client, error) {
	conf, err := google.JWTConfigFromJSON(credentialData, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	return conf.Client(ctx), nil
}

type SpreadsheetService interface {
	Ping(context.Context, string) error
	GetValues(context.Context, string, ...string) ([]*sheets.ValueRange, error)
	SetValues(context.Context, string, ...*sheets.ValueRange) error
	ClearValues(context.Context, string, ...string) error
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

func (s *googleSpreadsheetService) Ping(ctx context.Context, sheetID string) error {
	// Check if the sheet can be accessed as an alternative to ping
	_, err := s.service.Values.BatchGet(sheetID).Context(ctx).Do()
	return err
}

func (s *googleSpreadsheetService) GetValues(ctx context.Context, sheetID string, ranges ...string) ([]*sheets.ValueRange, error) {
	resp, err := s.service.Values.BatchGet(sheetID).Ranges(ranges...).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return resp.ValueRanges, nil
}

func (s *googleSpreadsheetService) SetValues(ctx context.Context, sheetID string, valuesRanges ...*sheets.ValueRange) error {
	// https://developers.google.com/sheets/api/reference/rest/v4/ValueInputOption
	valueInputOption := "USER_ENTERED"

	req := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: valueInputOption,
		Data:             valuesRanges,
	}
	_, err := s.service.Values.BatchUpdate(sheetID, req).Context(ctx).Do()
	return err
}

func (s *googleSpreadsheetService) ClearValues(ctx context.Context, sheetID string, ranges ...string) error {
	req := &sheets.BatchClearValuesRequest{
		Ranges: ranges,
	}
	_, err := s.service.Values.BatchClear(sheetID, req).Context(ctx).Do()
	return err
}
