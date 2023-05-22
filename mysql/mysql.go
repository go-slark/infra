package mysql

import (
	"fmt"
	xlogger "github.com/go-slark/slark/logger"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var dbInst *gorm.DB

type MySqlConfig struct {
	Address       string `json:"address"`
	MaxIdleConn   int    `json:"max_idle_conn"`
	MaxOpenConn   int    `json:"max_open_conn"`
	MaxLifeTime   int    `json:"max_life_time"`
	MaxIdleTime   int    `json:"max_idle_time"`
	LogMode       int    `json:"log_mode"` //默认warn
	CustomizedLog bool   `json:"customized_log"`
	xlogger.Logger
}

func InitMySql(c *MySqlConfig) {
	db, err := createDB(c)
	if err != nil {
		panic(errors.New(fmt.Sprintf("use %+v create mysql error %+v", c, err)))
	}
	dbInst = db
}

func createDB(c *MySqlConfig) (*gorm.DB, error) {
	var l logger.Interface
	if c.CustomizedLog {
		l = newCustomizedLogger(WithLogLevel(logger.LogLevel(c.LogMode)), WithLogger(c.Logger))
	} else {
		l = logger.Default.LogMode(logger.LogLevel(c.LogMode))
	}
	cfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         l,
	}
	db, err := gorm.Open(mysql.Open(c.Address), cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	if c.MaxLifeTime != 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifeTime) * time.Second)
	}
	if c.MaxIdleTime != 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(c.MaxIdleTime) * time.Second)
	}

	if err = sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, errors.WithStack(err)
	}

	if db == nil {
		return nil, errors.New("db is nil")
	}
	return db, nil
}

func GetDB() *gorm.DB {
	return dbInst
}

func Close() {
	if dbInst == nil {
		return
	}

	sqlDB, err := dbInst.DB()
	if err != nil {
		return
	}
	_ = sqlDB.Close()
}
