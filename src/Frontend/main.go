package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

//var tutorRoute string = "http://localhost:803?/"
//var moduleRoute string = "http://localhost:811?/"
var marksEntryRoute string = "http://localhost:8121/api/V1/marksSubmit"

type TemplateRenderer struct {
	templates *template.Template
}

type Module struct {
	ModuleCode string
	ModuleName string
}

type ModuleList struct {
	List []Module
}

type Class struct {
	ClassCode string
	Schedule  string
	Tutor     string
	Capacity  string
}

type ClassList struct {
	List []Class
}

type Student struct {
	StudentID      string
	StudentName    string
	DOB            string
	Address        string
	PhoneNumber    string
	ModuleEnrolled string
	ClassForModule string
}

type StudentList struct {
	List []Student
}

type Tutor struct {
	TutorID        string
	TutorName      string
	ModulesTaught  ModuleList
	ClassesTaught  ClassList
	StudentsTaught StudentList
}

// Function to initialize template rendered with echo framework
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

// Provides echo server with a header
func ServeHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Uni_LMS_Marks_Entry/1.0")

		return next(c)
	}
}

// Handles the get endpoint http://localhost:8120/api/V1/marksDashboard/:tutorID
func marksDashboard(c echo.Context) error {

	tutorID := c.Param("tutorID")
	fmt.Println("Tutor: " + tutorID + " accessed Marks Dashboard")

	// Get tutor data [Tutor Service not implemented]
	/*
		postBody, _ := json.Marshal(map[string]string{
			"TutorID": tutorID,
		})

		responsebody := bytes.NewBuffer(postBody)

		resp, err := http.Post(tutorRoute, "application/json", responsebody)

		if err != nil {
			log.FatalF("An error occured %s", err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatalln("An error occured %s", err)
		}

		sb := string(body)
	*/

	// Temporary data until tutor service is working
	TutorInfo := Tutor{
		TutorID:   "T0001",
		TutorName: "Mr Wesley",
	}

	// Get Module data [Module Service not implemented]
	/*
		postBody, _ := json.Marshal(map[string]string{
			"TutorID": tutorID,
		})

		responsebody := bytes.NewBuffer(postBody)

		resp, err := http.Post(tutorRoute, "application/json", responsebody)

		if err != nil {
			log.FatalF("An error occured %s", err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatalln("An error occured %s", err)
		}

		sb := string(body)
	*/

	// Initialize variables
	var Modules ModuleList
	var Classes ClassList
	var Students StudentList

	// Temporary Data
	tempModule := Module{
		ModuleCode: "ADB",
		ModuleName: "Advanced Databases",
	}

	tempModule2 := Module{
		ModuleCode: "ETI",
		ModuleName: "Emerging Trends in IT",
	}
	Modules.List = append(Modules.List, tempModule)
	Modules.List = append(Modules.List, tempModule2)

	tempClass := Class{
		ClassCode: "ADB001",
		Schedule:  "17-01-2022",
		Capacity:  "10",
	}

	tempClass2 := Class{
		ClassCode: "ADB002",
		Schedule:  "17-01-2022",
		Capacity:  "20",
	}

	tempClass3 := Class{
		ClassCode: "ETI001",
		Schedule:  "17-01-2022",
		Capacity:  "10",
	}

	tempClass4 := Class{
		ClassCode: "ETI002",
		Schedule:  "17-01-2022",
		Capacity:  "5",
	}

	Classes.List = append(Classes.List, tempClass)
	Classes.List = append(Classes.List, tempClass2)
	Classes.List = append(Classes.List, tempClass3)
	Classes.List = append(Classes.List, tempClass4)

	tempStudent := Student{
		StudentID:      "S001",
		StudentName:    "Guy1",
		DOB:            "Someday",
		Address:        "Someplace",
		PhoneNumber:    "1234",
		ModuleEnrolled: "ADB",
		ClassForModule: "ADB001",
	}

	tempStudent2 := Student{
		StudentID:      "S002",
		StudentName:    "Guy2",
		DOB:            "Someday",
		Address:        "Someplace",
		PhoneNumber:    "5678",
		ModuleEnrolled: "ADB",
		ClassForModule: "ADB002",
	}

	tempStudent3 := Student{
		StudentID:      "S003",
		StudentName:    "Guy3",
		DOB:            "Someday",
		Address:        "Someplace",
		PhoneNumber:    "7890",
		ModuleEnrolled: "ETI",
		ClassForModule: "ETI001",
	}

	tempStudent4 := Student{
		StudentID:      "S004",
		StudentName:    "Guy4",
		DOB:            "Someday",
		Address:        "Someplace",
		PhoneNumber:    "6543",
		ModuleEnrolled: "ETI",
		ClassForModule: "ETI002",
	}

	Students.List = append(Students.List, tempStudent)
	Students.List = append(Students.List, tempStudent2)
	Students.List = append(Students.List, tempStudent3)
	Students.List = append(Students.List, tempStudent4)

	return c.Render(http.StatusOK, "marksDashboard.html", map[string]interface{}{
		"TutorID":      TutorInfo.TutorID,
		"Name":         TutorInfo.TutorName,
		"ClassesList":  Classes.List,
		"StudentsList": Students.List,
		"ModulesList":  Modules.List,
	})
}

func marksEntry(c echo.Context) error {
	studentID := c.Param("studentID")
	fmt.Println("Posting marks for Student:" + studentID)

	marksEntered := c.FormValue("Marks")
	fmt.Println("Marks Entered: ", marksEntered)

	postBody, _ := json.Marshal(map[string]string{
		"StudentID": studentID,
		"Marks":     marksEntered,
	})

	responsebody := bytes.NewBuffer(postBody)

	url := "http://localhost:8121/api/V1/marksSubmit/" + studentID

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

	return c.String(http.StatusOK, "Test")
}

func main() {

	fmt.Println("Starting Frontend Service")

	//Create Echo HTTP Server
	e := echo.New()

	//Use custom server header dispalying applciation version
	e.Use(ServeHeader)

	// Initialzie renderer for echo server to use
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("html/*.html")),
	}

	e.Renderer = renderer

	fmt.Println("Server Renderer Initialized")

	//Group API version one routes together
	g := e.Group("/api/V1")

	// Routes the server is handling
	g.GET("/marksDashboard/:tutorID", marksDashboard)
	g.POST("/marksEntry/:studentID", marksEntry)

	// Use goroutine to run http server synchronoulsy with other functions
	go func() {
		if err := e.Start(":8120"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	fmt.Println("HTTP Server Created")

	//Gracefully shutdown the server if an error happens
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println("Shutting down Maarks Entry Frontend Service")
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
