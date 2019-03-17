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

func (s *TestSpreadsheetService) GetValues(ctx context.Context, spreadsheetID string, ranges ...string) ([][]interface{}, error) {
	return s.values, nil
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

func TestSpreadsheet_SetSchema(t *testing.T) {
	var (
		ctx            = context.Background()
		dsn            = ""
		sheetName      = ""
		columnRowIndex = 0
		fakeDriverName = "SetSchema"
	)

	// Prepare
	fakeSchema := &driver.Schema{
		Name: sheetName,
		PrimaryKey: &driver.Key{
			KeyType:     driver.KeyTypePrimary,
			ColumnNames: []string{"id"},
		},
		Columns: []*driver.Column{
			driver.NewColumn("id", 0, driver.ColumnTypeInt, true, false),
			driver.NewColumn("name", 1, driver.ColumnTypeString, true, false),
		},
	}
	fakeService := &TestSpreadsheetService{
		values: [][]interface{}{},
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

	// Setting schema
	if assert.NoError(t, ds.SetSchema(ctx, sheetName, fakeSchema)) {

	}
}
