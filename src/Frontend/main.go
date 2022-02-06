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
//var marksEntryRoute string = "http://localhost:8121/api/V1/marksSubmit"

// Struct to render html templates
type TemplateRenderer struct {
	templates *template.Template
}

// Struct to Initialize Module Information
type Module struct {
	ModuleCode string
	ModuleName string
}

// Struct to Initialize List of Modules
type ModuleList struct {
	List []Module
}

// Struct to Initialize Class
type Class struct {
	ClassCode string
	Schedule  string
	Tutor     string
	Capacity  string
}

// Struct to Initialize list of Class
type ClassList struct {
	List []Class
}

// Struct to Initialize Students
type Student struct {
	StudentID      string
	StudentName    string
	DOB            string
	Address        string
	PhoneNumber    string
	ModuleEnrolled string
	ClassForModule string
}

// Struct to Initialize List of Students
type StudentList struct {
	List []Student
}

// Struct to Initialize Tutor
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

// Generate Date of Next Semester's Start Date
func generateNextSemStartDate() time.Time {
	var currentDate = time.Now()
	var daysUntilMon = (1 - int(currentDate.Weekday()) + 7) % 7
	var semStartDate = currentDate.AddDate(0, 0, daysUntilMon)

	return semStartDate
}

// Generate Date of Next Semester's End Date
func generateNextSemEndDate(semStartDate time.Time) time.Time {
	semEndDate := semStartDate.AddDate(0, 0, 4)
	return semEndDate
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

	// Tutor T0002
	// Initialize variables
	var Modules2 ModuleList
	var Classes2 ClassList
	var Students2 StudentList

	// Temporary Data
	temp2Module := Module{
		ModuleCode: "SCS",
		ModuleName: "Server and Cloud Security",
	}

	temp2Module2 := Module{
		ModuleCode: "WEB",
		ModuleName: "Web Developnment",
	}

	temp2Module3 := Module{
		ModuleCode: "DL",
		ModuleName: "Deep Learning",
	}

	Modules2.List = append(Modules.List, temp2Module)
	Modules2.List = append(Modules.List, temp2Module2)
	Modules2.List = append(Modules.List, temp2Module3)

	temp2Class := Class{
		ClassCode: "SCS001",
		Schedule:  "17-01-2022",
		Capacity:  "10",
	}

	temp2Class2 := Class{
		ClassCode: "SCS002",
		Schedule:  "17-01-2022",
		Capacity:  "20",
	}

	temp2Class3 := Class{
		ClassCode: "SCS001",
		Schedule:  "17-01-2022",
		Capacity:  "10",
	}

	temp2Class4 := Class{
		ClassCode: "WEB002",
		Schedule:  "17-01-2022",
		Capacity:  "5",
	}

	Classes2.List = append(Classes.List, temp2Class)
	Classes2.List = append(Classes.List, temp2Class2)
	Classes2.List = append(Classes.List, temp2Class3)
	Classes2.List = append(Classes.List, temp2Class4)

	temp2Student := Student{
		StudentID:      "S003",
		StudentName:    "Caleb Goh",
		DOB:            "Someday",
		Address:        "Someplace",
		PhoneNumber:    "1234",
		ModuleEnrolled: "ADB",
		ClassForModule: "ADB001",
	}

	temp2Student2 := Student{
		StudentID:      "S004",
		StudentName:    "Danny Chan",
		DOB:            "Someday",
		Address:        "Someplace",
		PhoneNumber:    "5678",
		ModuleEnrolled: "ADB",
		ClassForModule: "ADB002",
	}

	temp2Student3 := Student{
		StudentID:      "S005",
		StudentName:    "Kenneth Teo",
		DOB:            "Someday",
		Address:        "Someplace",
		PhoneNumber:    "7890",
		ModuleEnrolled: "ETI",
		ClassForModule: "ETI001",
	}

	temp2Student4 := Student{
		StudentID:      "S006",
		StudentName:    "Kah Ho",
		DOB:            "Someday",
		Address:        "Someplace",
		PhoneNumber:    "6543",
		ModuleEnrolled: "ETI",
		ClassForModule: "ETI002",
	}

	temp2Student5 := Student{
		StudentID:      "S007",
		StudentName:    "Dong Kiat",
		DOB:            "Someday",
		Address:        "Someplace",
		PhoneNumber:    "6543",
		ModuleEnrolled: "ETI",
		ClassForModule: "ETI002",
	}

	Students2.List = append(Students.List, temp2Student)
	Students2.List = append(Students.List, temp2Student2)
	Students2.List = append(Students.List, temp2Student3)
	Students2.List = append(Students.List, temp2Student4)
	Students2.List = append(Students.List, temp2Student5)

	TutorName := TutorInfo.TutorID
	TutorID := tutorID
	ClassesVar := Classes2
	StudentsVar := Students2
	ModulesVar := Modules2

	if tutorID == "T0002" {
		TutorName = "Mr Andy Tan"
		TutorID = "T0002"
		ClassesVar = Classes2
		StudentsVar = Students2
		ModulesVar = Modules2
	} else {
		TutorName = "Mr Wesley Tan"
		TutorID = "T0001"
		ClassesVar = Classes
		StudentsVar = Students
		ModulesVar = Modules
	}

	// Return render of marksDashboard.html
	return c.Render(http.StatusOK, "marksDashboard.html", map[string]interface{}{
		"TutorID":      TutorID,
		"Name":         TutorName,
		"ClassesList":  ClassesVar.List,
		"StudentsList": StudentsVar.List,
		"ModulesList":  ModulesVar.List,
	})
}

// Function to send marks entry by tutor back to frontend service to be processed
func marksEntry(c echo.Context) error {
	// Get studentID
	studentID := c.Param("studentID")
	fmt.Println("Marks for Student received:" + studentID)

	// Get Marks Entered by tutor
	marksEntered := c.FormValue("Marks")
	fmt.Println("Marks Entered received: ", marksEntered)

	// Get Student ID and Marks entered, and fake data of Schedule and TutorID
	postBody, _ := json.Marshal(map[string]string{
		"StudentID": studentID,
		"Marks":     marksEntered,
		"Schedule":  "31-01-2022",
		"TutorID":   "T0001",
	})

	// Buffer the json
	responsebody := bytes.NewBuffer(postBody)

	// Send json to the following url
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

// Function to check if service is up and running
func checkAPI(c echo.Context) error {
	fmt.Println("Frontend Service has been pinged!")
	fmt.Println("Sending reply...")
	return c.String(http.StatusOK, "Service is up and running")
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

	// Routes the server is handling
	e.GET("/marksDashboard/:tutorID", marksDashboard)
	e.POST("/marksEntry/:studentID", marksEntry)
	e.GET("/checkapi", checkAPI)

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
	fmt.Println("Shutting down Marks Entry Frontend Service")
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
