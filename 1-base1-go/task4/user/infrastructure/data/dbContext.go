package data

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GromDatabase struct {
	db *gorm.DB
}

func (g *GromDatabase) GetDb() *gorm.DB {
	return g.db
}

func InitDb() (Database, error) {
	host := "localhost"
	port := 3306
	user := "root"
	password := "1"
	database := "web3_task3"
	timeout := "10s"
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", user, password, host, port, database, timeout)

	db, err := gorm.Open(mysql.Open(connStr))

	// var users []entities.UserEntity
	// db.Model(&entities.UserEntity{}).Find(&users)
	// fmt.Println(users)

	if err != nil {
		return nil, err
	}

	return &GromDatabase{db: db}, nil
}
