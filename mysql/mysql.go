package mysql

import (
	"fmt"
	xlogger "github.com/go-slark/slark/logger"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"sync"
	"time"
)

var (
	mysqlPools = make(map[string]*gorm.DB)
	once       sync.Once
)

type MySqlConfig struct {
	Alias         string        `json:"alias"`
	Address       string        `json:"address"`
	MaxIdleConn   int           `json:"max_idle_conn"`
	MaxOpenConn   int           `json:"max_open_conn"`
	MaxLifeTime   time.Duration `json:"max_life_time"`
	MaxIdleTime   time.Duration `json:"max_idle_time"`
	LogMode       int           `json:"log_mode"` //默认warn
	CustomizedLog bool          `json:"customized_log"`
	xlogger.Logger
}

func InitMySql(configs []*MySqlConfig) {
	once.Do(func() {
		for _, c := range configs {
			if _, ok := mysqlPools[c.Alias]; ok {
				panic(errors.New("duplicate mysql pool: " + c.Alias))
			}
			p, err := createNewMySqlPool(c)
			if err != nil {
				panic(errors.New(fmt.Sprintf("mysql pool %+v error %v", c, err)))
			}
			mysqlPools[c.Alias] = p
		}
	})
}

func createNewMySqlPool(c *MySqlConfig) (*gorm.DB, error) {
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
		sqlDB.SetConnMaxLifetime(c.MaxLifeTime)
	}
	if c.MaxIdleTime != 0 {
		sqlDB.SetConnMaxIdleTime(c.MaxIdleTime)
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

func GetMySql(alias string) *gorm.DB {
	return mysqlPools[alias]
}

func CloseMysql() {
	for _, db := range mysqlPools {
		if db == nil {
			continue
		}
		sqlDB, err := db.DB()
		if err != nil {
			continue
		}
		_ = sqlDB.Close()
	}
}
