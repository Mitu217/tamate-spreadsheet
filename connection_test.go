package spreadsheet

import (
	"context"
	"testing"

	"github.com/go-tamate/tamate"
	"github.com/go-tamate/tamate/driver"
	"github.com/stretchr/testify/assert"
	sheets "google.golang.org/api/sheets/v4"
)

type TestSpreadsheetService struct {
	valueRanges []*sheets.ValueRange
}

func (s *TestSpreadsheetService) Ping(ctx context.Context, sheetID string) error {
	return nil
}

func (s *TestSpreadsheetService) GetValues(ctx context.Context, sheetID string, ranges ...string) ([]*sheets.ValueRange, error) {
	return s.valueRanges, nil
}

func (s *TestSpreadsheetService) SetValues(ctx context.Context, sheetID string, valueRanges ...*sheets.ValueRange) error {
	s.valueRanges = valueRanges
	return nil
}

func (s *TestSpreadsheetService) ClearValues(ctx context.Context, sheetID string, ranges ...string) error {
	// FIXME
	s.valueRanges = []*sheets.ValueRange{}
	return nil
}

type fakeSpreadsheetDriver struct {
	spreadsheetlDriver
	FakeOpen func(ctx context.Context, dsn string) (driver.Conn, error)
}

func (fd *fakeSpreadsheetDriver) Open(ctx context.Context, dsn string) (driver.Conn, error) {
	return fd.FakeOpen(ctx, dsn)
}

func Test_GetSchema(t *testing.T) {
	var (
		ctx            = context.Background()
		dsn            = ""
		sheetName      = ""
		schemaRowIndex = 0
		fakeDriverName = "GetSchema"
	)

	// Prepare
	fakeService := &TestSpreadsheetService{
		valueRanges: []*sheets.ValueRange{
			&sheets.ValueRange{
				Values: [][]interface{}{
					[]interface{}{"id", "name", "age"},
				},
			},
		},
	}
	fakeDriver := &fakeSpreadsheetDriver{
		FakeOpen: func(ctx context.Context, dsn string) (driver.Conn, error) {
			return &SpreadsheetConn{
				sheetID:        dsn,
				schemaRowIndex: schemaRowIndex,
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
		assert.Equal(t, fakeService.valueRanges[0].Values[schemaRowIndex][0], sc.Columns[0].Name)
		assert.Equal(t, driver.ColumnTypeString, sc.Columns[0].Type)
		assert.Equal(t, fakeService.valueRanges[0].Values[schemaRowIndex][1], sc.Columns[1].Name)
		assert.Equal(t, driver.ColumnTypeString, sc.Columns[1].Type)
		assert.Equal(t, fakeService.valueRanges[0].Values[schemaRowIndex][2], sc.Columns[2].Name)
		assert.Equal(t, driver.ColumnTypeString, sc.Columns[2].Type)
	}
}

func Test_SetSchema(t *testing.T) {
	var (
		ctx            = context.Background()
		dsn            = ""
		sheetName      = ""
		schemaRowIndex = 0
		fakeDriverName = "SetSchema"
	)

	// Prepare
	schema := &driver.Schema{
		Name: sheetName,
		PrimaryKey: &driver.Key{
			KeyType:     driver.KeyTypePrimary,
			ColumnNames: []string{},
		},
		Columns: []*driver.Column{
			driver.NewColumn("id", 0, driver.ColumnTypeInt, true, false),
			driver.NewColumn("name", 1, driver.ColumnTypeString, true, false),
		},
	}
	fakeService := &TestSpreadsheetService{
		valueRanges: []*sheets.ValueRange{
			&sheets.ValueRange{
				Values: [][]interface{}{},
			},
		},
	}
	fakeDriver := &fakeSpreadsheetDriver{
		FakeOpen: func(ctx context.Context, dsn string) (driver.Conn, error) {
			return &SpreadsheetConn{
				sheetID:        dsn,
				schemaRowIndex: schemaRowIndex,
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
	if assert.NoError(t, ds.SetSchema(ctx, sheetName, schema)) {
		sc, err := ds.GetSchema(ctx, sheetName)
		if assert.NoError(t, err) {
			assert.Equal(t, "id", sc.Columns[0].Name)
			assert.Equal(t, driver.ColumnTypeString, sc.Columns[0].Type)
			assert.Equal(t, "name", sc.Columns[1].Name)
			assert.Equal(t, driver.ColumnTypeString, sc.Columns[1].Type)
		}
	}
}
