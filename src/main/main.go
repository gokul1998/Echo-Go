package main

import (
	"fmt"
	"net/http"
	"time"

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
func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "you are in admin page")
}
func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Custom Header")
		c.Response().Header().Set("Location", "BLR")
		return next(c)
	}
}
func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	if username == "gokul" && password == "12345" {
		cookie := &http.Cookie{}
		cookie.Name = "SessionID"
		cookie.Value = "some_value"
		cookie.Expires = time.Now().Add(10 * time.Minute)
		c.SetCookie(cookie)
		return c.String(http.StatusOK, "you are logged in")
	}
	return c.String(http.StatusUnauthorized, "enter the correct credentials")
}
func main() {
	fmt.Println("hellllo world")
	e := echo.New()
	g := e.Group("/admin")
	//cookiegroup := e.Group("/cookie")
	g.Use(serverHeader)
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
	e.GET("/login", login)
	g.GET("/main", mainAdmin)
	e.GET("/", home)
	e.Start(":8000")
}
