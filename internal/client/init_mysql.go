package client

//import (
//	"fmt"
//
//	"github.com/dqixuan/stock_info/internal/conf"
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//)
//
//// InitMySQL initializes MySQL client
//func InitMySQL(c *conf.Data_Database) (*gorm.DB, error) {
//	if c == nil || c.Source == "" {
//		return nil, fmt.Errorf("database config is empty")
//	}
//
//	db, err := gorm.Open(mysql.Open(c.Source), &gorm.Config{})
//	if err != nil {
//		return nil, fmt.Errorf("failed to connect to database: %w", err)
//	}
//
//	return db, nil
//}
//
