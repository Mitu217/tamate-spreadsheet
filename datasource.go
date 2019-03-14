package spreadsheet

import (
	"context"

	"github.com/Mitu217/tamate"
	"github.com/Mitu217/tamate/driver"
)

const driverName = "spreadsheet"

type mysqlDriver struct{}

func (md *mysqlDriver) Open(ctx context.Context, dsn string) (driver.Conn, error) {
	sc, err := newSpreadsheetConn(nil, dsn, 0)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func init() {
	tamate.Register(driverName, &mysqlDriver{})
}
