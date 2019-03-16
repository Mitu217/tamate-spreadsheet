package spreadsheet

import (
	"context"

	"github.com/go-tamate/tamate"
	"github.com/go-tamate/tamate/driver"
)

const driverName = "spreadsheet"

type spreadsheetlDriver struct{}

func (md *spreadsheetlDriver) Open(ctx context.Context, dsn string) (driver.Conn, error) {
	return newSpreadsheetConn(ctx, dsn, 0)
}

func init() {
	tamate.Register(driverName, &spreadsheetlDriver{})
}
