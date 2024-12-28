package bootstrap

import (
	"fmt"
	"github.com/jeancaardo/go-app-event-notifier/pkg/utils"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectLocal func
func ConnectLocal(l utils.Logger) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		_ = l.CatchError(err)
		os.Exit(-1)
	}
	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}
	if os.Getenv("DATABASE_MIGRATE") == "true" {
		err = db.AutoMigrate(
		/* domain.domain{} */
		)
		_ = l.CatchError(err)
	}
	return db
}

func InitLogger() utils.Logger {
	loggers := []interface{}{utils.LogOption{Debug: true}}
	fmt.Println("Initializing logger...")
	if os.Getenv("SENTRY_ENABLED") == "true" {
		fmt.Println("Sentry enabled")
		sentry := utils.SentryOption{
			Dsn:         os.Getenv("SENTRY_DSN"),
			Environment: os.Getenv("ENVIRONMENT"),
		}
		loggers = append(loggers, sentry)
	}

	return utils.New(loggers...)
}
