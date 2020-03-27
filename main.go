package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	core "github.com/ledinhbao/ldbcore"
	strava "github.com/ledinhbao/ldbstrava"
)

const (
	dbInstance = string("GlobalDatabase")
)

func setDatabaseForGlobalUse(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(dbInstance, db)
		c.Next()
	}
}

func main() {
	router := gin.Default()

	db, err := gorm.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panic("Error to connect to database")
	}

	db.AutoMigrate(&core.Setting{})

	router.Use(setDatabaseForGlobalUse(db))

	strava.SetConfig(strava.Config{
		ClientID:          "44814",
		ClientSecret:      "c44a13c4308b3b834320ae5e3648d6c7855980a3",
		PathPrefix:        "/",
		PathSubscription:  "/subscription",
		SubscriptionDBKey: "strava_subscription_key",
		GlobalDatabase:    dbInstance,
	})
	strava.InitializeRoutes(router)
	go strava.CreateSubscription(db)
	router.Run(":9098")
}
