package spreadsheet

import (
	"context"
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	sheets "google.golang.org/api/sheets/v4"
)

func Test_Connect(t *testing.T) {
	var (
		ctx = context.Background()
	)

	// Prepare
	strCredentialData := os.Getenv(KeyCredentialData)
	credentialFilePath := os.Getenv(KeyCredentialFilePath)
	if strCredentialData == "" && credentialFilePath == "" {
		t.Skip(fmt.Printf("env: %s, %s not set", KeyCredentialData, KeyCredentialFilePath))
	}
	sheetId := os.Getenv(KeySheetId)
	if sheetId == "" {
		t.Skip(fmt.Printf("env: %s not set", KeySheetId))
	}
	sheetName := os.Getenv(KeySheetName)
	if sheetName == "" {
		t.Skip(fmt.Printf("env: %s not set", KeySheetName))
	}

	// Connect
	credentialData := *(*[]byte)(unsafe.Pointer(&strCredentialData))
	client, err := getClient(ctx, credentialData)
	if assert.NoError(t, err) {
		service, err := sheets.New(client)
		assert.NoError(t, err)

		// Check if the sheet can be accessed as an alternative to ping
		_, err = service.Spreadsheets.Values.Get(sheetId, sheetName).Context(ctx).Do()
		assert.NoError(t, err)
	}
}

func Test_GetValues(t *testing.T) {
	// Prepare
	strCredentialData := os.Getenv(KeyCredentialData)
	credentialFilePath := os.Getenv(KeyCredentialFilePath)
	if strCredentialData == "" && credentialFilePath == "" {
		t.Skip(fmt.Printf("env: %s, %s not set", KeyCredentialData, KeyCredentialFilePath))
	}
	sheetId := os.Getenv(KeySheetId)
	if sheetId == "" {
		t.Skip(fmt.Printf("env: %s not set", KeySheetId))
	}
	sheetName := os.Getenv(KeySheetName)
	if sheetName == "" {
		t.Skip(fmt.Printf("env: %s not set", KeySheetName))
	}

	// TODO:
}

func Test_SetValues(t *testing.T) {
	// Prepare
	strCredentialData := os.Getenv(KeyCredentialData)
	credentialFilePath := os.Getenv(KeyCredentialFilePath)
	if strCredentialData == "" && credentialFilePath == "" {
		t.Skip(fmt.Printf("env: %s, %s not set", KeyCredentialData, KeyCredentialFilePath))
	}
	sheetId := os.Getenv(KeySheetId)
	if sheetId == "" {
		t.Skip(fmt.Printf("env: %s not set", KeySheetId))
	}
	sheetName := os.Getenv(KeySheetName)
	if sheetName == "" {
		t.Skip(fmt.Printf("env: %s not set", KeySheetName))
	}

	// TODO:
}
