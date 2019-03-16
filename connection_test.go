package spreadsheet

import (
	"context"
	"testing"

	"github.com/go-tamate/tamate"
	"github.com/go-tamate/tamate/driver"
	"github.com/stretchr/testify/assert"
)

type TestSpreadsheetService struct {
	values [][]interface{}
}

func (s *TestSpreadsheetService) GetValues(ctx context.Context, sheetId string, sheetName string) ([][]interface{}, error) {
	return s.values, nil
}

func (s *TestSpreadsheetService) SetValues(ctx context.Context, sheetId string, sheetName string, values [][]interface{}) error {
	s.values = values
	return nil
}

type fakeSpreadsheetDriver struct {
	spreadsheetlDriver
	FakeOpen func(ctx context.Context, dsn string) (driver.Conn, error)
}

func (fd *fakeSpreadsheetDriver) Open(ctx context.Context, dsn string) (driver.Conn, error) {
	return fd.FakeOpen(ctx, dsn)
}

func TestSpreadsheet_GetSchema(t *testing.T) {
	var (
		ctx            = context.Background()
		dsn            = ""
		sheetName      = ""
		columnRowIndex = 0
		fakeDriverName = "GetSchema"
	)

	// Prepare
	fakeService := &TestSpreadsheetService{
		values: [][]interface{}{
			[]interface{}{"id", "name", "age"},
		},
	}
	fakeDriver := &fakeSpreadsheetDriver{
		FakeOpen: func(ctx context.Context, dsn string) (driver.Conn, error) {
			return &SpreadsheetConn{
				SpreadsheetID:  sheetName,
				ColumnRowIndex: columnRowIndex,
				service:        fakeService,
			}, nil
		},
	}
	tamate.Register(fakeDriverName, fakeDriver)

	// Open datasource
	ds, err := tamate.Open(fakeDriverName, dsn)
	defer ds.Close()
	assert.NoError(t, err)

	// Getting schema
	sc, err := ds.GetSchema(ctx, sheetName)
	if assert.NoError(t, err) {
		assert.Equal(t, fakeService.values[columnRowIndex][0], sc.Columns[0].Name)
		assert.Equal(t, driver.ColumnTypeString, sc.Columns[0].Type)
		assert.Equal(t, fakeService.values[columnRowIndex][1], sc.Columns[1].Name)
		assert.Equal(t, driver.ColumnTypeString, sc.Columns[1].Type)
		assert.Equal(t, fakeService.values[columnRowIndex][2], sc.Columns[2].Name)
		assert.Equal(t, driver.ColumnTypeString, sc.Columns[2].Type)
	}
}
