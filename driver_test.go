package spreadsheet

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-tamate/tamate"
	"github.com/stretchr/testify/assert"
)

func Test_Init(t *testing.T) {
	drivers := tamate.Drivers()
	d, has := drivers[driverName]
	assert.EqualValues(t, &spreadsheetlDriver{}, d)
	assert.True(t, has)
}

func Test_Open(t *testing.T) {
	var (
		keySheetId = "TAMATE_SPREADSHEET_SHEET_ID"
	)

	sheetID := os.Getenv(keySheetId)
	if sheetID == "" {
		t.Skip(fmt.Printf("env: %s not set", keySheetId))
	}

	dsn := fmt.Sprintf("%s", sheetID)
	ds, err := tamate.Open(driverName, dsn)
	defer func() {
		err := ds.Close()
		assert.NoError(t, err)
	}()
	assert.NoError(t, err)
}
