package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

// Struct to Initialize Student
type Student struct {
	StudentID      string
	StudentName    string
	DOB            string
	Address        string
	PhoneNumber    string
	ModuleEnrolled string
	ClassForModule string
}

// Struct to Initialize StudentMarks
type StudentMarks struct {
	StudentID string `json: StudentID`
	Marks     string `json: Marks`
	Schedule  string `json: Schedule`
	TutorID   string `json: TutorID`
}

// Function to set header for HTTP server
func ServeHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Uni_LMS_Marks_Entry/1.0")

		return next(c)
	}
}

// Function to submit marks to Database Service and Marks Wallet Service
func marksSubmit(c echo.Context) error {
	studentID := c.Param("studentID")
	fmt.Println("Posting marks for Student:" + studentID)

	// Get Student information from student service
	tempStudent := Student{
		StudentID:   "S001",
		StudentName: "Pritheev Roshan",
		DOB:         "08-07-2002",
		Address:     "Someplace",
		PhoneNumber: "12345678",
	}

	// Temporary Fake Data
	tempModuleCode := "ADB"

	fmt.Println(tempStudent.StudentID)

	StudentMark := StudentMarks{}

	// Decode incoming student data
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&StudentMark)
	log.Printf("StudentID Retrieved: " + StudentMark.StudentID)
	log.Printf("Marks Retrieved: " + StudentMark.Marks)

	if err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return c.String(http.StatusFailedDependency, "Faield reading the request body")
	} else {
		//Post data to database
		fmt.Println("Posting data to database")
		postBody, _ := json.Marshal(map[string]string{
			"StudentID":   StudentMark.StudentID,
			"StudentName": tempStudent.StudentName,
			"DOB":         tempStudent.DOB,
			"Address":     tempStudent.Address,
			"PhoneNumber": tempStudent.PhoneNumber,
			"Mark":        StudentMark.Marks,
			"Schedule":    StudentMark.Schedule,
			"TutorID":     StudentMark.TutorID,
		})

		responsebody := bytes.NewBuffer(postBody)

		url := "http://10.31.11.11:8129/api/V1/database/postMarks"

		resp, err := http.Post(url, "application/json", responsebody)

		if err != nil {
			log.Fatalf("An error occured %s", err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		sb := string(body)

		fmt.Println(sb)

		// Post to Marks Wallet
		fmt.Println("Posting data to marks wallet")

		postBodywallet, _ := json.Marshal(map[string]string{
			"ttype": "Mark",
			"sid":   "Mark",
			"rid":   tempStudent.StudentID,
			"ts":    "2022-02-06 12:51:34",
			"tysm":  tempModuleCode,
			"ta":    StudentMark.Marks,
			"stat":  "ping",
		})

		responsebodywallet := bytes.NewBuffer(postBodywallet)

		urlwallet := "http://10.31.11.11:8053/Transaction/new"

		respwallet, errwallet := http.Post(urlwallet, "application/json", responsebodywallet)

		if errwallet != nil {
			log.Fatalf("An error occured %s", errwallet)
		}

		defer respwallet.Body.Close()

		bodywallet, err := ioutil.ReadAll(respwallet.Body)
		if err != nil {
			log.Fatalln(err)
		}
		sbwallet := string(bodywallet)

		fmt.Println(sbwallet)
	}

	return c.String(http.StatusOK, "Marks Entered")

	return c.String(http.StatusOK, "Posting Marks for StudentID: "+studentID)
}

// Function to check if service is up and running
func checkAPI(c echo.Context) error {
	fmt.Println("Marks Service has been pinged!")
	fmt.Println("Sending reply...")
	return c.String(http.StatusOK, "Service is up and running")
}

func main() {

	fmt.Println("Starting Marks Entry Service")

	//Create Echo HTTP Server
	e := echo.New()
	//Use custom server header dispalying applciation version
	e.Use(ServeHeader)
	fmt.Println("HTTP Server Created")

	//Group API version one routes together
	g := e.Group("/api/V1")

	g.POST("/marksSubmit/:studentID", marksSubmit)
	g.GET("/checkapi", checkAPI)

	// Use goroutine to run http server synchronoulsy with other functions
	go func() {
		if err := e.Start(":8121"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	//Gracefully shutdown the server if an error happens
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println("Shutting down Marks Entry Marks Service")
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
