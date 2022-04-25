# gorm-with-logrus

## Dependency

See `go.mod`.

## Usage

external/db.go

```go
package external

import (
	"os"
	"time"

	"github.com/2754github/gorm-with-logrus/gorm"
	"github.com/2754github/gorm-with-logrus/logger"
	"github.com/2754github/gorm-with-logrus/mysql"
)

func NewDB() *gorm.DB {
	db, err := mysql.New(
		&mysql.Config{
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Database: os.Getenv("DB_DATABASE"),
		},
		&gorm.Config{
			Logger: logger.New(200*time.Millisecond, false, logger.Warn),
		},
		&mysql.Option{
			MaxIdleConns:    100,
			MaxOpenConns:    100,
			ConnMaxLifetime: 100,
		},
	)
	if err != nil {
		panic(err)
	}

	return db
}

```

main.go

```go
package main

import (
	...
)

func main() {
    ...

    db := external.NewDB()

    ...
}

```
