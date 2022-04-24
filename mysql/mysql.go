package mysql

import (
	"database/sql"

	m "gorm.io/driver/mysql"
	g "gorm.io/gorm"

	"github.com/2754github/gorm-with-logrus/gorm"
)

func New(mConfig config, gConfig *gorm.Config, opt *Option) (*gorm.DB, error) {
	db, err := g.Open(m.Open(mConfig.ToDSN()), gConfig)
	if err != nil {
		return nil, err
	}

	if err := setConnPool(db, opt); err != nil {
		return nil, err
	}

	return db, nil
}

func setConnPool(db *g.DB, opt *Option) error {
	if opt == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(opt.MaxIdleConns)
	sqlDB.SetMaxOpenConns(opt.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(opt.ConnMaxLifetime)

	return nil
}

func NewForTest(sqlDB *sql.DB) (*gorm.DB, error) {
	return g.Open(m.New(m.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true, // https://github.com/go-gorm/gorm/issues/3565#issuecomment-712113474
	}))
}
