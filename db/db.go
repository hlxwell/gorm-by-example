package db

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/hlxwell/gorm-by-example/plugin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Conn *gorm.DB
var ScopedConns map[uint64]*gorm.DB = make(map[uint64]*gorm.DB)

func InitDB() {
	makeConn("gorm_by_example")
}

func InitTestDB() {
	makeConn("gorm_by_example_test")
}

func ScopedConn() *gorm.DB {
	return ScopedConns[GoID()]
}

func SetScope(currentUserID uint) {
	// FIXME: I need a solution to expire old keys.
	ScopedConns[GoID()] = Conn.Set(plugin.CurrentUserIDKey, currentUserID)
}

// Helper Methods ============================

func makeConn(name string) {
	logLevel := logger.Error

	// custom logger
	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logLevel,
			Colorful:      true,
		},
	)

	// data source name
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		"root",
		"",
		"localhost",
		"3306",
		name,
	)

	// Init conn
	var err error
	Conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   customLogger,
	})

	if err != nil {
		panic(err)
	}

	Conn.Use(plugin.New(plugin.Config{
		DB:          Conn,
		AutoMigrate: true,
		Tables: []string{
			"User",
		},
	}))
}

func GoID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
