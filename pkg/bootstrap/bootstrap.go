package bootstrap

import (
	"fmt"
	"github.com/jeancaardo/go-app-event-notifier/pkg/domain"
	"github.com/jeancaardo/go-app-event-notifier/pkg/utils/sentrykit"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/getsentry/sentry-go"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DatabaseMaxIdleConns = 2
const DatabaseMaxOpenConns = 10

// LoadEnv load env settings
func LoadEnv(logger log.Logger, source string) {
	_ = level.Info(logger).Log("boot", "Loading .env file...")
	var err error
	if source != "" {
		err = godotenv.Load(source)
	} else {
		err = godotenv.Overload()
	}

	if err != nil {
		_ = level.Info(logger).Log("boot", "Error loading .env file")
	}
}

// ConnectLocal func
func ConnectLocal() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbConfig, _ := db.DB()
	dbConfig.SetMaxIdleConns(DatabaseMaxIdleConns)
	dbConfig.SetMaxOpenConns(DatabaseMaxOpenConns)
	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}
	if os.Getenv("DATABASE_MIGRATE") == "true" {
		_ = db.AutoMigrate(&domain.User{})
		_ = db.AutoMigrate(&domain.Event{})
	}
	return db
}

// InitLogger -
func InitLogger() log.Logger {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", os.Getenv("APP_NAME"),
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	defer func() {
		_ = level.Info(logger).Log("boot", "service ended")
	}()
	_ = level.Info(logger).Log("boot", "service started")
	return logger
}

// InitSentry -
func InitSentry() log.Logger {
	fmt.Println("Initializing Sentry...")
	if os.Getenv("SENTRY_ENABLED") == "true" {
		fmt.Println("Sentry enabled")
		client, err := sentry.NewClient(sentry.ClientOptions{
			Dsn:         os.Getenv("SENTRY_DSN"),
			Environment: os.Getenv("ENVIRONMENT"),
		})
		if err != nil {
			fmt.Printf("sentry.Init: %s", err)
		}
		logger := sentrykit.NewSentryLogger(client)
		client.Flush(time.Second * 5)
		return logger
	}
	// dev prevent nil exception without sentry
	return InitLogger()
}

func InitSNS() *sns.SNS {
	sessionOptions := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if os.Getenv("ENVIRONMENT") == "local" {
		sessionOptions.Config = aws.Config{
			Region:      aws.String("us-east-1"),
			Endpoint:    aws.String("http://host.docker.internal:4566"),
			Credentials: credentials.NewStaticCredentials("foo", "bar", ""),
		}
	}
	sess := session.Must(session.NewSessionWithOptions(sessionOptions))
	return sns.New(sess)
}
