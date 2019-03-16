package spreadsheet

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-tamate/tamate/driver"
)

type SpreadsheetConn struct {
	SpreadsheetID  string
	ColumnRowIndex int
	service        SpreadsheetService
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

func (c *SpreadsheetConn) GetSchema(ctx context.Context, sheetName string) (*driver.Schema, error) {
	values, err := c.service.GetValues(ctx, c.SpreadsheetID, sheetName)
	if err != nil {
		return nil, err
	}
	row := values[c.ColumnRowIndex]

	schema := &driver.Schema{
		Name:    sheetName,
		Columns: make([]*driver.Column, 0),
		PrimaryKey: &driver.Key{
			KeyType:     driver.KeyTypePrimary,
			ColumnNames: make([]string, 0),
		},
	}

	// Setting Columns
	for i, val := range row {
		spreadsheetColName, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("interface conversion: interface {} cannot cast string: %v", val)
		}
		genericColName, err := columnNameFromSpreadSheetToGeneric(spreadsheetColName)
		if err != nil {
			return nil, err
		}
		col := &driver.Column{
			Name:            genericColName,
			OrdinalPosition: i,
			Type:            driver.ColumnTypeString,
		}
		schema.Columns = append(schema.Columns, col)
	}

	// Setting PrimaryKey
	for _, val := range row {
		spreadsheetColName, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("interface conversion: interface {} cannot cast string: %v", val)
		}
		genericColName, err := columnNameFromSpreadSheetToGeneric(spreadsheetColName)
		if err != nil {
			return nil, err
		}
		if spreadsheetColName != genericColName {
			schema.PrimaryKey.ColumnNames = append(schema.PrimaryKey.ColumnNames, genericColName)
		}
	}

	return schema, nil
}

func (c *SpreadsheetConn) SetSchema(ctx context.Context, sheetName string, schema *driver.Schema) error {
	return nil
}

func (c *SpreadsheetConn) GetRows(ctx context.Context, sheetName string) ([]*driver.Row, error) {
	values, err := c.service.GetValues(ctx, c.SpreadsheetID, sheetName)
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

func columnNameFromSpreadSheetToGeneric(columnName string) (string, error) {
	reg := regexp.MustCompile("\\((.+?)\\)")
	if res := reg.FindStringSubmatch(columnName); len(res) >= 2 {
		return res[1], nil
	}
	return columnName, nil
}

func columnNameFromGenericToSpreadSheet(col *driver.Column, pk *driver.Key) (string, error) {
	if pk != nil && pk.KeyType == driver.KeyTypePrimary {
		for _, v := range pk.ColumnNames {
			if v == col.Name {
				return "(" + col.Name + ")", nil
			}
		}
	}
	return col.Name, nil
}
