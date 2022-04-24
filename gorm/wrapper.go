package gorm

import (
	"gorm.io/gorm"
)

type DB = gorm.DB
type Config = gorm.Config

var ErrRecordNotFound = gorm.ErrRecordNotFound
