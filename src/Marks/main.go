package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

func ServeHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Uni_LMS_Marks_Entry/1.0")

		return next(c)
	}
}

func marksDashboard(c echo.Context) error {
	tutorID := c.Param("tutorID")
	fmt.Println("Tutor ID received: " + tutorID)

	return c.String(http.StatusOK, "Marks Dashboard Accessed")
}

func marksSubmit(c echo.Context) error {
	// Get studentID and module from the query string
	studentID := c.QueryParam("studentID")
	moduleCode := c.QueryParam("moduleCode")
	marks := c.QueryParam("marks")

	// TODO
	// Post to database and receive reply

	return c.String(http.StatusOK, "Posting Marks for StudentID:"+studentID+", Module:"+moduleCode+", Marks: "+marks)
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

	g.GET("/:tutorID", marksDashboard)
	g.POST("/marksSubmit", marksSubmit)

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
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
