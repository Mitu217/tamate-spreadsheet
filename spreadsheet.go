package spreadsheet

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
)

const (
	KeyCredentialFilePath = "TAMATE_SPREADSHEET_CREDENTIAL_FILE_PATH"
	KeySheetId            = "TAMATE_SPREADSHEET_SHEET_ID"
	KeySheetName          = "TAMATE_SPREADSHEET_SHEET_Name"
)

func getClient(credentialFilePath string) (*http.Client, error) {
	data, err := ioutil.ReadFile(credentialFilePath)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	return conf.Client(oauth2.NoContext), nil
}

type googleSpreadsheetService struct {
	service *sheets.SpreadsheetsService
}

func newGoogleSpreadsheetService(ctx context.Context) (*googleSpreadsheetService, error) {
	path := os.Getenv(KeyCredentialFilePath)
	if path == "" {
		return nil, fmt.Errorf("env: %s not set", KeyCredentialFilePath)
	}
	client, err := getClient(path)
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

func (s *googleSpreadsheetService) Get(ctx context.Context, spreadsheetID string, sheetName string) ([][]interface{}, error) {
	valueRange, err := s.service.Values.Get(spreadsheetID, sheetName).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return valueRange.Values, nil
}
