package database

import (
	"github.com/voice-of-colombo/service/graph/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func OpenDatabase() {
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/testvalues?charset=utf8mb4&parseTime=True&loc=Local"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	Db.AutoMigrate(&model.Meeting{})
	Db.AutoMigrate(&model.User{})
	Db.AutoMigrate(&model.Speech{})
	Db.AutoMigrate(&model.AhCounts{})
	Db.AutoMigrate(&model.Club{})
	Db.AutoMigrate(&model.ClubUserMapping{})
}
