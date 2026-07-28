package data

import (
	"github.com/dqixuan/stock_info/configs"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *configs.Mysql, loggerHelper log.Logger) (*Data, func(), error) {
	db, err := initDB(c, loggerHelper)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.NewHelper(loggerHelper).Info("closing the data resources")
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}

	return &Data{
		db: db,
	}, cleanup, nil
}

// DB returns the database instance
func (d *Data) DB() *gorm.DB {
	return d.db
}

// initDB initializes MySQL connection
func initDB(c *configs.Mysql, loggerHelper log.Logger) (*gorm.DB, error) {
	if c == nil {
		return nil, nil
	}

	db, err := gorm.Open(mysql.Open(c.Root), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.NewHelper(loggerHelper).Errorf("failed to connect database: %v", err)
		return nil, err
	}

	log.NewHelper(loggerHelper).Info("database connected successfully")
	return db, nil
}
