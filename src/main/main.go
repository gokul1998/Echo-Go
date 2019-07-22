package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Cat struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}
type Dog struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func home(c echo.Context) error {
	return c.String(http.StatusOK, "home page")
}
func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catAge := c.QueryParam("age")
	dataType := c.Param("type")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is %s \n your cat's age is %s", catName, catAge))
	}
	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"age":  catAge,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "please enter either string or json format",
	})
}
func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "you are in admin page")
}
func main() {
	fmt.Println("hellllo world")
	e := echo.New()
	g := e.Group("/admin")
	//g.Use(middleware.Logger()) ->simple logger
	//this is a logger with customized output
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host} ${latency_human}` + "\n",
	}))
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "gokul" && password == "12345" {
			return true, nil
		}
		return false, nil
	}))
	g.GET("/main", mainAdmin)
	e.GET("/", home)
	e.GET("/cats/:type", getCats)
	e.Start(":8000")
}
