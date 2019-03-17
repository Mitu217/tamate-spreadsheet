package spreadsheet

import (
	"context"
	"fmt"
	"os"
	"testing"

	"google.golang.org/api/sheets/v4"

	"github.com/stretchr/testify/assert"
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
	service, err := newGoogleSpreadsheetService(ctx)
	assert.NoError(t, err)

	// Connect
	assert.NoError(t, service.Ping(ctx, sheetId))
}

func Test_GetValues(t *testing.T) {
	var (
		ctx       = context.Background()
		sheetName = "GetValues"
		// FIXME: setting before test
		fakeValues = [][]interface{}{
			[]interface{}{"id", "name", "age"},
			[]interface{}{"1", "hana", "15"},
		}
	)

	// Prepare
	strCredentialData := os.Getenv(KeyCredentialData)
	credentialFilePath := os.Getenv(KeyCredentialFilePath)
	if strCredentialData == "" && credentialFilePath == "" {
		t.Skip(fmt.Printf("env: %s, %s not set", KeyCredentialData, KeyCredentialFilePath))
	}
	sheetID := os.Getenv(KeySheetId)
	if sheetID == "" {
		t.Skip(fmt.Printf("env: %s not set", KeySheetId))
	}
	service, err := newGoogleSpreadsheetService(ctx)
	assert.NoError(t, err)

	// getting values
	valueRanges, err := service.GetValues(ctx, sheetID, sheetName)
	if assert.NoError(t, err) {
		assert.Equal(t, fakeValues[0][0], valueRanges[0].Values[0][0])
		assert.Equal(t, fakeValues[0][1], valueRanges[0].Values[0][1])
		assert.Equal(t, fakeValues[0][2], valueRanges[0].Values[0][2])
		assert.Equal(t, fakeValues[1][0], valueRanges[0].Values[1][0])
		assert.Equal(t, fakeValues[1][1], valueRanges[0].Values[1][1])
		assert.Equal(t, fakeValues[1][2], valueRanges[0].Values[1][2])
	}
}

func Test_SetValues(t *testing.T) {
	var (
		ctx        = context.Background()
		sheetName  = "SetValues"
		fakeValues = [][]interface{}{
			[]interface{}{"id", "name", "age"},
			[]interface{}{"1", "hana", "15"},
		}
	)

	// Prepare
	strCredentialData := os.Getenv(KeyCredentialData)
	credentialFilePath := os.Getenv(KeyCredentialFilePath)
	if strCredentialData == "" && credentialFilePath == "" {
		t.Skip(fmt.Printf("env: %s, %s not set", KeyCredentialData, KeyCredentialFilePath))
	}
	sheetID := os.Getenv(KeySheetId)
	if sheetID == "" {
		t.Skip(fmt.Printf("env: %s not set", KeySheetId))
	}
	service, err := newGoogleSpreadsheetService(ctx)
	assert.NoError(t, err)
	assert.NoError(t, service.ClearValues(ctx, sheetID, sheetName))

	// setting values
	valueRange := &sheets.ValueRange{
		Range:  sheetName,
		Values: fakeValues,
	}
	if assert.NoError(t, service.SetValues(ctx, sheetID, valueRange)) {
		valueRanges, err := service.GetValues(ctx, sheetID, sheetName)
		if assert.NoError(t, err) {
			assert.Equal(t, fakeValues[0][0], valueRanges[0].Values[0][0])
			assert.Equal(t, fakeValues[0][1], valueRanges[0].Values[0][1])
			assert.Equal(t, fakeValues[0][2], valueRanges[0].Values[0][2])
			assert.Equal(t, fakeValues[1][0], valueRanges[0].Values[1][0])
			assert.Equal(t, fakeValues[1][1], valueRanges[0].Values[1][1])
			assert.Equal(t, fakeValues[1][2], valueRanges[0].Values[1][2])
		}
	}
}

func Test_ClearValues(t *testing.T) {
	var (
		ctx       = context.Background()
		sheetName = "SetValues"
	)

	// Prepare
	strCredentialData := os.Getenv(KeyCredentialData)
	credentialFilePath := os.Getenv(KeyCredentialFilePath)
	if strCredentialData == "" && credentialFilePath == "" {
		t.Skip(fmt.Printf("env: %s, %s not set", KeyCredentialData, KeyCredentialFilePath))
	}
	sheetID := os.Getenv(KeySheetId)
	if sheetID == "" {
		t.Skip(fmt.Printf("env: %s not set", KeySheetId))
	}
	service, err := newGoogleSpreadsheetService(ctx)
	assert.NoError(t, err)
	assert.NoError(t, service.ClearValues(ctx, sheetID, sheetName+"!R1C1:R1C3"))
}
