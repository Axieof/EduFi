package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

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
	Lsit []Class
}

type Student struct {
	StudentID   string
	StudentName string
	DOB         string
	Address     string
	PhoneNumber string
}

type StudentList struct {
	List []Student
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func ServeHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Uni_LMS_Marks_Entry/1.0")

		return next(c)
	}
}

func marksDashboard(c echo.Context) error {

	tutorID := c.Param("tutorID")

	return c.Render(http.StatusOK, "marksDashboard.html", map[string]interface{}{
		"ModulesList":  ModulesList.List,
		"ClassesList":  ClassesList.List,
		"StudentsList": StudentsList.List,
	})
}

func main() {

	fmt.Println("Starting Frontend Service")

	//Create Echo HTTP Server
	e := echo.New()
	//Use custom server header dispalying applciation version
	e.Use(ServeHeader)
	fmt.Println("HTTP Server Created")

	//Group API version one routes together
	g := e.Group("/api/V1")

	g.GET("/marksDashboard/:tutorID", marksDashboard)

	// Use goroutine to run http server synchronoulsy with other functions
	go func() {
		if err := e.Start(":8120"); err != nil && err != http.ErrServerClosed {
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
