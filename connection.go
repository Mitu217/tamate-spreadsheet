package spreadsheet

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-tamate/tamate/driver"
	sheets "google.golang.org/api/sheets/v4"
)

type SpreadsheetConn struct {
	sheetID        string
	schemaRowIndex int
	service        SpreadsheetService
}

func newSpreadsheetConn(ctx context.Context, sheetID string, schemaRowIndex int) (*SpreadsheetConn, error) {
	if schemaRowIndex < 0 {
		return nil, fmt.Errorf("columnRowIndex is invalid value: %d", schemaRowIndex)
	}

	service, err := newGoogleSpreadsheetService(ctx)
	if err != nil {
		return nil, err
	}
	return &SpreadsheetConn{
		sheetID:        sheetID,
		schemaRowIndex: schemaRowIndex,
		service:        service,
	}, nil
}

func (c *SpreadsheetConn) Close() error {
	return nil
}

func (c *SpreadsheetConn) GetSchema(ctx context.Context, sheetName string) (*driver.Schema, error) {
	values, err := c.service.GetValues(ctx, c.sheetID, sheetName)
	if err != nil {
		return nil, err
	}
	rows := values[0].Values

	row := rows[c.schemaRowIndex]
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
		genericColName, err := columnNameFromSpreadsheetToGeneric(spreadsheetColName)
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
		genericColName, err := columnNameFromSpreadsheetToGeneric(spreadsheetColName)
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
	// delete old data
	values, err := c.service.GetValues(ctx, c.sheetID, sheetName)
	if err != nil {
		return err
	}
	rawRows := values[0].Values
	if len(rawRows) > c.schemaRowIndex {
		rawSchemaRow := rawRows[c.schemaRowIndex]
		clearLeftTopCell := cellFromGenericToSpreadsheet(c.schemaRowIndex, 0)
		clearRightDownCell := cellFromGenericToSpreadsheet(c.schemaRowIndex, len(rawSchemaRow))
		clearRange := fmt.Sprintf("%s!%s:%s", sheetName, clearLeftTopCell, clearRightDownCell)
		if err := c.service.ClearValues(ctx, c.sheetID, clearRange); err != nil {
			return err
		}
	}

	// setting new data
	columnNames := make([]interface{}, len(schema.Columns))
	for i, col := range schema.Columns {
		columnName, err := columnNameFromGenericToSpreadsheet(col, schema.PrimaryKey)
		if err != nil {
			return err
		}
		columnNames[i] = columnName
	}
	setLeftTopCell := cellFromGenericToSpreadsheet(c.schemaRowIndex, 0)
	setRightBottomCell := cellFromGenericToSpreadsheet(c.schemaRowIndex, len(columnNames))
	setRange := fmt.Sprintf("%s!%s:%s", sheetName, setLeftTopCell, setRightBottomCell)
	if err != nil {
		return err
	}
	valueRange := &sheets.ValueRange{
		Range: setRange,
		Values: [][]interface{}{
			columnNames,
		},
	}
	return c.service.SetValues(ctx, sheetName, valueRange)
}

func (c *SpreadsheetConn) GetRows(ctx context.Context, sheetName string) ([]*driver.Row, error) {
	values, err := c.service.GetValues(ctx, c.sheetID, sheetName)
	if err != nil {
		return nil, err
	}
	rawRows := values[0].Values

	// extract only row
	if len(rawRows) > c.schemaRowIndex {
		rawRowsWithoutColumn := make([][]interface{}, len(rawRows)-1)
		for rowIndex, rowValue := range rawRows {
			if rowIndex < c.schemaRowIndex {
				rawRowsWithoutColumn[rowIndex] = rowValue
			} else if rowIndex > c.schemaRowIndex {
				rawRowsWithoutColumn[rowIndex-1] = rowValue
			}
		}
		rawRows = rawRowsWithoutColumn
	}

	schema, err := c.GetSchema(ctx, sheetName)
	if err != nil {
		return nil, err
	}

	rows := make([]*driver.Row, len(rawRows))
	for rowIndex, rawRow := range rawRows {
		rowValues := make(driver.RowValues, len(schema.Columns))
		groupByKey := make(driver.GroupByKey)
		for colIndex, col := range schema.Columns {
			var colValue *driver.GenericColumnValue
			if colIndex < len(rawRow) {
				colValue = driver.NewGenericColumnValue(col, rawRow[colIndex].(string))
			} else {
				colValue = driver.NewGenericColumnValue(col, "")
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
}

func (c *SpreadsheetConn) SetRows(ctx context.Context, sheetName string, rows []*driver.Row) error {
	// delete old data
	values, err := c.service.GetValues(ctx, c.sheetID, sheetName)
	if err != nil {
		return err
	}
	rawRows := values[0].Values
	if err := c.service.ClearValues(ctx, c.sheetID, sheetName); err != nil {
		return err
	}

	// setting new data
	setValues := make([][]interface{}, 0)
	if len(rawRows) > c.schemaRowIndex {
		setValues = append(setValues, rawRows[c.schemaRowIndex])
	}
	for _, row := range rows {
		setValue := make([]interface{}, 0)
		for _, val := range row.Values {
			setValue = append(setValue, val.String())
		}
		setValues = append(setValues, setValue)
	}
	valueRange := &sheets.ValueRange{
		Range:  sheetName,
		Values: setValues,
	}
	return c.service.SetValues(ctx, c.sheetID, valueRange)
}

func columnNameFromSpreadsheetToGeneric(columnName string) (string, error) {
	reg := regexp.MustCompile("\\((.+?)\\)")
	if res := reg.FindStringSubmatch(columnName); len(res) >= 2 {
		return res[1], nil
	}
	return columnName, nil
}

func columnNameFromGenericToSpreadsheet(col *driver.Column, pk *driver.Key) (string, error) {
	if pk != nil && pk.KeyType == driver.KeyTypePrimary {
		for _, v := range pk.ColumnNames {
			if v == col.Name {
				return "(" + col.Name + ")", nil
			}
		}
	}
	return col.Name, nil
}

func cellFromGenericToSpreadsheet(row, column int) string {
	return fmt.Sprintf("R%dC%d", row+1, column+1)
}
