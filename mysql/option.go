package mysql

import (
	"time"
)

type Option struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}
