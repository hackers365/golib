package db_server

import (
	"fmt"
	"math"
	"sync"
	"time"
)

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbServer struct {
	MysqlCon *gorm.DB
}

var dbServer *DbServer

func NewDbServer(conf string) (*DbServer, error) {
	if dbServer != nil {
		return dbServer
	}

	db, err := gorm.Open(sqlite.Open(conf), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	dbServer = db

	return dbServer, nil

}
