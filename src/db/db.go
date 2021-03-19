package db

import (
	"fmt"
	"log"

	"github.com/ToshihitoKon/slack-ctf-score-bot/src/constants"
	"github.com/ToshihitoKon/slack-ctf-score-bot/src/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Tmdb struct {
	gorm *gorm.DB
}

var db *Tmdb

func DB() *Tmdb {
	if db != nil {
		return db
	}
	var (
		err    error
		dbhost = constants.DbHost
		dbport = constants.DbPort
		dbuser = constants.DbUser
		dbpass = constants.DbPass
		dbname = constants.DbName
		dsn    = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			dbuser, dbpass, dbhost, dbport, dbname,
		)
	)

	gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("err gorm.Open ", err.Error(), dsn)
	}
	db = &Tmdb{
		gorm: gormdb,
	}
	return db
}

func (db Tmdb) Insert(data interface{}) error {
	result := db.gorm.Create(data)
	if result.Error != nil {
		log.Println("db.Create: ", result)
	}
	return result.Error
}

// 使いたくねぇ
func (db Tmdb) GetGorm() *gorm.DB {
	return db.gorm
}

func (db Tmdb) Migrate() error {
	return db.gorm.AutoMigrate(&types.Transaction{})
}
