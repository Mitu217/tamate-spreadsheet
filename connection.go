package spreadsheet

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-tamate/tamate/driver"
)

type SpreadsheetConn struct {
	SpreadsheetID  string
	ColumnRowIndex int
	service        *googleSpreadsheetService
}

func newSpreadsheetConn(ctx context.Context, sheetID string, columnRowIndex int) (*SpreadsheetConn, error) {
	if columnRowIndex < 0 {
		return nil, fmt.Errorf("columnRowIndex is invalid value: %d", columnRowIndex)
	}
	service, err := newGoogleSpreadsheetService(ctx)
	if err != nil {
		return nil, err
	}
	return &SpreadsheetConn{
		SpreadsheetID:  sheetID,
		ColumnRowIndex: columnRowIndex,
		service:        service,
	}, nil
}

func (c *SpreadsheetConn) Close() error {
	return nil
}

func (c *SpreadsheetConn) GetSchema(ctx context.Context, name string) (*driver.Schema, error) {
	values, err := c.getValues(ctx, name)
	if err != nil {
		return nil, err
	}
	primaryKey := &driver.Key{
		KeyType: driver.KeyTypePrimary,
	}
	cols := make([]*driver.Column, 0)
	for rowIndex, row := range values {
		if rowIndex != c.ColumnRowIndex {
			continue
		}
		for colIndex := range row {
			colName, ok := row[colIndex].(string)
			if !ok {
				return nil, errors.New("interface conversion: interface {} is not string")
			}
			// check primarykey
			reg := regexp.MustCompile("\\((.+?)\\)")
			if res := reg.FindStringSubmatch(colName); len(res) >= 2 {
				colName = res[1]
				primaryKey.ColumnNames = append(primaryKey.ColumnNames, colName)
			}
			cols = append(cols, &driver.Column{
				Name:            colName,
				OrdinalPosition: colIndex,
				Type:            driver.ColumnTypeString,
			})
		}
		break
	}
	return &driver.Schema{
		Name:       name,
		PrimaryKey: primaryKey,
		Columns:    cols,
	}, nil
}

func (c *SpreadsheetConn) SetSchema(ctx context.Context, name string, schema *driver.Schema) error {
	return fmt.Errorf("feature support")
}

func (c *SpreadsheetConn) GetRows(ctx context.Context, sheetName string) ([]*driver.Row, error) {
	values, err := c.getValues(ctx, sheetName)
	if err != nil {
		return nil, err
	}
	if len(values) > c.ColumnRowIndex {
		valuesWithoutColumn := make([][]interface{}, len(values)-1)
		for rowIndex, row := range values {
			if rowIndex < c.ColumnRowIndex {
				valuesWithoutColumn[rowIndex] = row
			} else if rowIndex > c.ColumnRowIndex {
				valuesWithoutColumn[rowIndex-1] = row
			}
		}
		values = valuesWithoutColumn
	}
	/*
		rows := make([]*driver.Row, len(values))
		for rowIndex, row := range values {
			rowValues := make(driver.RowValues, len(schema.Columns))
			groupByKey := make(driver.GroupByKey)
			for colIndex, col := range schema.Columns {
				var colValue *driver.GenericColumnValue
				if colIndex < len(row) {
					colValue = NewGenericColumnValue(col, row[colIndex].(string))
				} else {
					colValue = NewGenericColumnValue(col, "")
				}
				rowValues[col.Name] = colValue
				// grouping primarykey
				for i := range schema.PrimaryKey.ColumnNames {
					if schema.PrimaryKey.ColumnNames[i] == col.Name {
						key := schema.PrimaryKey.String()
						groupByKey[key] = append(groupByKey[key], colValue)
					}
				}
			}
			rows[rowIndex] = &driver.Row{GroupByKey: groupByKey, Values: rowValues}
		}
		return rows, nil
	*/
	return nil, nil
}

func (c *SpreadsheetConn) SetRows(ctx context.Context, name string, rows []*driver.Row) error {
	return fmt.Errorf("feature support")
}

func (c *SpreadsheetConn) getValues(ctx context.Context, sheetName string) ([][]interface{}, error) {
	return c.service.Get(ctx, c.SpreadsheetID, sheetName)
}
