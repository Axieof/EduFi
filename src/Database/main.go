package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StudentMarks struct {
	StudentID   string `json: StudentID`
	StudentName string `json: StudentName`
	DOB         string `json: DOB`
	Address     string `json: Address`
	PhoneNumber string `json: PhoneNumber`
	Schedule    string `json: Schedule`
	TutorID     string `json: TutorID`
	Mark        string `json: Mark`
	dateUpdated string `json: dateUpdated`
}

type DatabaseClient struct {
	Client  mongo.Client
	Context context.Context
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

func getCurrentSemStart() time.Time {
	var nowdate = time.Now()
	var nowint = int(nowdate.Weekday())
	var current = nowint - 1
	var after = nowdate.AddDate(0, 0, -current)

	fmt.Println("Current sem start is: ", after)

	return after
}

func ServeHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Uni_LMS_Marks_Entry/1.0")

		return next(c)
	}
}

func checkCollectionExists(semester string) bool {
	DBClient := connectToDB()

	names, err := DBClient.Client.Database("StudentMarks").ListCollectionNames(DBClient.Context, bson.D{{"options.capped", true}})
	if err != nil {
		log.Printf("Failed to get collection names: %v", err)
		return false
	}

	for _, name := range names {
		if name == semester {
			log.Printf("The Collection exists!")
			return true
		} else {
			return false
		}
	}

	return false
}

func postMarks(c echo.Context) error {

	StudentMarkEntry := StudentMarks{}

	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&StudentMarkEntry)

	if err != nil {
		log.Fatalf("Failed reading the request body %s", err)
	} else {
		DBClient := connectToDB()

		currentSemester := StudentMarkEntry.Schedule

		fmt.Println(currentSemester)

		colExists := checkCollectionExists(currentSemester)

		if !colExists {
			DBClient.Client.Database("StudentMarks").CreateCollection(DBClient.Context, currentSemester)
		}

		StudentMarkscollection := DBClient.Client.Database("StudentMarks").Collection(currentSemester)

		result, err := StudentMarkscollection.InsertOne(DBClient.Context, StudentMarkEntry)
		if err != nil {
			fmt.Println("An error occured %s", err)
		} else {
			fmt.Println("Insert result type: ", reflect.TypeOf(result))
			fmt.Println("Insert APi result: ", result)
		}
	}

	return c.String(http.StatusOK, "Marks Entered into DB")
}

/*
func toBSON(v interface{}) (doc *bson.Document, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.UnMarshal(data, &doc)
	return
}
*/

func connectToDB() DatabaseClient {
	credential := options.Credential{
		Username: "admin",
		Password: "admin",
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://0.0.0.0:8128").SetAuth(credential))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	DBClient := DatabaseClient{
		Client:  *client,
		Context: ctx,
	}

	return DBClient
}

func checkNewSem(prevsem time.Time, cursem time.Time) time.Time {
	if prevsem != cursem {
		prevsem = cursem
	}

	return prevsem
}

func checkAPI(c echo.Context) error {
	fmt.Println("Database Service has been pinged!")
	fmt.Println("Sending reply...")
	return c.String(http.StatusOK, "Service is up and running")
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
	g := e.Group("/api/V1")

	g.POST("/database/postMarks", postMarks)
	g.GET("/checkapi", checkAPI)

	// Use goroutine to run http server synchronoulsy with other functions
	go func() {
		if err := e.Start(":8129"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	previoussemester := time.Now()
	currentsemesterchan := make(chan time.Time)
	currentsemester := <-currentsemesterchan

	go func() {
		currentsem := getCurrentSemStart()
		currentsemesterchan <- currentsem
	}()

	previoussemester = checkNewSem(previoussemester, currentsemester)

	go func() {
		fmt.Println(currentsemester)
		fmt.Println(previoussemester)
	}()

	//Gracefully shutdown the server if an error happens
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println("Shutting down Marks Entry Database Service")
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
