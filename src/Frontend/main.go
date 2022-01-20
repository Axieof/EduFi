package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func ServeHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Uni_LMS_Marks_Entry/1.0")

		return next(c)
	}
}

func main() {

	fmt.Println("Starting Frontend Service")

	//Create Echo HTTP Server
	e := echo.New()
	//Use custom server header dispalying applciation version
	e.Use(ServeHeader)
	fmt.Println("HTTP Server Created")

	//Group API version one routes together
	g := e.Group("frontend")

	//g.GET("/:tutorID", marksDashboard)
	//g.POST("/marksSubmit", marksSubmit)

}
