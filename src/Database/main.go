package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StudentMarks struct {
	StudentID   string
	TutorID     string
	Mark        string
	dateUpdated string
}

func generateNextSemStartDate() time.Time {
	var currentDate = time.Now()
	var daysUntilMon = (1 - int(currentDate.Weekday()) + 7) % 7
	var semStartDate = currentDate.AddDate(0, 0, daysUntilMon)

	return semStartDate
}

func generateNextSemEndDate(semStartDate time.Time) time.Time {
	semEndDate := semStartDate.AddDate(0, 0, 4)
	return semEndDate
}

func ServeHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Uni_LMS_Marks_Entry/1.0")

		return next(c)
	}
}

func main() {

	fmt.Println("Starting Database Service")

	//Create Echo HTTP Server
	e := echo.New()
	//Use custom server header dispalying applciation version
	e.Use(ServeHeader)
	fmt.Println("HTTP Server Created")

	fmt.Println(generateNextSemStartDate())
	fmt.Println(generateNextSemEndDate(generateNextSemStartDate()))

	//Group API version one routes together
	//g := e.Group("/api/V1")

	credential := options.Credential{
		Username: "admin",
		Password: "admin",
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://0.0.0.0:8059").SetAuth(credential))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases)

	currentSemester := "17-01-2022"

	collection := client.Database("StudentMarks").Collection(currentSemester)
	fmt.Println(collection)

	// Use goroutine to run http server synchronoulsy with other functions
	go func() {
		if err := e.Start(":8129"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	//Gracefully shutdown the server if an error happens
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
