package mysql

import (
	"fmt"
)

type config interface {
	ToDSN() string
}

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func (c *Config) ToDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}
