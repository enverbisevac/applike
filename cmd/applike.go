package main

import (
	"applike/pkg"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"os"
	"time"
)

func main() {
	var err error

	// Database
	log.Info("Connecting to database...")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	migrations := os.Getenv("DB_MIGRATIONS")
	// I forgot to put parsetime and loose time
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, dbName)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// migrations
	// Run migrations
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}
	log.Info(fmt.Sprintf("file://%s", migrations))
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrations), // file://path/to/directory
		"mysql", driver)

	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	log.Info("Database migrated")

	// repository
	repoAccount := pkg.NewMysqlRepository(db)
	service := pkg.NewService(repoAccount)
	log.Info("Connected")

	// sqs
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Endpoint: aws.String(os.Getenv("SQS_ENDPOINT")),
	})
	svc := sqs.New(sess)
	queue := pkg.NewSqsQueue(svc, os.Getenv("SQS_QUEUE_URL"))

	// main staff
	// start inserting
	go InsertTodoItem(service, queue)

	// routes
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.Lvl(log.DEBUG))
	// router
	pkg.MakeHandlers(e, service)
	e.Logger.Fatal(e.Start(":3000"))
}

func InsertTodoItem(service *pkg.Service, queue pkg.Queue) {
	for {
		item := &pkg.TodoItem{
			Description: "Some description",
			DueDate:     time.Now().Add(time.Hour * 3),
		}
		err := service.CreateTodoItem(item)
		if err == nil {
			err = queue.PublishMessage(*item)
			if err != nil {
				log.Error(err)
			}
		} else {
			log.Error(err)
		}
		time.Sleep(3 * time.Second)
	}
}
