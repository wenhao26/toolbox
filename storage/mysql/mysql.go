package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Connection configuration
type Option struct {
	Username    string
	Password    string
	Host        string
	Port        int
	Dbname      string
	TablePrefix string
	MaxIdleConn int
	MaxOpenConn int
	MaxLifetime time.Duration
}

// Storage structure
type Storage struct {
	DB    *gorm.DB
	SqlDB *sql.DB
}

// Create MySQL storage instance
func NewMySQLStorage(option *Option) (*Storage, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		option.Username,
		option.Password,
		option.Host,
		option.Port,
		option.Dbname,
	)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   option.TablePrefix,
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
	}

	// Set the number of
	// connection pools,
	// the maximum number of connections,
	// and the maximum reusable time
	sqlDB.SetMaxIdleConns(option.MaxIdleConn)
	sqlDB.SetMaxOpenConns(option.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(option.MaxLifetime)

	return &Storage{DB: db, SqlDB: sqlDB}, nil
}
