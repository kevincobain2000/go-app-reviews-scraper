package app

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbMaster *gorm.DB

func NewDB() *gorm.DB {
	return NewMasterDB()
}

func setMasterDB() {
	c := NewConfig()

	if c.DBConfig.DBConnection == "mysql" {
		dbMaster = sqlConnection(c.DBConfig.DBHost)
		sqlDB, err := dbMaster.DB()

		if err != nil {
			panic("cannot connect to sql database")
		}
		configureSQL(sqlDB)
	}

	if c.DBConfig.DBConnection == "sqlite" {
		dbMaster = sqliteConnection()
	}
}

func sqlConnection(host string) *gorm.DB {
	c := NewConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DBConfig.DBUsername, c.DBConfig.DBPassword, host, c.DBConfig.DBPort, c.DBConfig.DBDatabase)
	log.Println("[info] dsn: ", dsn)

	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(c.DBConfig.DBLogLevel)),
	})
	if err != nil {
		panic("cannot connect to database")
	}
	return connection
}

func sqliteConnection() *gorm.DB {
	c := NewConfig()
	connection, err := gorm.Open(sqlite.Open(c.DBConfig.DBHost), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(c.DBConfig.DBLogLevel)),
	})
	if err != nil {
		panic("cannot connect to database")
	}
	return connection
}

func NewMasterDB() *gorm.DB {
	if dbMaster == nil {
		setMasterDB()
	}
	sqlDB, err := dbMaster.DB()
	_ = sqlDB

	if err != nil {
		panic("cannot get to sql database")
	}

	return dbMaster
}

func configureSQL(sqlDB *sql.DB) {
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(30)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(30)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// Default DB close on mysql is 8 hours, so we set way before that (1 min)
	// This can be increased to 1 hour as well
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(1))
}
