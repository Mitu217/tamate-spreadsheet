# tamate-spreadsheet

[![LICENSE](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/go-tamate/tamate-spreadsheet?status.svg)](https://godoc.org/github.com/go-tamate/tamate-spreadsheet)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-tamate/tamate-spreadsheet)](https://goreportcard.com/report/github.com/go-tamate/tamate-spreadsheet)

[![CircleCI](https://circleci.com/gh/go-tamate/tamate-spreadsheet.svg?style=svg)](https://circleci.com/gh/go-tamate/tamate-spreadsheet)

A Spreadsheet-Driver for [go-tamate/tamate](https://godoc.org/github.com/go-tamate/tamate) package

---------------------------------------

  * [Features](#features)
  * [Requirements](#requirements)
  * [Installation](#installation)
  * [Usage](#usage)
    * [DSN](#dsn-data-source-name)
      * [Sheet ID](#sheet-id)
    * [Credential Data](#credential-data)
  * [Testing / Development](#testing--development)
  * [License](#license)

---------------------------------------

## Features
 * Column's row index make valiable
 * Setting ignore row

## Requirements
 * Go 1.12 or higher. We aim to support the 3 latest versions of Go.

---------------------------------------

## Installation
Simple install the package to your [$GOPATH](https://github.com/golang/go/wiki/GOPATH "GOPATH") with the [go tool](https://golang.org/cmd/go/ "go command") from shell:
```bash
$ go get -u github.com/go-tamate/tamate-spreadsheet
```
Make sure [Git is installed](https://git-scm.com/downloads) on your machine and in your system's `PATH`.

## Usage
_Tamate Driver_ is an implementation of `tamate/driver` interface.

Use `spreadsheet` as `driverName` and a valid [DSN](#dsn-data-source-name)  as `dataSourceName`:
```go
import  "github.com/go-tamate/tamate"
import  _ "github.com/go-tamate/tamate-spreadsheet"

ds, err := tamate.Open("spreadseet", "xxxxxxxxxxxx")
```

Use this to `Get`, `Set`, `GettingDiff`, etc.

### DSN (Data Source Name)

#### Sheet ID

```
https://docs.google.com/spreadsheets/d/xxxxxxxxxxxx
```

### Credential Data

1. Create _ServiceAccount_ at Google Cloud Console.

2. Getting credential file from Google Cloud Console.

3. Allow _ServiceAccount_ to access spreadsheet. For example add _ServiceAccount's_ email address to sharer.

4. Setting credential file path or credential data to env as _TAMATE_SPREADSHEET_CREDENTIAL_FILE_PATH_ or _TAMATE_SPREADSHEET_CREDENTIAL_DATA_.

## Testing / Development

Please execute the following command at the root of the project

```bash
docker-compose up -d
go test ./...
docker-compose down
```

---------------------------------------

## License
* MIT
    * see [LICENSE](./LICENSE)